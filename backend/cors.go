package main

import "net/http"

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Vary", "Origin")
		if request.Header.Get("Origin") == "http://localhost:5173" {
			writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
			writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}

		if request.Method == http.MethodOptions {
			writer.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(writer, request)
	})
}
