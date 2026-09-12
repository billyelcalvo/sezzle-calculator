package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/diegomora/sezzle-calculator/internal/handlers"
)

func TestCORSMiddleware(tester *testing.T) {
	origins := []struct {
		name    string
		origin  string
		allowed bool
	}{
		{name: "allowed origin", origin: "http://localhost:5173", allowed: true},
		{name: "different port", origin: "http://localhost:5174"},
		{name: "different host", origin: "http://127.0.0.1:5173"},
		{name: "untrusted origin", origin: "https://example.com"},
		{name: "origin prefix", origin: "http://localhost:5173.example.com"},
		{name: "null origin", origin: "null"},
		{name: "no origin"},
	}

	for _, origin := range origins {
		for _, method := range []string{http.MethodPost, http.MethodOptions} {
			tester.Run(origin.name+"/"+method, func(tester *testing.T) {
				called := false
				next := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					called = true
					writer.WriteHeader(http.StatusAccepted)
				})
				request := httptest.NewRequest(method, "/api/calculate", nil)
				if origin.origin != "" {
					request.Header.Set("Origin", origin.origin)
				}
				if method == http.MethodOptions {
					request.Header.Set("Access-Control-Request-Method", http.MethodPost)
					request.Header.Set("Access-Control-Request-Headers", "content-type")
				}
				recorder := httptest.NewRecorder()

				corsMiddleware(next).ServeHTTP(recorder, request)

				if vary := recorder.Header().Get("Vary"); vary != "Origin" {
					tester.Errorf("Vary = %q; want Origin", vary)
				}
				for header, allowedValue := range map[string]string{
					"Access-Control-Allow-Origin":  "http://localhost:5173",
					"Access-Control-Allow-Methods": "POST, OPTIONS",
					"Access-Control-Allow-Headers": "Content-Type",
				} {
					expected := ""
					if origin.allowed {
						expected = allowedValue
					}
					if actual := recorder.Header().Get(header); actual != expected {
						tester.Errorf("%s = %q; want %q", header, actual, expected)
					}
				}

				if method == http.MethodOptions {
					if recorder.Code != http.StatusNoContent {
						tester.Errorf("status = %d; want 204", recorder.Code)
					}
					if called {
						tester.Error("preflight must not call the next handler")
					}
					if recorder.Body.Len() != 0 {
						tester.Errorf("preflight body must be empty; got %q", recorder.Body.String())
					}
				} else {
					if !called {
						tester.Error("request must reach the next handler")
					}
					if recorder.Code != http.StatusAccepted {
						tester.Errorf("status = %d; want 202", recorder.Code)
					}
				}
			})
		}
	}
}

func TestCORSRouter(tester *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/calculate", handlers.Calculate)
	handler := corsMiddleware(mux)
	tests := []struct {
		name         string
		method       string
		path         string
		body         string
		status       int
		expectedBody string
	}{
		{name: "success", method: http.MethodPost, path: "/api/calculate", body: `{"operation":"add","a":10,"b":5}`, status: http.StatusOK, expectedBody: `{"result":15}`},
		{name: "invalid JSON", method: http.MethodPost, path: "/api/calculate", body: `{`, status: http.StatusBadRequest, expectedBody: `{"error":"body must be a valid JSON object with valid fields and types"}`},
		{name: "division by zero", method: http.MethodPost, path: "/api/calculate", body: `{"operation":"divide","a":10,"b":0}`, status: http.StatusBadRequest, expectedBody: `{"error":"division by zero"}`},
		{name: "unsupported method", method: http.MethodGet, path: "/api/calculate", status: http.StatusMethodNotAllowed, expectedBody: `{"error":"method must be POST"}`},
		{name: "unknown route", method: http.MethodPost, path: "/missing", status: http.StatusNotFound, expectedBody: "404 page not found"},
		{name: "preflight", method: http.MethodOptions, path: "/api/calculate", status: http.StatusNoContent},
	}

	for _, test := range tests {
		tester.Run(test.name, func(tester *testing.T) {
			request := httptest.NewRequest(test.method, test.path, strings.NewReader(test.body))
			request.Header.Set("Origin", "http://localhost:5173")
			if test.method == http.MethodOptions {
				request.Header.Set("Access-Control-Request-Method", http.MethodPost)
				request.Header.Set("Access-Control-Request-Headers", "content-type")
			} else {
				request.Header.Set("Content-Type", "application/json")
			}
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, request)

			if recorder.Code != test.status {
				tester.Errorf("status = %d; want %d", recorder.Code, test.status)
			}
			if origin := recorder.Header().Get("Access-Control-Allow-Origin"); origin != "http://localhost:5173" {
				tester.Errorf("Access-Control-Allow-Origin = %q; want http://localhost:5173", origin)
			}
			if body := strings.TrimSpace(recorder.Body.String()); body != test.expectedBody {
				tester.Errorf("body = %q; want %q", body, test.expectedBody)
			}
		})
	}
}
