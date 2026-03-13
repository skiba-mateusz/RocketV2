package templater

import "io"

type Templater interface {
	Render(w io.Writer, layout string, data any, templates []string) error
}