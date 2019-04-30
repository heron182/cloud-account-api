package handlers

import "net/http"

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)

	w.Write([]byte(`{ "error": "Method not allowed"`))
}

func RouteNotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	w.Write([]byte(`{ "error": "Not found"`))
}
