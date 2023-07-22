package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/msavanan/bookings/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTest = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"gq", "/generals-quarters", "GET", http.StatusOK},
	{"major", "/majors-suite", "GET", http.StatusOK},
	{"sa", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"non-existent", "/green/eggs/hams", "GET", http.StatusNotFound},

	//New Routes
	{"login", "/user/login", "GET", http.StatusOK},
	{"logout", "/user/logout", "GET", http.StatusOK},
	{"dashboard", "/admin/dashboard", "GET", http.StatusOK},
	{"new res", "/admin/reservations-new", "GET", http.StatusOK},
	{"new all", "/admin/reservations-all", "GET", http.StatusOK},
	{"show res", "/admin/reservations/new/1/show", "GET", http.StatusOK},

	// {"post-sa", "/search-availability", "POST", []postData{
	// 	{key: "start", value: "2023-11-09"},
	// 	{key: "end", value: "2023-11-09"},
	// }, http.StatusOK},
	// {"post-sa-json", "/search-availability-json", "POST", []postData{
	// 	{key: "start", value: "2023-11-09"},
	// 	{key: "end", value: "2023-11-09"},
	// }, http.StatusOK},
	// {"post-sa-json", "/make-reservation", "POST", []postData{
	// 	{key: "first_name", value: "smith"},
	// 	{key: "last_name", value: "John"},
	// 	{key: "email", value: "me@here.com"},
	// 	{key: "phone", value: "5555-5555-5555"},
	// }, http.StatusOK},
}

func TestHandlers(t *testing.T) {

	routes := getRoutes()

	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	app.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	for _, e := range theTest {

		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal()
		}

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s, expected %d, but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}

	}

}

type form struct {
	StartDate string
	EndDate   string
	FirstName string
	LastName  string
	Email     string
	Phone     string
	RoomId    string
}

func TestRepository_PostReservations(t *testing.T) {

	formVal := form{
		StartDate: "2050-01-02",
		EndDate:   "2050-01-02",
		FirstName: "john",
		LastName:  "Smith",
		Email:     "john@smith.com",
		Phone:     "5555-5555-5555",
		RoomId:    "1",
	}

	// With valid data
	postReservationTest(formVal, t)

	// without reservation data
	postReservationTest(form{}, t)

	// Invalid form data
	req, _ := http.NewRequest("POST", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got %d wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// invalid first Name
	reqBody := fmt.Sprintf("start_date=%s", formVal.StartDate)
	reqBody = fmt.Sprintf("%s&end_date=%s", reqBody, formVal.EndDate)
	reqBody = fmt.Sprintf("%s&first_name=%s", reqBody, "j")
	reqBody = fmt.Sprintf("%s&last_name=%s", reqBody, formVal.LastName)
	reqBody = fmt.Sprintf("%s&email=%s", reqBody, formVal.Email)
	reqBody = fmt.Sprintf("%s&phone=%s", reqBody, formVal.Phone)
	reqBody = fmt.Sprintf("%s&room_id=%s", reqBody, formVal.RoomId)

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	reservation := models.Reservation{
		Id:        1,
		FirstName: "dff",
		LastName:  "ffff",
		Email:     "fff@gmail.com",
		Phone:     "66666",
	}

	session.Put(ctx, "reservation", reservation)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("PostReservation handler returned wrong response code: got %d wanted %d", rr.Code, http.StatusOK)
	}

	// insert reservations
	reqBody = fmt.Sprintf("start_date=%s", formVal.StartDate)
	reqBody = fmt.Sprintf("%s&end_date=%s", reqBody, formVal.EndDate)
	reqBody = fmt.Sprintf("%s&first_name=%s", reqBody, formVal.FirstName)
	reqBody = fmt.Sprintf("%s&last_name=%s", reqBody, formVal.LastName)
	reqBody = fmt.Sprintf("%s&email=%s", reqBody, formVal.Email)
	reqBody = fmt.Sprintf("%s&phone=%s", reqBody, formVal.Phone)
	reqBody = fmt.Sprintf("%s&room_id=%s", reqBody, "1000")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	reservation = models.Reservation{
		Id:        1,
		FirstName: "dff",
		LastName:  "ffff",
		Email:     "fff@gmail.com",
		Phone:     "66666",
		RoomId:    2,
	}

	session.Put(ctx, "reservation", reservation)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got %d wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// insert restriction id
	reqBody = fmt.Sprintf("start_date=%s", formVal.StartDate)
	reqBody = fmt.Sprintf("%s&end_date=%s", reqBody, formVal.EndDate)
	reqBody = fmt.Sprintf("%s&first_name=%s", reqBody, formVal.FirstName)
	reqBody = fmt.Sprintf("%s&last_name=%s", reqBody, formVal.LastName)
	reqBody = fmt.Sprintf("%s&email=%s", reqBody, formVal.Email)
	reqBody = fmt.Sprintf("%s&phone=%s", reqBody, formVal.Phone)
	reqBody = fmt.Sprintf("%s&room_id=%s", reqBody, "100")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	reservation = models.Reservation{
		Id:        1,
		FirstName: "dff",
		LastName:  "ffff",
		Email:     "fff@gmail.com",
		Phone:     "66666",
		RoomId:    100,
	}

	session.Put(ctx, "reservation", reservation)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got %d wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func postReservationTest(formVal form, t *testing.T) {
	reqBody := fmt.Sprintf("start_date=%s", formVal.StartDate)
	reqBody = fmt.Sprintf("%s&end_date=%s", reqBody, formVal.EndDate)
	reqBody = fmt.Sprintf("%s&first_name=%s", reqBody, formVal.FirstName)
	reqBody = fmt.Sprintf("%s&last_name=%s", reqBody, formVal.LastName)
	reqBody = fmt.Sprintf("%s&email=%s", reqBody, formVal.Email)
	reqBody = fmt.Sprintf("%s&phone=%s", reqBody, formVal.Phone)
	reqBody = fmt.Sprintf("%s&room_id=%s", reqBody, formVal.RoomId)

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got %d wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRpository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomId: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d wanted %d", rr.Code, http.StatusOK)
	}

	//Test case were reservation is not in the session (reset everything)

	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d wanted %d", rr.Code, http.StatusOK)
	}

	// test with non-esistent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()
	reservation.RoomId = 100
	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d wanted %d", rr.Code, http.StatusOK)
	}

}

var loginTest = []struct {
	name               string
	email              string
	expectedStatusCode int
	expectedHtml       string
	expectedLocation   string
}{
	{
		"valid-credentials",
		"me@here.ca",
		http.StatusSeeOther,
		"",
		"/",
	},
	{
		"invalid-credentials",
		"jack@nimble.com",
		http.StatusSeeOther,
		"",
		"/user/login",
	},

	{
		"invalid-credentials",
		"j",
		http.StatusOK,
		`action="user/login"`,
		"",
	},
}

func TestLogin(t *testing.T) {
	for _, e := range loginTest {
		postedData := url.Values{}
		postedData.Add("email", e.email)
		postedData.Add("password", "password")

		req, _ := http.NewRequest("POST", "/user/login", strings.NewReader(postedData.Encode()))
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.PostShowLogin)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("failed %s expected code %d but code %d", e.name, e.expectedStatusCode, rr.Code)
		}

		if e.expectedLocation != "" {
			actualLocation, _ := rr.Result().Location()

			if actualLocation.String() != e.expectedLocation {
				t.Errorf("failed %s expected location %s but location %s", e.name, e.expectedLocation, actualLocation.String())
			}
		}

		if e.expectedLocation != "" {
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHtml) {
				t.Errorf("failed %s expected html %s but got html %s", e.name, e.expectedHtml, html)
			}
		}

	}

}

func getCtx(req *http.Request) context.Context {

	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}

	return ctx

}
