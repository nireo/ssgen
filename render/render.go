package render

import (
	"os"
	"path/filepath"
)

type Renderer struct {
	// all of the posts are cached in memory since they don't take up that much space
	// and they're quicker to retrieve.
	posts     map[string][]byte
	directory string
}

func New(dir string) (*Renderer, error) {
	rdr := &Renderer{
		posts:     make(map[string][]byte),
		directory: dir,
	}

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		return nil
	})

	return rdr, nil
}
