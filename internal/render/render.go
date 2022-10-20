package render

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	"github.com/justinas/nosurf"
	"github.com/keroda/bookings/internal/config"
	"github.com/keroda/bookings/internal/models"
)

var functions = template.FuncMap{
	"niceDate":   NiceDate,
	"formatDate": FormatDate,
	"iterate":    Iterate,
	"add":        Add,
}

var app *config.AppConfig

func NewRenderer(a *config.AppConfig) {
	app = a
}

func NiceDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func FormatDate(t time.Time, f string) string {
	return t.Format(f)
}

func Add(a, b int) int {
	return a + b
}

func Iterate(count int) []int {
	var i int
	var items []int
	for i = 0; i < count; i++ {
		items = append(items, i)
	}
	return items
}

var pathToTemplates = "./templates"

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1
	}
	return td
}

func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
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
		log.Fatal("could not get page template from cache(" + tmpl + ")")
		return errors.New("could not get page template from cache (" + tmpl + ")")
	} //else {
	//log.Println("Rendering", tmpl)
	//}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)

	// render the page
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println("error writing page to browser", err)
		return err
	}

	return nil
}

func CreateTemplatecache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all html files
	pages, err := filepath.Glob(pathToTemplates + "/*.page.html")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(pathToTemplates + "/*.layout.html")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(pathToTemplates + "/*.layout.html")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts

	}

	return myCache, nil
}
