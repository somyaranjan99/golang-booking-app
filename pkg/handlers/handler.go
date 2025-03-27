package handlers

import (
	"fmt"
	"github/somyaranjan99/basic-go-project/pkg/config"
	"github/somyaranjan99/basic-go-project/pkg/models"
	"github/somyaranjan99/basic-go-project/pkg/render"
	"net/http"
)

type Repository struct {
	Repo *config.AppConfig
}

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{Repo: a}
}

func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr

	repo.Repo.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.RenderTemplate(w, r, "home.page.tmpl", repo.Repo, &models.TemplateData{})
}

func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {
	data := &models.TemplateData{
		StringMap: map[string]string{},
	}
	// Populate .StringMap with data
	data.StringMap["remote_ip"] = repo.Repo.Session.GetString(r.Context(), "remote_ip")
	render.RenderTemplate(w, r, "about.page.tmpl", repo.Repo, data)
}
func (repo *Repository) Contact(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "contact.page.tmpl", repo.Repo, &models.TemplateData{})
}
func (repo *Repository) Generals(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "contact.page.tmpl", repo.Repo, &models.TemplateData{})
}
func (repo *Repository) Majors(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "majors.page.tmpl", repo.Repo, &models.TemplateData{})
}
func (repo *Repository) Reservation(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "make-reservation.page.tmpl", repo.Repo, &models.TemplateData{})
}
func (repo *Repository) Aavailability(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "search-availability.page.tmpl", repo.Repo, &models.TemplateData{})
}

func (respo *Repository) PostSearchAvailability(w http.ResponseWriter, r *http.Request) {

}

func (repo *Repository) PostAavailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	w.Write([]byte(fmt.Sprintf("start date is %s end date is %s", start, end)))
}
