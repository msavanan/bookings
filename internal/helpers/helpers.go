package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/msavanan/bookings/internal/config"
)

var app *config.AppConfig

func NewHelper(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, statusCode int) {
	app.InfoLog.Println("Client error with status of ", statusCode)
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func IsAuthenticated(r *http.Request) bool {
	exists := app.Session.Exists(r.Context(), "user_id")
	return exists
}
