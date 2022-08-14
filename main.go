package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nireo/ssgen/render"
	"github.com/nireo/ssgen/setup"
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

	rdr, err := render.New(*ssgenDir)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/post", rdr.RenderPostPage)
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
