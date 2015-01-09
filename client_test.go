package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestAddsGivenParameters(t *testing.T) {

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, fmt.Sprintf("%s", r.URL))
	}))

	parameters := map[string]string{
		"custom": "param",
	}

	client := New(testServer.URL, "mykey", "private")
	response, _ := client.Get("/trello/api", parameters)

	u, _ := url.Parse(string(response))
	q := u.Query()

	if strings.TrimSpace(q.Get("custom")) != "param" {
		t.Error("Should have included the parameter")
	}
}

func TestAddNoParameters(t *testing.T) {

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, fmt.Sprintf("%s", r.URL))
	}))

	client := New(testServer.URL, "mykey", "private")
	response, _ := client.Get("/trello/api")

	u, _ := url.Parse(string(response))

	var expected = "/trello/api?key=mykey&token=private"
	if strings.TrimSpace(u.String()) != expected {
		t.Errorf("Expected '%s' got '%s'", expected, strings.TrimSpace(u.String()))
	}
}

func TestReturnsErrorOn500(t *testing.T) {

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Oops", http.StatusInternalServerError)
	}))

	client := New(testServer.URL, "user", "pass")
	_, err := client.Get("/api/agents")

	var expected = "Response code 500"
	if err.Error() != expected {
		t.Errorf("Expected '%s' got '%s'", expected, err.Error())
	}
}

func TestReturnsErrorOn404(t *testing.T) {

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Oops", http.StatusNotFound)
	}))

	client := New(testServer.URL, "user", "pass")
	_, err := client.Get("/api/agents")

	var expected = "Response code 404"
	if err.Error() != expected {
		t.Errorf("Expected '%s' got '%s'", expected, err.Error())
	}
}
