package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/msavanan/bookings/internal/config"
	"github.com/msavanan/bookings/internal/forms"
	"github.com/msavanan/bookings/internal/helpers"
	"github.com/msavanan/bookings/internal/models"
	"github.com/msavanan/bookings/internal/render"
)

var Repo *Repository

// Repository pattern
type Repository struct {
	App *config.AppConfig
}

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})

}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{})

}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// GET service Availablity
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// POST service Availablity
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	w.Write([]byte(fmt.Sprintf("Posted to search Availability \n start date: %s, end date: %s", start, end)))
}

type JsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

// service Availablity-json
func (m *Repository) AvailabilityJson(w http.ResponseWriter, r *http.Request) {
	resp := JsonResponse{
		Ok:      true,
		Message: "Available",
	}
	if out, err := json.MarshalIndent(resp, "", ""); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	} else {
		helpers.ServerError(w, err)
	}

}

// Get Reservation
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	data["reservation"] = models.Reservation{}
	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// Post Reservation
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	//err := errors.New("checkeing helper functions")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	fields := [4]string{"first_name", "last_name", "email", "phone"}

	reservation := models.Reservation{
		FirstName: r.Form.Get(fields[0]),
		LastName:  r.Form.Get(fields[1]),
		Email:     r.Form.Get(fields[2]),
		Phone:     r.Form.Get(fields[3]),
	}

	form := forms.New(r.PostForm)

	form.Required(fields[:3]...)
	form.MinimumLength(fields[0], 3)
	form.IsEmail(fields[2])

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return

	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		log.Println("Cannot get data from session")
		m.App.ErrorLog.Println("cn't get reservation from session")
		m.App.Session.Put(r.Context(), "ttttt", "Cannot get reservation from session")
		m.App.Session.Put(r.Context(), "warning", "Cannot get reservation from session")
		//m.App.Session.Put(r.Context(), "error", "Cannot get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.RenderTemplate(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})

}
