package handlers

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTest = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"gq", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"major", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"sa", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"mr", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"post-sa", "/search-availability", "POST", []postData{
		{key: "start", value: "2023-11-09"},
		{key: "end", value: "2023-11-09"},
	}, http.StatusOK},
	{"post-sa-json", "/search-availability-json", "POST", []postData{
		{key: "start", value: "2023-11-09"},
		{key: "end", value: "2023-11-09"},
	}, http.StatusOK},
	{"post-sa-json", "/make-reservation", "POST", []postData{
		{key: "first_name", value: "smith"},
		{key: "last_name", value: "John"},
		{key: "email", value: "me@here.com"},
		{key: "phone", value: "5555-5555-5555"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {

	routes := getRoutes()

	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	app.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	for _, e := range theTest {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal()
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d, but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {
			value := url.Values{}

			for _, x := range e.params {
				value.Add(x.key, x.value)
			}

			resp, err := ts.Client().PostForm(ts.URL+e.url, value)
			if err != nil {
				t.Log(err)
				t.Fatal()
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d, but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}

		}
	}

}
