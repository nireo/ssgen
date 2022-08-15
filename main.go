package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
	"github.com/nireo/ssgen/render"
	"github.com/nireo/ssgen/setup"
	"gopkg.in/ini.v1"
)

func main() {
	// do all of the necessary setup
	ssgenDir := flag.String("dir", "./", "The directory in which ssgen settings and posts are stored.")
	genDirectory := flag.Bool("gen", false, "Should the given directory be generated with the basic files.")
	flag.Parse()

	if ssgenDir == nil {
		fmt.Println("ssgen --dir=<directory-path> --gen=<optional|generate default directory>")
		os.Exit(1)
	}

	if *genDirectory {
		if err := setup.SetupDirectory(*ssgenDir); err != nil {
			fmt.Printf("[ssgen] cannot setup directory: %s", err)
			os.Exit(1)
		}
	}

	metadataCfg, err := ini.Load(filepath.Join(*ssgenDir, "metadata.ini"))
	if err != nil {
		fmt.Printf("[ssgen] cannot read metadata file: %s", err)
		os.Exit(1)
	}

	rendererMetadata := render.Metadata{
		Title:  metadataCfg.Section("site-meta").Key("title").String(),
		Author: metadataCfg.Section("site-meta").Key("author").String(),
	}

	rdr, err := render.New(*ssgenDir, rendererMetadata)
	if err != nil {
		panic(err)
	}

	router := httprouter.New()
	router.GET("/", rdr.RenderHomePage)
	router.GET("/post/:post", rdr.RenderPostPage)

	log.Fatalln(http.ListenAndServe(":8080", router))
}
