package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/msavanan/bookings/internal/config"
	"github.com/msavanan/bookings/internal/models"
)

// func RenderTest(w http.ResponseWriter, templ string) {

// 	parsedFiles, err := template.ParseFiles(filepath.Join("../../templates", templ),
// 		filepath.Join("../../templates", "base.layout.tmpl"))
// 	if err != nil {
// 		log.Println("error failed to parse the template file", err)
// 	}

// 	if err = parsedFiles.Execute(w, nil); err != nil {
// 		log.Println(err.Error())
// 	}

// }

// var tc = make(map[string]*template.Template)

// func RenderMethod2(w http.ResponseWriter, t string) {
// 	var tmpl *template.Template
// 	var err error

// 	_, inMap := tc[t]

// 	if !inMap {
// 		log.Println("Create templates......")

// 		err = createTemplate(t)
// 		if err != nil {
// 			log.Println("error failed to parse the template file", err)
// 		}

// 	} else {

// 		log.Println("Cached templates......")

// 	}

// 	tmpl = tc[t]

// 	err = tmpl.Execute(w, nil)
// 	if err != nil {
// 		log.Println("error failed to execute template file", err.Error())
// 	}

// }

// func createTemplate(t string) error {
// 	//When running from the root directory
// 	// go run cmd/web/main.go
// 	templates := []string{filepath.Join("templates", t),
// 		filepath.Join("templates", "base.layout.tmpl")}
// 	//When running from the cmd/web directory, directory where main file exist
// 	// go run cmd/web/main.go
// 	// templates := []string{filepath.Join("../../templates", t),
// 	// 	filepath.Join("../../templates", "base.layout.tmpl")}
// 	parsedFiles, err := template.ParseFiles(templates...)
// 	if err != nil {
// 		return err
// 	}

// 	tc[t] = parsedFiles

// 	return nil
// }

var app *config.AppConfig

var pathToTemplates = "./templates"

func NewTemplate(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.CSRFToken = nosurf.Token(r)
	return td
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var tc map[string]*template.Template
	var err error

	//Create a template cache
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, err = CreateTemplateCache()
		if err != nil {
			log.Fatal(err)
			return err
		}

	}

	//Get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Println("Failed to fetch template")
		return errors.New("Can't fetch template")
	}

	// Render the template
	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	err = t.Execute(buf, td)
	if err != nil {
		log.Println(err)
		return err
	}

	if _, err = buf.WriteTo(w); err != nil {
		log.Println(err)
		return err
	}

	return nil

}

func CreateTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	//get all of the files named *,page.html from ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return cache, err
	}

	// range through all files ending with *.page.html
	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return cache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return cache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return cache, err
			}
		}

		cache[name] = ts

	}
	return cache, nil

}
