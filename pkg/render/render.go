package render

import (
	"bytes"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/keroda/bookings/pkg/config"
	"github.com/keroda/bookings/pkg/models"
)

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template

	//get template cache from config
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplatecache()
	}

	// get requested page from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("could not get page template from cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

	_ = t.Execute(buf, td)

	// render the page
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println("error writing page to browser", err)
	}
}

func CreateTemplatecache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all html files
	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.html")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts

	}

	return myCache, nil
}