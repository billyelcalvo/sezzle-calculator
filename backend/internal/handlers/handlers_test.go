package handlers_test

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/diegomora/sezzle-calculator/internal/handlers"
)

func TestCalculateSuccess(tester *testing.T) {
	tests := []struct {
		name     string
		body     string
		expected float64
	}{
		{name: "addition", body: `{"operation":"add","a":10,"b":5}`, expected: 15},
		{name: "subtraction", body: `{"operation":"subtract","a":5,"b":10}`, expected: -5},
		{name: "multiplication", body: `{"operation":"multiply","a":-10,"b":5}`, expected: -50},
		{name: "division", body: `{"operation":"divide","a":5,"b":2}`, expected: 2.5},
		{name: "decimals", body: `{"operation":"add","a":0.1,"b":0.2}`, expected: 0.3},
		{name: "zero operands", body: `{"operation":"add","a":0,"b":0}`, expected: 0},
		{name: "zero dividend", body: `{"operation":"divide","a":0,"b":5}`, expected: 0},
		{name: "trailing whitespace", body: "{\"operation\":\"add\",\"a\":10,\"b\":5}\n\t ", expected: 15},
	}

	for _, test := range tests {
		tester.Run(test.name, func(tester *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/api/calculate", strings.NewReader(test.body))
			request.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			handlers.Calculate(recorder, request)

			payload := checkJSONResponse(tester, recorder, http.StatusOK)
			var result float64
			if err := json.Unmarshal(payload["result"], &result); err != nil {
				tester.Fatalf("decode result: %v", err)
			}
			if string(payload["result"]) == "null" {
				tester.Fatal("result must be a number, got null")
			}
			if math.IsNaN(result) || math.Abs(result-test.expected) > 1e-9 {
				tester.Errorf("result = %v; want %v", result, test.expected)
			}
		})
	}
}

func TestCalculateInvalidInput(tester *testing.T) {
	tests := []struct {
		name   string
		body   string
		status int
	}{
		{name: "empty body", body: "", status: http.StatusBadRequest},
		{name: "malformed JSON", body: `{"operation":"add","a":10,"b":`, status: http.StatusBadRequest},
		{name: "array body", body: `[]`, status: http.StatusBadRequest},
		{name: "null body", body: `null`, status: http.StatusBadRequest},
		{name: "string body", body: `"add"`, status: http.StatusBadRequest},
		{name: "empty object", body: `{}`, status: http.StatusBadRequest},
		{name: "unknown field", body: `{"operation":"add","a":10,"b":5,"extra":1}`, status: http.StatusBadRequest},
		{name: "multiple objects", body: `{"operation":"add","a":10,"b":5} {}`, status: http.StatusBadRequest},
		{name: "trailing null", body: `{"operation":"add","a":10,"b":5} null`, status: http.StatusBadRequest},
		{name: "trailing garbage", body: `{"operation":"add","a":10,"b":5} invalid`, status: http.StatusBadRequest},
		{name: "missing operation", body: `{"a":10,"b":5}`, status: http.StatusBadRequest},
		{name: "empty operation", body: `{"operation":"","a":10,"b":5}`, status: http.StatusBadRequest},
		{name: "null operation", body: `{"operation":null,"a":10,"b":5}`, status: http.StatusBadRequest},
		{name: "non-string operation", body: `{"operation":1,"a":10,"b":5}`, status: http.StatusBadRequest},
		{name: "unsupported operation", body: `{"operation":"power","a":10,"b":5}`, status: http.StatusBadRequest},
		{name: "missing a", body: `{"operation":"add","b":5}`, status: http.StatusBadRequest},
		{name: "missing b", body: `{"operation":"add","a":10}`, status: http.StatusBadRequest},
		{name: "null a", body: `{"operation":"add","a":null,"b":5}`, status: http.StatusBadRequest},
		{name: "null b", body: `{"operation":"add","a":10,"b":null}`, status: http.StatusBadRequest},
		{name: "string operand", body: `{"operation":"add","a":"10","b":5}`, status: http.StatusBadRequest},
		{name: "boolean operand", body: `{"operation":"add","a":10,"b":true}`, status: http.StatusBadRequest},
		{name: "operand overflow", body: `{"operation":"add","a":1e400,"b":5}`, status: http.StatusBadRequest},
		{name: "division by zero", body: `{"operation":"divide","a":10,"b":0}`, status: http.StatusUnprocessableEntity},
		{name: "zero divided by zero", body: `{"operation":"divide","a":0,"b":0}`, status: http.StatusUnprocessableEntity},
		{name: "division by negative zero", body: `{"operation":"divide","a":10,"b":-0}`, status: http.StatusUnprocessableEntity},
		{name: "result overflow", body: `{"operation":"multiply","a":1e308,"b":10}`, status: http.StatusUnprocessableEntity},
	}

	for _, test := range tests {
		tester.Run(test.name, func(tester *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/api/calculate", strings.NewReader(test.body))
			request.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			handlers.Calculate(recorder, request)

			checkJSONResponse(tester, recorder, test.status)
		})
	}
}

func TestCalculateUnsupportedMethods(tester *testing.T) {
	for _, method := range []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions} {
		tester.Run(method, func(tester *testing.T) {
			request := httptest.NewRequest(method, "/api/calculate", strings.NewReader(`{"operation":"add","a":10,"b":5}`))
			request.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			handlers.Calculate(recorder, request)

			checkJSONResponse(tester, recorder, http.StatusMethodNotAllowed)
			if allow := recorder.Header().Get("Allow"); allow != http.MethodPost {
				tester.Errorf("Allow = %q; want %q", allow, http.MethodPost)
			}
		})
	}
}

func TestCalculateContentType(tester *testing.T) {
	tests := []struct {
		name        string
		contentType string
		status      int
	}{
		{name: "JSON", contentType: "application/json", status: http.StatusOK},
		{name: "JSON with charset", contentType: "application/json; charset=utf-8", status: http.StatusOK},
		{name: "missing", contentType: "", status: http.StatusUnsupportedMediaType},
		{name: "plain text", contentType: "text/plain", status: http.StatusUnsupportedMediaType},
		{name: "malformed", contentType: "application/json; charset", status: http.StatusUnsupportedMediaType},
	}

	for _, test := range tests {
		tester.Run(test.name, func(tester *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/api/calculate", strings.NewReader(`{"operation":"add","a":10,"b":5}`))
			if test.contentType != "" {
				request.Header.Set("Content-Type", test.contentType)
			}
			recorder := httptest.NewRecorder()

			handlers.Calculate(recorder, request)

			payload := checkJSONResponse(tester, recorder, test.status)
			if test.status == http.StatusOK && string(payload["result"]) != "15" {
				tester.Errorf("result = %s; want 15", payload["result"])
			}
		})
	}
}

func checkJSONResponse(tester *testing.T, recorder *httptest.ResponseRecorder, expectedStatus int) map[string]json.RawMessage {
	tester.Helper()
	if recorder.Code != expectedStatus {
		tester.Fatalf("status = %d; want %d; body = %s", recorder.Code, expectedStatus, recorder.Body.String())
	}
	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
		tester.Errorf("Content-Type = %q; want application/json", contentType)
	}
	var payload map[string]json.RawMessage
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		tester.Fatalf("decode response: %v", err)
	}
	if len(payload) != 1 {
		tester.Fatalf("response must contain exactly one field; got %s", recorder.Body.String())
	}
	if expectedStatus != http.StatusOK {
		var message string
		if err := json.Unmarshal(payload["error"], &message); err != nil {
			tester.Fatalf("decode error message: %v", err)
		}
		if strings.TrimSpace(message) == "" {
			tester.Error("error message must not be empty")
		}
	}
	return payload
}
