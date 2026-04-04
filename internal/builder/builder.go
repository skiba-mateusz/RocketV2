package builder

import (
	"context"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/skiba-mateusz/RocketV2/internal/config"
	"github.com/skiba-mateusz/RocketV2/internal/parser"
	"github.com/skiba-mateusz/RocketV2/internal/templater"
	"github.com/skiba-mateusz/RocketV2/pkg/logger"
	"golang.org/x/sync/errgroup"
)

type kind int

const (
	single kind = iota
	list
)

type Node struct {
	Page	 	 *parser.Page
	Parent		 *Node
	Children 	 []*Node
	Kind 	 	 kind
	Permalink    template.URL
	Section		 string
	SourcePath	 string
	OutPath		 string
	IsDir 		 bool	
}

type Builder interface {
	Build(ctx context.Context) error
}

type DefaultBuilder struct {
	logger 			*logger.Logger
	config 			*config.Config
	pageParser 		parser.Parser
	templater  		templater.Templater
	root 			*Node
	nodes			map[string]*Node
	counter 		*atomic.Uint32
}

func NewDefault(logger *logger.Logger, config *config.Config, pageParser parser.Parser, templater templater.Templater) *DefaultBuilder {
	return &DefaultBuilder{
		logger: logger,
		config: config,
		pageParser: pageParser,
		templater: templater,
		root: nil,
		nodes: make(map[string]*Node),
		counter: &atomic.Uint32{},
	}
}

func (b *DefaultBuilder) Build(ctx context.Context) error {
	b.reset()
	
	start := time.Now()

	b.logger.Info("Starting build process...")

	if err := b.cleanBuildDir(); err != nil {
		return err
	}

	if err := b.prepareNodes(); err != nil {
		return err
	}

	if err := b.processContent(ctx); err != nil {
		return err
	}
	
	b.sortNodes()
	
	if err := b.build(ctx); err != nil {
		return err
	}

	if err := b.copyStaticAssets(); err != nil {
		return err
	}

	elapsed := time.Since(start)
	b.logger.Success("Build completed in %.2f, Total pages: %d", elapsed.Seconds(), b.counter.Load())

	return nil
}

func (b *DefaultBuilder) prepareNodes() error {
	b.logger.Debug("Preparing nodes...")

	return filepath.WalkDir(b.config.ContentDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && d.Name() == "_index.md" {
			parentPath := filepath.Dir(path)
			if parent, ok := b.nodes[parentPath]; ok {
				parent.SourcePath = path
				parent.Kind = list
			}

			return nil
		}

		node := &Node{
			SourcePath: path,
			IsDir: d.IsDir(),
			Kind: single,
		}

		b.nodes[path] = node

		if path != b.config.ContentDir {
			parentPath := filepath.Dir(path)
			if parent, ok := b.nodes[parentPath]; ok {
				node.Parent = parent
				parent.Children = append(parent.Children, node)
			}
		} else {
			b.root = node
		}

		return nil
	})
}

func (b *DefaultBuilder) processContent(ctx context.Context) error {
	b.logger.Debug("Processing content...")

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(runtime.NumCPU())

	for _, node := range b.nodes {
		n := node
		if node.SourcePath == "" || (n.IsDir && n.Kind != list) {
			continue
		}
		
		g.Go(func() error {
			if ctx.Err() != nil {
				if ctx.Err() == context.Canceled {
					return nil
				}
				return ctx.Err()
			}

			if err := b.loadContent(n); err != nil {
				b.logger.Warn("Skipping %s: %v", node.SourcePath, err)
			}

			return nil
		})
	}

	return g.Wait()
}

func (b *DefaultBuilder) loadContent(node *Node) error {
	file, err := os.Open(node.SourcePath)
	if err != nil {
		return fmt.Errorf("failed to open %s: %v", node.SourcePath, err)
	}

	page, err := b.pageParser.Parse(file)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %v", node.SourcePath, err)
	}

	node.Page = page

	return b.populateNode(node)
}

func (b *DefaultBuilder) populateNode(node *Node) error {
    relPath, err := filepath.Rel(b.config.ContentDir, node.SourcePath)
    if err != nil {
        return err
    }

    cleanPath := filepath.ToSlash(strings.TrimSuffix(relPath, filepath.Ext(relPath)))

    var dir string
    base := path.Base(cleanPath)

	section := path.Dir(cleanPath)
	if section == "." {
		section = ""
	}
	node.Section = section

    if base == "index" || base == "_index" || cleanPath == "." {
        dir = path.Dir(cleanPath)
    } else {
        dir = cleanPath
    }

    permalink := path.Join("/", dir)
    if permalink != "/" {
        permalink += "/"
    }

    node.Permalink = template.URL(permalink)
    node.OutPath = filepath.Join(b.config.BuildDir, filepath.FromSlash(dir), "index.html")

	return nil
}

func (b *DefaultBuilder) build(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(runtime.NumCPU())

	for _, node := range b.nodes {
		n := node
		if n.Page == nil {
			continue
		}

		g.Go(func() error {
			if ctx.Err() != nil {
				if ctx.Err() == context.Canceled {
					return nil
				}
				return ctx.Err()
			}	

			if err := b.renderNode(n); err != nil {
				b.logger.Warn("Skipping %s: %v", node.SourcePath, err)
			}

			return nil
		})
	}

	return g.Wait()
}

func (b *DefaultBuilder) renderNode(node *Node) error {
	b.logger.Debug("Rendering page %s -> %s", node.SourcePath, node.OutPath)

	var err error

	if node.Kind == list {
		err = b.renderPaginated(node)
	} else {
		err = b.renderSingle(node)
	}

	if err == nil {
		b.counter.Add(1)
	}

	return err
}

func (b *DefaultBuilder) renderSingle(node *Node) error {

	data := struct {
		Config  *config.Config
		Root 	*Node
		*Node
	} {
		b.config,
		b.root,
		node,
	}

	return b.render(node, data)
}

func (b *DefaultBuilder) renderPaginated(node *Node) error {
	pageSize := b.config.Paginate
	if pageSize <= 0 {
		pageSize = 5
	}

	totalItems := len(node.Children)
	totalPages := (totalItems + pageSize - 1) / pageSize

	pages := make([]int, totalPages)
	for p := range totalPages {
		pages[p] = p + 1
	}

	outPath := node.OutPath

	for i := range totalPages {
		start := pageSize * i
		end := start + pageSize
		if end > totalItems {
			end = totalItems
		}

		page := i + 1
		items := node.Children[start:end]

		outPathPage := outPath 
		if page > 1 {
			outPathPage = filepath.Join(filepath.Dir(outPath), "page", fmt.Sprint(page), "index.html")
		}

		data := struct {
			Config  	*config.Config
			Root 		*Node
			Items		[]*Node
			Pages		[]int
			TotalItems 	int
			CurrentPage int
			TotalPages	int
			HasPrevPage bool
			HasNextPage bool
			*Node
		} {
			b.config,
			b.root,
			items,
			pages,
			totalItems,
			page,
			totalPages,
			page > 1,
			page < totalPages,
			node,
		}

		tmp := *node
		tmp.OutPath = outPathPage
		if err := b.render(&tmp, data); err != nil {
			return fmt.Errorf("failed to render: %v", err)
		}
	}

	return nil
}

func (b *DefaultBuilder) render(node *Node, data any) error {
	if err := os.MkdirAll(filepath.Dir(node.OutPath), 0755); err != nil {
		return fmt.Errorf("failed to create dir %s: %v", filepath.Dir(node.OutPath), err)
	}

	file, err := os.Create(node.OutPath)
	if err != nil {
		return fmt.Errorf("failed to create %s: %v", node.OutPath, err)
	}
	defer file.Close()

	layout := "single.html"
	if node.Kind == list {
		layout = "list.html"
	}

	return b.templater.Render(file, "baseof.html", data, []string{filepath.Join(b.config.LayoutDir, node.Section, layout)})
}

func (b *DefaultBuilder) sortNodes() {
	for _, node := range b.nodes {
		if len(node.Children) > 1 {
			sort.Slice(node.Children, func(i, j int) bool {
				if node.Children[i].Page == nil {
					return false
				}

				if node.Children[j].Page == nil {
					return true
				}

				dateI, err := time.Parse(time.RFC3339, node.Children[i].Page.Meta.Date)
				if err != nil {
					return false
				}

				dateJ, err := time.Parse(time.RFC3339, node.Children[j].Page.Meta.Date)
				if err != nil {
					return true
				}

				return dateI.After(dateJ)
			})
		}
	}
}

func (b *DefaultBuilder) cleanBuildDir() error {
	b.logger.Debug("Cleaning build directory...")

	if _, err := os.Stat(b.config.BuildDir); os.IsNotExist(err) {
		return nil
	}

	if err := os.RemoveAll(b.config.BuildDir); err != nil {
		return fmt.Errorf("Failed to remove %s: %v", b.config.BuildDir, err)
	}

	return nil
}

func (b *DefaultBuilder) copyStaticAssets() error {
	b.logger.Debug("Copying static assets...")

	wd, _ := os.Getwd()
	source := filepath.Join(wd, b.config.StaticDir)
	destRoot := filepath.Join(wd, b.config.BuildDir, "static")

	return filepath.WalkDir(source, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return  err
		}

		rel, _ := filepath.Rel(source, path)
		dest := filepath.Join(destRoot, rel)

		if d.IsDir() {
			if err := os.MkdirAll(dest, 0755); err != nil {
				return fmt.Errorf("failed to create dir %s: %v", dest, err)
			}
			return nil
		} 

		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return fmt.Errorf("failed to create dir %s: %v", filepath.Dir(dest), err)
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %v", path, err)
		}
		if err := os.WriteFile(dest, content, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %v", dest, err)
		}

		return nil
	})
}

func (b *DefaultBuilder) reset() {
	b.root = nil
	b.nodes = make(map[string]*Node)
	b.counter.Store(0)
	b.templater.ClearCache()
}