package uberproxy

import (
	"encoding/json"
	"log"
	"net/http"
)

func readJsonBody[T any](w http.ResponseWriter, r *http.Request) (v T, err error) {
	err = json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}
	return
}
func writeJsonBody(w http.ResponseWriter, v interface{}) (err error) {
	w.Header().Add("content-type", "application/json")
	err = json.NewEncoder(w).Encode(v)
	if err != nil {
		return
	}
	return
}
func wrapHandler(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			log.Println(err)
			w.WriteHeader(500)
		}
	}
}
