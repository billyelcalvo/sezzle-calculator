package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"math"
	"mime"
	"net/http"

	"github.com/diegomora/sezzle-calculator/internal/calculator"
)

type calculateRequest struct {
	Operation string   `json:"operation"`
	OperandA  *float64 `json:"a"`
	OperandB  *float64 `json:"b"`
}

type calculateResponse struct {
	Result float64 `json:"result"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func Calculate(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.Header().Set("Allow", http.MethodPost)
		writeJSON(writer, http.StatusMethodNotAllowed, errorResponse{Error: "method must be POST"})
		return
	}

	mediaType, _, err := mime.ParseMediaType(request.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		writeJSON(writer, http.StatusUnsupportedMediaType, errorResponse{Error: "Content-Type must be application/json"})
		return
	}

	var input calculateRequest
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		writeJSON(writer, http.StatusBadRequest, errorResponse{Error: "body must be a valid JSON object with operation, a and b"})
		return
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		writeJSON(writer, http.StatusBadRequest, errorResponse{Error: "body must contain exactly one JSON object"})
		return
	}

	if input.Operation == "" {
		writeJSON(writer, http.StatusBadRequest, errorResponse{Error: "operation is required"})
		return
	}
	if input.OperandA == nil || input.OperandB == nil {
		writeJSON(writer, http.StatusBadRequest, errorResponse{Error: "a and b are required numbers"})
		return
	}

	var result float64
	switch input.Operation {
	case "add":
		result = calculator.Add(*input.OperandA, *input.OperandB)
	case "subtract":
		result = calculator.Subtract(*input.OperandA, *input.OperandB)
	case "multiply":
		result = calculator.Multiply(*input.OperandA, *input.OperandB)
	case "divide":
		result, err = calculator.Divide(*input.OperandA, *input.OperandB)
	default:
		writeJSON(writer, http.StatusBadRequest, errorResponse{Error: "operation must be add, subtract, multiply or divide"})
		return
	}

	if errors.Is(err, calculator.ErrDivisionByZero) {
		writeJSON(writer, http.StatusUnprocessableEntity, errorResponse{Error: "division by zero"})
		return
	}
	if err != nil {
		writeJSON(writer, http.StatusInternalServerError, errorResponse{Error: "calculation failed"})
		return
	}
	if math.IsNaN(result) || math.IsInf(result, 0) {
		writeJSON(writer, http.StatusUnprocessableEntity, errorResponse{Error: "result is outside the supported numeric range"})
		return
	}

	writeJSON(writer, http.StatusOK, calculateResponse{Result: result})
}

func writeJSON(writer http.ResponseWriter, status int, payload any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	if err := json.NewEncoder(writer).Encode(payload); err != nil {
		log.Printf("write JSON response: %v", err)
	}
}
