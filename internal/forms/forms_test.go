package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("post", "/whatever", nil)
	form := New(r.Form)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("post", "/whatever", nil)
	form := New(r.Form)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required field is missing")

	}

	postData := url.Values{}
	postData.Add("a", "a")
	postData.Add("b", "b")
	postData.Add("c", "c")

	r = httptest.NewRequest("post", "/whatever", nil)
	r.PostForm = postData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("form shows invalid when the form is valid")
	}
}

func TestForm_MinimumLength(t *testing.T) {
	r := httptest.NewRequest("post", "/whatever", nil)

	postData := url.Values{}
	postData.Add("a", "a")
	//postData.Add("b", "b")
	//postData.Add("c", "c")
	r.PostForm = postData
	form := New(r.PostForm)

	form.MinimumLength("a", 3)
	if form.Valid() {
		t.Error("Expected Invalid length err but got No invalid length err")
	}

	isError := form.Errors.Get("a")
	if isError == "" {
		t.Error("should not have an error, but got an error")
	}

	postData = url.Values{}
	postData.Add("a", "abc")
	//postData.Add("b", "b")
	//postData.Add("c", "c")
	r.PostForm = postData
	form = New(r.PostForm)

	form.MinimumLength("a", 3)
	if !form.Valid() {
		t.Error("Expected Invalid length err but got No invalid length err")
	}

	isError = form.Errors.Get("a")
	if isError != "" {
		t.Error("should have an error, but no error")
	}

}

func TestForm_IsEmail(t *testing.T) {

	postData := url.Values{}
	postData.Add("a", "a@her.com")

	form := New(postData)

	form.IsEmail("a")

	if !form.Valid() {
		t.Error("Expected valid email err but got invalid email err")
	}

	postData = url.Values{}
	postData.Add("a", "a.com")

	form = New(postData)

	form.IsEmail("a")

	if form.Valid() {
		t.Error("Expected valid email err but got invalid email err")
	}

}

func TestForm_Has(t *testing.T) {

	postData := url.Values{}
	postData.Add("a", "a")
	postData.Add("b", "b")
	postData.Add("c", "c")

	form := New(postData)

	a := form.Has("a")

	if !a {
		t.Error("Expected valid but got err")
	}

	postData = url.Values{}
	postData.Add("a", "")

	form = New(postData)

	a = form.Has("a")

	if a {
		t.Error("Expected err but got valid ")
	}

}
