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
