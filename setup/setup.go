package setup

import (
	"os"
	"path/filepath"
)

const basicMessage = `
[site-meta]
title = "ssgen default"
author = "not set"
`

func genSiteMetadata(baseDirectory string) error {
	file, err := os.Create(filepath.Join(baseDirectory, "metadata.ini"))
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write([]byte(basicMessage)); err != nil {
		return err
	}

	return nil
}

func genCustomTheme(baseDirectory string) error {
	file, err := os.Create(filepath.Join(baseDirectory, "theme.css"))
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

func genPostDirectory(baseDirectory string) error {
	return os.Mkdir(filepath.Join(baseDirectory, "posts"), os.ModePerm)
}

func SetupDirectory(baseDirectory string) error {
	if _, err := os.Stat(baseDirectory); os.IsNotExist(err) {
		if err := os.Mkdir(baseDirectory, os.ModePerm); err != nil {
			return err
		}
	}

	if err := genSiteMetadata(baseDirectory); err != nil {
		return err
	}

	if err := genPostDirectory(baseDirectory); err != nil {
		return err
	}

	if err := genCustomTheme(baseDirectory); err != nil {
	}

	return nil
}
