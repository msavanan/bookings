package render

import (
	"net/http"
	"testing"

	"github.com/msavanan/bookings/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSessionData()
	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "123")
	result := AddDefaultData(&td, r)

	if result.Flash != "123" {
		t.Errorf("Failed, expected flash value is  %s but available value %s", "123", result.Flash)
	}

}

func getSessionData() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		//log.Fatal(err)
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil

}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	var ww myWriter
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Errorf("Failed to create Template cache %e", err)
	}

	r, err := getSessionData()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	err = RenderTemplate(&ww, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Errorf("Failed to render the template %e", err)
	}

	err = RenderTemplate(&ww, r, "non-existence.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Errorf("rendered template that doesn't exists %e", err)
	}
}

func TestNewTemplates(t *testing.T) {
	NewTemplate(app)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

}
