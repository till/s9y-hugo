package hugo

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const tplSection string = "section.tmpl"

// struct to supply data to the template
type tplData struct {
	Name        string
	Description string
	Link        string
}

type Section struct {
	base string
	fs   embed.FS
}

func New(base string, fs embed.FS) *Section {
	return &Section{
		base: base,
		fs:   fs,
	}
}

func (s *Section) Create(name, description, link string) error {
	path := filepath.Join(s.base, strings.ToLower(name))
	if err := os.Mkdir(path, 0755); err != nil {
		if !os.IsExist(err) {
			return fmt.Errorf("os.Mkdir (section): %v", err)
		}
	}

	file, err := os.Create(filepath.Join(path, "_index.md"))
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return fmt.Errorf("os.Create (%s): %v", path, err)
	}
	defer file.Close()

	data := tplData{
		Name:        name,
		Description: description,
		Link:        link,
	}

	tpl := template.Must(template.New(tplSection).ParseFS(s.fs, tplSection))
	err = tpl.Execute(file, data)
	if err != nil {
		log.Fatal(fmt.Errorf("tpl.Execute (%s): %v", tplSection, err))
	}

	return nil
}
