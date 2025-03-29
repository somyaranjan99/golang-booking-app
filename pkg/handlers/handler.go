package handlers

import (
	"fmt"
	"github/somyaranjan99/basic-go-project/cmd/web/middleware/forms"
	"github/somyaranjan99/basic-go-project/pkg/config"
	"github/somyaranjan99/basic-go-project/pkg/models"
	"github/somyaranjan99/basic-go-project/pkg/render"
	"log"
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

func (repo *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation
	render.RenderTemplate(w, r, "make-reservation.page.tmpl", repo.Repo, &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (repo *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}
	reservation := &models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}
	form := forms.New(&r.PostForm)
	form.Has("first_name", r)
	form.Has("last_name", r)
	form.Has("email", r)
	form.Has("phone", r)
	// form.Required("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 3, r)
	form.MinLength("last_name", 3, r)
	form.IsValidEmail("email")
	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.RenderTemplate(w, r, "make-reservation.page.tmpl", repo.Repo, &models.TemplateData{
			Data: data,
			Form: form,
		})
		return
	}
	repo.Repo.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}
func (repo *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := repo.Repo.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		log.Println("can't get item from session")
		repo.Repo.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	repo.Repo.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.RenderTemplate(w, r, "reservation-summary.page.tmpl", repo.Repo, &models.TemplateData{
		Data: data,
	})

}
