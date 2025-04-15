package handlers

import (
	"fmt"
	"github/somyaranjan99/basic-go-project/cmd/web/middleware/forms"
	"github/somyaranjan99/basic-go-project/internal/condriver"
	"github/somyaranjan99/basic-go-project/internal/helpers"
	"github/somyaranjan99/basic-go-project/internal/repository"
	"github/somyaranjan99/basic-go-project/internal/repository/dbrepo"
	"github/somyaranjan99/basic-go-project/internal/reservationtypes"
	"github/somyaranjan99/basic-go-project/pkg/config"
	"github/somyaranjan99/basic-go-project/pkg/models"
	"github/somyaranjan99/basic-go-project/pkg/render"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Repository struct {
	Repo *config.AppConfig
	Db   repository.DatabaseRepo
}

func NewRepo(a *config.AppConfig, db *condriver.DB) *Repository {
	return &Repository{Repo: a, Db: dbrepo.NewRepositoryDBHandler(a, db.SQL)}
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

	render.RenderTemplate(w, r, "generals.page.tmpl", repo.Repo, &models.TemplateData{})
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
	err := r.ParseForm()
	if err != nil {
		repo.Repo.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	layout := "2006-01-02"

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")
	start_date, err := time.Parse(layout, sd)

	if err != nil {
		repo.Repo.Session.Put(r.Context(), "error", "can't parse start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	end_date, err := time.Parse(layout, ed)
	if err != nil {
		repo.Repo.Session.Put(r.Context(), "error", "can't parse start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	rooms, err := repo.Db.SearchAvailabilityForAllRooms(start_date, end_date)
	if len(rooms) == 0 {
		repo.Repo.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}
	data := make(map[string]interface{})
	data["rooms"] = rooms
	res := models.Reservation{
		StartDate: start_date,
		EndDate:   end_date,
	}
	repo.Repo.Session.Put(r.Context(), "reservation", res)

	render.RenderTemplate(w, r, "choose-room.page.tmpl", repo.Repo, &models.TemplateData{
		Data: data,
	})

	//w.Write([]byte(fmt.Sprintf("start date is %s end date is %s", start, end)))
}

func (repo *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	res, _ := repo.Repo.Session.Get(r.Context(), "reservation").(models.Reservation)

	// if !ok {
	// 	repo.Repo.Session.Put(r.Context(), "error", "can't get reservation from session")
	// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	// 	return
	// }
	room, err := repo.Db.GetRoomByID(res.RoomID)
	if err != nil {
		repo.Repo.Session.Put(r.Context(), "error", "can't get reservation from db")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	res.Room.RoomName = room.RoomName
	repo.Repo.Session.Put(r.Context(), "reservation", res)
	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed
	data := make(map[string]interface{})
	data["reservation"] = res
	render.RenderTemplate(w, r, "make-reservation.page.tmpl", repo.Repo, &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

func (repo *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	res, ok := repo.Repo.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		repo.Repo.Session.Put(r.Context(), "error", "can't parse session data")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		log.Println(err)
		return
	}
	//now := time.Now()
	reservation := &models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: res.StartDate,
		EndDate:   res.EndDate,
		Room:      reservationtypes.Room{RoomName: res.Room.RoomName},
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
	layout := "2006-01-02"
	// startDate := now.Truncate(24 * time.Hour)
	// endDate := now.AddDate(0, 0, 1).Truncate(24 * time.Hour)
	startDate, err := time.Parse(layout, r.Form.Get("start_date"))
	if err != nil {
		repo.Repo.Session.Put(r.Context(), "error", "can't parse start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	endDate, err := time.Parse(layout, r.Form.Get("end_date"))
	if err != nil {
		repo.Repo.Session.Put(r.Context(), "error", "can't parse end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	room_id, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		repo.Repo.Session.Put(r.Context(), "error", "can't parse room id")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	postReservation := reservationtypes.Reservation{
		FirstName: reservation.FirstName,
		LastName:  reservation.LastName,
		Email:     reservation.Email,
		Phone:     reservation.Phone,
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    room_id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	response, err := repo.Db.BookReservation(&postReservation)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	fmt.Print(response)
	//repo.(&postReservation)
	repo.Repo.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}
func (repo *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := repo.Repo.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		repo.Repo.ErrorLog.Println("can't get item from session")
		log.Println("can't get item from session")
		repo.Repo.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Println(reservation)
	repo.Repo.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation
	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed
	render.RenderTemplate(w, r, "reservation-summary.page.tmpl", repo.Repo, &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})

}
func (repo *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	splittedString := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(splittedString[2])
	if err != nil {
		repo.Repo.Session.Put(r.Context(), "error", "Missing url parameters")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	res, ok := repo.Repo.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		repo.Repo.Session.Put(r.Context(), "error", "Can't get reservation from session")
		fmt.Print("testing")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	res.RoomID = id
	repo.Repo.Session.Put(r.Context(), "reservation", res)
	// res, _ = repo.Repo.Session.Get(r.Context(), "reservation").(models.Reservation)
	// fmt.Println(res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}
func (repo *Repository) Login(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	data["user_login"] = data

	render.RenderTemplate(w, r, "login.page.tmpl", repo.Repo, &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}
func (repo *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("post login data")
	//render.RenderTemplate(w, r, "login.page.tmpl", repo.Repo, &models.TemplateData{})
	err := r.ParseForm()
	if err != nil {
		repo.Repo.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	user := reservationtypes.User{
		Email:    r.Form.Get("email"),
		Password: r.Form.Get("password"),
	}
	form := forms.New(&r.PostForm)
	form.Has("email", r)
	form.Has("password", r)
	form.IsValidEmail("email")
	form.MinLength("password", 6, r)
	// fmt.Println(form.Errors)
	if !form.Valid() {
		data := make(map[string]interface{})
		data["user_login"] = user
		render.RenderTemplate(w, r, "login.page.tmpl", repo.Repo, &models.TemplateData{
			Data: data,
			Form: form,
		})
		return
	}
	isAuth, err := repo.Db.IsAuthenticatedUser(user.Email, user.Password)
	if err != nil {
		repo.Repo.Session.Put(r.Context(), "error", "Internal server Error")
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}
	if !isAuth {
		repo.Repo.Session.Put(r.Context(), "error", "Internal server Error")
		http.Redirect(w, r, "/", http.StatusInternalServerError)
		return
	}

	repo.Repo.Session.Put(r.Context(), "success", "Loged in success")
	repo.Repo.Session.Put(r.Context(), "user_login", isAuth)
	render.RenderTemplate(w, r, "/", repo.Repo, &models.TemplateData{})
}
func (repo *Repository) UserSignup(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	data["user_signup"] = reservationtypes.User{}
	render.RenderTemplate(w, r, "/signup.page.tmpl", repo.Repo, &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}
func (repo *Repository) PostUserSignup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		repo.Repo.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	email := r.Form.Get("email")
	res, _ := repo.Db.GetUserByEmail(email)
	// if err != nil {
	// 	fmt.Println(err)
	// 	repo.Repo.Session.Put(r.Context(), "error", "Internal server error")
	// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	// 	return
	// }
	if res == true {
		repo.Repo.Session.Put(r.Context(), "error", "User already exist")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Println(res, "GetUserByEmail")
	form := forms.New(&r.PostForm)
	form.Has("first_name", r)
	form.Has("last_name", r)
	form.Has("email", r)
	form.Has("password", r)
	form.IsValidEmail("email")
	form.MinLength("password", 6, r)
	user := reservationtypes.User{}
	if !form.Valid() {

		data := make(map[string]interface{})
		data["user_signup"] = user
		repo.Repo.Session.Put(r.Context(), "error", "Please enter all the fields")
		render.RenderTemplate(w, r, "/user/signup", repo.Repo, &models.TemplateData{
			Data: data,
			Form: form,
		})
		return
	}
	password, err := bcrypt.GenerateFromPassword([]byte(r.Form.Get("password")), 12)

	if err != nil {
		fmt.Println(err)
		data := make(map[string]interface{})
		data["user_signup"] = user
		repo.Repo.Session.Put(r.Context(), "error", "internal server error")
		render.RenderTemplate(w, r, "/user/signup", repo.Repo, &models.TemplateData{
			Data: data,
			Form: form,
		})
		return
	}
	access_level, _ := strconv.Atoi(r.Form.Get("access_level"))
	fmt.Println(access_level, "access_level")

	user = reservationtypes.User{
		FirstName:   r.Form.Get("first_name"),
		LastName:    r.Form.Get("last_name"),
		Email:       email,
		Password:    string(password),
		AccessLevel: access_level,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	res, err = repo.Db.SignupUser(&user)
	fmt.Println(err, "SignupUser")

	if err != nil {
		data := make(map[string]interface{})
		data["user_signup"] = user
		repo.Repo.Session.Put(r.Context(), "error", "internal server error")
		render.RenderTemplate(w, r, "/user/signup", repo.Repo, &models.TemplateData{
			Data: data,
			Form: form,
		})
		return
	}
	if !res {
		data := make(map[string]interface{})
		data["user_signup"] = user
		repo.Repo.Session.Put(r.Context(), "error", "internal server error")
		render.RenderTemplate(w, r, "/user/signup", repo.Repo, &models.TemplateData{
			Data: data,
			Form: form,
		})
		return
	}

	repo.Repo.Session.Put(r.Context(), "success", "User sign up success")
	http.Redirect(w, r, "/", http.StatusCreated)
	return

}
