package components

import (
	"fmt"
	"io"
	"sort"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type Styles map[string]any

func (c Styles) Render(w io.Writer) error {
	included := make([]string, 0, len(c))
	for format, v := range c {
		included = append(included, fmt.Sprintf(format, v))
	}
	sort.Strings(included)
	return h.Style(strings.Join(included, ";")).Render(w)
}

func (c Styles) Type() g.NodeType { return g.AttributeType }
