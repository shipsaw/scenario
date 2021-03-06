package views

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
)

var LayoutFilePath string = "views/layouts/*.gohtml"

type View struct {
	Template *template.Template
	Layout   string
}

func NewView(layout string, files ...string) *View {
	files = append(files, layoutFiles()...)

	template, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{
		Template: template,
		Layout:   layout,
	}
}

func (v *View) Render(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	switch data.(type) {
	case Data:
		// Do nothing
	default:
		data = Data{
			Yield: data,
		}
	}
	var buf bytes.Buffer
	if err := v.Template.ExecuteTemplate(&buf, v.Layout, data); err != nil {
		http.Error(w, "Something went wrong. If the problem presists please email support", http.StatusInternalServerError)
	}
	io.Copy(w, &buf)
}

func layoutFiles() []string {
	files, err := filepath.Glob(LayoutFilePath)
	if err != nil {
		panic(err)
	}
	return files
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, nil)
}
