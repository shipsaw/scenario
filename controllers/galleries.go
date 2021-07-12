package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/shipsaw/scenario/models"
	"github.com/shipsaw/scenario/views"
)

type Galleries struct {
	NewView *views.View
	gs      models.GalleryService
}

type GalleryForm struct {
	Title string `schema:"title"`
}

// POST /galleries
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.NewView.Render(w, vd)
		return
	}
	gallery := models.Gallery{
		Title: form.Title,
	}
	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.NewView.Render(w, vd)
		return
	}
	fmt.Fprint(w, gallery)
}
func NewGalleries(gs models.GalleryService) *Galleries {
	return &Galleries{
		NewView: views.NewView("bootstrap", "views/galleries/new.gohtml"),
		gs:      gs,
	}
}
