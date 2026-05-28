package httpclient

import (
	"net/http"
	"testing"
)

func TestOption(t *testing.T) {
	t.Run("should be able to set Authorization header via RequestOption", func(t *testing.T) {
		token := "eyJheader.payload.signature"
		want := "Bearer " + token

		req, _ := http.NewRequest("GET", "http://localhost:8080", nil)
		opt := RequestOption(func(r *http.Request) {
			r.Header.Add("Authorization", "Bearer "+token)
		})

		opt(req)

		got := req.Header.Get("Authorization")
		if got != want {
			t.Errorf("Authorization header got: %s, want: %s", got, want)
		}
	})

	t.Run("should be able to set custom header", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "http://localhost:8080", nil)
		opt := HeaderOption("X-Custom", "value")
		opt(req)

		got := req.Header.Get("X-Custom")
		if got != "value" {
			t.Errorf("X-Custom header got: %s, want: %s", got, "value")
		}
	})

	t.Run("should be able to set Authorization header via BearerTokenOption", func(t *testing.T) {
		token := "eyJheader.payload.signature"
		want := "Bearer " + token

		req, _ := http.NewRequest("GET", "http://localhost:8080", nil)
		opt := BearerTokenOption(token)
		opt(req)

		got := req.Header.Get("Authorization")
		if got != want {
			t.Errorf("Authorization header got: %s, want: %s", got, want)
		}
	})

	t.Run("jsonOption should not override existing Content-Type", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "http://localhost:8080", nil)
		req.Header.Set("Content-Type", "application/custom")
		jsonOption(req)

		got := req.Header.Get("Content-Type")
		want := "application/custom"
		if got != want {
			t.Errorf("Content-Type got: %s, want: %s", got, want)
		}
	})

	t.Run("jsonOption should set default Content-Type", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "http://localhost:8080", nil)
		jsonOption(req)

		got := req.Header.Get("Content-Type")
		want := "application/json"
		if got != want {
			t.Errorf("Content-Type got: %s, want: %s", got, want)
		}
	})

	t.Run("RequestOption should run wrapped function", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "http://localhost:8080", nil)
		opt := RequestOption(func(r *http.Request) {
			r.Header.Set("X-Req", "1")
		})
		opt(req)

		got := req.Header.Get("X-Req")
		if got != "1" {
			t.Errorf("X-Req header got: %s, want: %s", got, "1")
		}
	})
}
