package render

import (
	"github/somyaranjan99/basic-go-project/pkg/config"
	"github/somyaranjan99/basic-go-project/pkg/models"
	"html/template"
	"net/http"

	"github.com/justinas/nosurf"
)

func RenderTemplate(w http.ResponseWriter, r *http.Request, templ string, app *config.AppConfig, m *models.TemplateData) {
	csrfToken := nosurf.Token(r)
	m.CSRFToken = csrfToken
	m.Error = app.Session.PopString(r.Context(), "error")
	m.Flash = app.Session.PopString(r.Context(), "flash")
	m.Warning = app.Session.PopString(r.Context(), "warning")
	renderPage, err := template.ParseFiles("../../template/"+templ, "../../template/base.layout.tmpl")
	if err != nil {
		return
	}
	renderPage.Execute(w, m)
}

// var tc = make(map[string]*template.Template)

// func Render(w http.ResponseWriter, t string) {
// 	var templ *template.Template
// 	_, ok := tc[t]
// 	if !ok {
// 		err := creatingCache(t)
// 		if err != nil {
// 			fmt.Fprintf(w, "running from new %v", err)
// 			return
// 		}
// 	} else {
// 		fmt.Println("running from cache")

// 	}
// 	templ = tc[t]
// 	_ = templ.Execute(w, nil)
// }

// func creatingCache(t string) error {
// 	templates := []string{
// 		fmt.Sprintf("../../template/%s", t), "../../template/base.layout.tmpl",
// 	}
// 	templ, err := template.ParseFiles(templates...)
// 	if err != nil {
// 		return err
// 	}
// 	tc[t] = templ
// 	return nil

// }

// func RenderTemplate(w http.ResponseWriter, tmpl string, app *config.AppConfig, m *models.TemplateData) {
// 	var tc map[string]*template.Template

// 	if app.UseCache {
// 		// get the template cache from the app config
// 		tc = app.TemplateCache
// 	} else {
// 		tc, _ = CreateTemplateCache()
// 	}

// 	t, ok := tc[tmpl]
// 	if !ok {
// 		log.Fatal("Could not get template from template cache")
// 	}

// 	buf := new(bytes.Buffer)

// 	err := t.Execute(buf, nil)
// 	if err != nil {
// 		log.Println("Error executing template:", err)
// 	}

// 	_, err = buf.WriteTo(w)
// 	if err != nil {
// 		log.Println("Error writing template to browser:", err)
// 	}
// }

// func CreateTemplateCache() (map[string]*template.Template, error) {

// 	myCache := map[string]*template.Template{}

// 	// get all of the files named *.page.tmpl from ./templates
// 	pages, err := filepath.Glob("../../template/*.page.tmpl")
// 	if err != nil {
// 		return myCache, err
// 	}

// 	// range through all files ending with *.page.tmpl
// 	for _, page := range pages {
// 		name := filepath.Base(page)
// 		ts, err := template.New(name).ParseFiles(page)
// 		if err != nil {
// 			return myCache, err
// 		}

// 		matches, err := filepath.Glob("../../template/*.layout.tmpl")
// 		if err != nil {
// 			return myCache, err
// 		}

// 		if len(matches) > 0 {
// 			ts, err = ts.ParseGlob("../../template/*.layout.tmpl")
// 			if err != nil {
// 				return myCache, err
// 			}
// 		}

// 		myCache[name] = ts
// 	}

// 	return myCache, nil
// }
