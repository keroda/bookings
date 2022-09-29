package handlers

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/keroda/bookings/internal/models"
)

// type postData struct {
// 	key   string
// 	value string
// }

type theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}

func TestHandlers(t *testing.T) {
	var myTest = theTests{
		{"home", "/", "GET", http.StatusOK},
		{"about", "/about", "GET", http.StatusOK},
		{"gq", "/generals-quarters", "GET", http.StatusOK},
		{"ms", "/majors-suite", "GET", http.StatusOK},
		//for max coverage we would need to
		//test all the get-requests

		//we will test post requests separately
		// {"post-search-avail", "/search-availability", "POST", []postData{
		// 	{key: "start", value: "2022-01-01"},
		// 	{key: "end", value: "2022-01-02"},
		// }, http.StatusOK},
		// {"post-reservation", "/make-reservation-json", "POST", []postData{
		// 	{key: "fname", value: "Bo"},
		// 	{key: "lname", value: "Yo"},
		// 	{key: "email", value: "bo@yo.no"},
		// 	{key: "phone", value: "01234"},
		// }, http.StatusOK},
	}

	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range myTest {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}
	}
}

func TestRepositoryPostReservation(t *testing.T) {

	//build the for's body
	// reqBody := "start_date=2030-01-01"
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2030-01-01")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Johnny")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=January")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "email=jo@ja.com")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=54321012345")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	//easier way to build the form body:
	postedData := url.Values{}
	postedData.Add("start_date", "2030-01-01")
	postedData.Add("end_date", "2030-01-02")
	postedData.Add("first_name", "Joe")
	postedData.Add("last_name", "Theman")
	postedData.Add("email", "joe@hey.oh")
	postedData.Add("phone", "555-444-3332")
	postedData.Add("room_id", "1")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-url-encoded")

	rr := httptest.NewRecorder()
	h := http.HandlerFunc(Repo.PostReservation)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post-reservation handler status error: expected %d but got %d", http.StatusSeeOther, rr.Code)
	}

	//test the same call w missing form body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-url-encoded")
	rr = httptest.NewRecorder()
	h = http.HandlerFunc(Repo.PostReservation)
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post-reservation handler status error: expected %d but got %d", http.StatusTemporaryRedirect, rr.Code)
	}

	//we could/should test the other if-tests (copy-paste)
	//start_date
	//end_date
	//roomID
	//form validation
	//add reservation
	//add restriction
}

func TestRepositoryReservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
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

	h := http.HandlerFunc(Repo.Reservation)
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler status error: expected %d but got %d", http.StatusOK, rr.Code)
	}

	//test case for reservation not in session (reset/re-initialize everything)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler status error: expected %d but got %d", http.StatusTemporaryRedirect, rr.Code)
	}

	//test non-existing room (id > 2)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler status error: expected %d but got %d", http.StatusTemporaryRedirect, rr.Code)
	}
}

func TestAvailabilityJSON(t *testing.T) {
	//test cases:

	//rooms NOT available
	//rooms available

	// //read data that is returned as JSON
	// var j jsonResponse
	// err := json.Unmarshal([]byte(rr.Body.String()), &j)
	// if err != nil {
	// 	t.Error("failed to parse JSON")
	// }
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
