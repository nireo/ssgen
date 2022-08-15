package render

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type Metadata struct {
	Title  string
	Author string
}

type Renderer struct {
	// all of the posts are cached in memory since they don't take up that much space
	// and they're quicker to retrieve.
	posts     map[string][]byte
	directory string
	md        goldmark.Markdown
	meta      Metadata
}

type PageRender struct {
	HTML template.HTML
}

func New(dir string, meta Metadata) (*Renderer, error) {
	rdr := &Renderer{
		posts:     make(map[string][]byte),
		directory: dir,
	}

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() || filepath.Ext(path) != ".md" { // ignore random directories and files
			return nil
		}

		name := info.Name()

		rdr.posts[name[:len(name)-3]], err = ioutil.ReadFile(path)
		if err != nil {
			return nil // just skip reading the file since something bad has happened.
		}

		return nil
	})

	rdr.md = goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	rdr.meta = meta

	return rdr, nil
}

func (re *Renderer) RenderPostPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	postName := ps.ByName("post")

	mdData, ok := re.posts[postName]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var buf bytes.Buffer
	if err := re.md.Convert(mdData, &buf); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pr := PageRender{
		HTML: template.HTML(buf.Bytes()),
	}

	tmpl := template.Must(template.ParseFiles("post.html"))

	if err := tmpl.Execute(w, pr); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type HomePost struct {
	Name string
}

type HomeRender struct {
	Posts []HomePost
	Meta  Metadata
}

func (re *Renderer) RenderHomePage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	posts := make([]HomePost, 0)
	for name := range re.posts {
		posts = append(posts, HomePost{
			Name: name,
		})
	}

	pr := HomeRender{
		Posts: posts,
		Meta:  re.meta,
	}

	tmpl := template.Must(template.ParseFiles("home.html"))
	if err := tmpl.Execute(w, pr); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
