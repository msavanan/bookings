package main

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"text/template"

	"github.com/msavanan/bookings/internal/models"
	"github.com/msavanan/bookings/internal/render"
)

var functions = template.FuncMap{
	"humanDate":  render.HumanDate,
	"formatDate": render.FormatDate,
	"iterate":    render.Iterate,
	"add":        render.Add,
}

func TestMain(m *testing.M) {
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Restriction{})
	//gob.Register(models.RoomRestriction{})
	gob.Register(models.Room{})
	gob.Register(map[string]int{})

	os.Exit(m.Run())
}

type myHandler struct{}

func (h *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
