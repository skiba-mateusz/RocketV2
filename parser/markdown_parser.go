package parser

import (
	"fmt"
	"html/template"
	"io"

	"github.com/adrg/frontmatter"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type Metadata struct {
	Title 	string 	 	`yaml:"title"`
	Tags 	[]string 	`yaml:"tags"`
	Date	string		`yaml:"date"`
}

type Page struct {
	Meta Metadata
	Content template.HTML
	Permalink template.URL
	Section string
}

type MarkdownParser struct{}

func NewMarkdwonParser() *MarkdownParser {
	return &MarkdownParser{}
}

func (p *MarkdownParser) Parse(r io.Reader) (*Page, error) {
	var meta Metadata

	body, err := frontmatter.Parse(r, &meta)
	if err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %v", err)
	}

	content := p.mdToHTML(body)

	return &Page{
		Meta: meta,
		Content: content,
	}, nil
}

func (p *MarkdownParser) mdToHTML(md []byte) template.HTML {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	parser := parser.NewWithExtensions(extensions)
	doc := parser.Parse(md)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	renderedHtml := markdown.Render(doc, renderer)

	return template.HTML(renderedHtml)
}