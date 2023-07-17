package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/msavanan/bookings/internal/config"
	"github.com/msavanan/bookings/internal/driver"
	"github.com/msavanan/bookings/internal/forms"
	"github.com/msavanan/bookings/internal/helpers"
	"github.com/msavanan/bookings/internal/models"
	"github.com/msavanan/bookings/internal/render"
	"github.com/msavanan/bookings/internal/repository"
	"github.com/msavanan/bookings/internal/repository/dbrepo"
)

var Repo *Repository

// Repository pattern
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// For Testing only
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestPostgresRepo(a),
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})

}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})

}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// GET service Availablity
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// POST service Availablity
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if len(rooms) == 0 {
		// No Availabilty
		m.App.InfoLog.Println("No Availability")
		m.App.Session.Put(r.Context(), "error", "No Availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
	}

	var data = make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose_room.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})

	for _, room := range rooms {
		m.App.InfoLog.Println("ROOM: ", room.ID, room.RoomName)
	}

	//w.Write([]byte(fmt.Sprintf("Posted to search Availability \n start date: %s, end date: %s", start, end)))
}

type JsonResponse struct {
	Ok        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomId    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// service Availablity-json
func (m *Repository) AvailabilityJson(w http.ResponseWriter, r *http.Request) {
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	roomId, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
	}

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
	}

	available, err := m.DB.SearchAvailabilityByDateByRoomId(startDate, endDate, roomId)
	if err != nil {
		helpers.ServerError(w, err)
	}

	log.Println("available.............", available)

	resp := JsonResponse{
		Ok: available,
		//Message: "",
		StartDate: sd,
		EndDate:   ed,
		RoomId:    strconv.Itoa(roomId),
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
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		// helpers.ServerError(w, errors.New("type assertion to models.reservation failed"))
		// return
		m.App.Session.Put(r.Context(), "error", "Can't get reservation form session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := m.DB.GetRoomById(res.RoomId)
	if err != nil {
		// helpers.ServerError(w, err)
		// return

		m.App.Session.Put(r.Context(), "error", "Can't find room")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res
	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// Post Reservation
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// helpers.ServerError(w, err)
		// return
		m.App.Session.Put(r.Context(), "error", "Can't parse form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		//helpers.ServerError(w, errors.New("type assertion to string failed"))
		m.App.Session.Put(r.Context(), "error", "Can't find reservation in session data")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fields := [4]string{"first_name", "last_name", "email", "phone"}

	reservation.FirstName = r.Form.Get(fields[0])
	reservation.LastName = r.Form.Get(fields[1])
	reservation.Email = r.Form.Get(fields[2])
	reservation.Phone = r.Form.Get(fields[3])

	form := forms.New(r.PostForm)

	form.Required(fields[:3]...)
	form.MinimumLength(fields[0], 3)
	form.IsEmail(fields[2])

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return

	}

	newReservationId, err := m.DB.InsertReservations(reservation)

	if err != nil {
		// helpers.ServerError(w, err)
		// return
		m.App.Session.Put(r.Context(), "error", "Can't find reservation id")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	restrictions := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomId:        reservation.RoomId,
		ReservationId: newReservationId,
		RestrictionId: 1,
	}

	err = m.DB.InsertRoomRestrictions(restrictions)
	if err != nil {
		// helpers.ServerError(w, err)
		// return
		m.App.Session.Put(r.Context(), "error", "Can't find restriction id")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Send Mail - First to the user
	htmlMessage := fmt.Sprintf(`
	<strong>Reservation Confirmation</strong><br>
	Dear %s, <br>
	This confirm the reservation for %s rooms from %s to %s
	`, reservation.FirstName, reservation.Room.RoomName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))

	msg := models.MailData{
		To:       reservation.Email,
		From:     "me@here.com",
		Subject:  "Reservation Confirmation",
		Content:  htmlMessage,
		Template: "basic.html",
	}

	m.App.MailChan <- msg

	//Send Mail - Property  Owner
	htmlMessage = fmt.Sprintf(`
	<strong>Reservation Notification</strong><br>
	Details of reservation <br>
	Room Name     : %s <br> 
	Arrival Date  : %s <br> 
	Departure Date: %s <br>
	`, reservation.Room.RoomName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))

	msg = models.MailData{
		To:       "me@here.com",
		From:     "me@here.com",
		Subject:  "Reservation Notification",
		Content:  htmlMessage,
		Template: "basic.html",
	}

	m.App.MailChan <- msg

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

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})

}

func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {

	roomId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("type assertion to models.reservation failed"))
		return
	}

	res.RoomId = roomId

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}

func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		helpers.ServerError(w, err)
	}
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	var res models.Reservation

	res.RoomId = id

	room, err := m.DB.GetRoomById(res.RoomId)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.Room.RoomName = room.RoomName

	res.StartDate, err = time.Parse("2006-01-02", sd)
	if err != nil {
		helpers.ServerError(w, err)
	}

	res.EndDate, err = time.Parse("2006-01-02", ed)
	if err != nil {
		helpers.ServerError(w, err)
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}

func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})

}

func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	//Sesion fixation attack
	err := m.App.Session.RenewToken(r.Context())
	if err != nil {
		helpers.ServerError(w, err)
	}

	err = r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	form := forms.New(r.PostForm)
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form.Required("email", "password")
	form.IsEmail("email")

	if !form.Valid() {
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (m Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{})
}

func (m *Repository) AdminReservationsAll(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.Template(w, r, "admin-reservations-all.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) AdminReservationsNew(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllNewReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations
	render.Template(w, r, "admin-reservations-new.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) AdminShowReservation(w http.ResponseWriter, r *http.Request) {
	explode := strings.Split(r.RequestURI, "/")
	src := explode[3]
	id, err := strconv.Atoi(explode[4])
	if err != nil {
		log.Println(err)
	}

	stringMap := make(map[string]string)
	stringMap["src"] = src

	reservation, err := m.DB.GetReservationById(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, r, "admin-reservations-show.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
		Form:      forms.New(nil),
	})
}

func (m *Repository) AdminPostShowReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	explode := strings.Split(r.RequestURI, "/")
	src := explode[3]
	id, err := strconv.Atoi(explode[4])
	if err != nil {
		log.Println(err)
	}

	var reservation models.Reservation

	reservation.Id = id
	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	err = m.DB.UpdateReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Changes saved")

	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
}

func (m *Repository) AdminReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	if r.URL.Query().Get("y") != "" {
		year, err := strconv.Atoi(r.URL.Query().Get("y"))
		if err != nil {
			log.Println("Failed to query year")
		}
		month, err := strconv.Atoi(r.URL.Query().Get("m"))
		if err != nil {
			log.Println("Failed to query month")
		}

		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}

	next := now.AddDate(0, 1, 0)
	last := now.AddDate(0, -1, 0)

	nextMonth := next.Format("01")
	nextMonthYear := next.Format("2006")

	lastMonth := last.Format("01")
	lastMonthYear := last.Format("2006")

	data := make(map[string]interface{})
	data["now"] = now

	stringMap := make(map[string]string)
	stringMap["next_month"] = nextMonth
	stringMap["last_month"] = lastMonth
	stringMap["next_year"] = nextMonthYear
	stringMap["last_year"] = lastMonthYear

	stringMap["this_month"] = now.Format("01")
	stringMap["this_month_year"] = now.Format("2006")

	render.Template(w, r, "admin-reservations-calendar.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
	})
}

func (m *Repository) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")

	err := m.DB.UpdateProcessedForReservation(id, 1)
	if err != nil {
		log.Println(err)
	}

	m.App.Session.Put(r.Context(), "flash", "Marked as processed")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)

}

func (m *Repository) AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")

	err := m.DB.DeleteReservation(id)
	if err != nil {
		log.Println(err)
	}

	m.App.Session.Put(r.Context(), "flash", "Reservation Deleted")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)

}
