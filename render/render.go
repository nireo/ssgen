package render

import "github.com/yuin/goldmark"

type Renderer struct {
	md        goldmark.Markdown
	directory string
}
