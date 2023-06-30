package handlers

import (
	"net/http"

	"github.com/msavanan/bookings/pkg/config"
	"github.com/msavanan/bookings/pkg/models"
	"github.com/msavanan/bookings/pkg/render"
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
	remoteIp := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIp)

	render.RenderTemplate(w, "home.page.html", &models.TemplateData{})

}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello again"

	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIp

	td := models.TemplateData{
		StringMap: stringMap,
	}
	render.RenderTemplate(w, "about.page.html", &td)

}
