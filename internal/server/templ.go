package server

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const HOST_ROOM = "host_room"

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
func newTemplate() *Templates {
	// using the function
	rootDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	tarDir := path.Join(rootDir, "/web/view/*html")
	return &Templates{
		templates: template.Must(template.ParseGlob(tarDir)),
	}
}
func (t *Templates) ServeMuxHandle(w io.Writer, name string, data interface{}) error {
	baseName := strings.TrimSuffix(name, filepath.Ext(name))
	switch baseName {
	case HOST_ROOM:
		t.Render(w, baseName, nil)
	default:
		return errors.New("Not found")
	}
	return nil
}
