package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got valid when should not have been")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "A")
	postedData.Add("b", "B")
	postedData.Add("c", "C")

	r, _ = http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows does not have required fileds when it does")
	}
}

func TestForm_Has(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	has := form.Has("whatever")
	if has {
		t.Error("should be false (field is empty)")
	}

	postedData.Add("a", "A")
	form = New(postedData)
	has = form.Has("a")
	if !has {
		t.Error("should be valid (field exists and is not empty)")
	}
}

func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	min := form.MinLength("x", 2)
	if min {
		t.Error("should be false (field is empty)")
	}

	//test errors.go too
	isError := form.Errors.Get("x")
	if isError == "" {
		t.Error("should be an error there now")
	}

	postedData.Add("a", "Ab")
	form = New(postedData)
	min = form.MinLength("a", 2)
	if !min {
		t.Error("should be valid (field exists and is not empty)")
	}

	//test errors.go too
	isError = form.Errors.Get("a")
	if isError != "" {
		t.Error("should not be any errors now")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.IsEmail("x")
	if form.Valid() {
		t.Error("should be invalid (filed does not exist)")
	}

	postedData.Add("email", "aaa@")
	form.IsEmail("email")
	if form.Valid() {
		t.Error("should be invalid")
	}

	postedData = url.Values{}
	postedData.Add("email", "aaa@bbb.com")
	form = New(postedData)
	form.IsEmail("email")
	if !form.Valid() {
		t.Error("should be valid")
	}

}
