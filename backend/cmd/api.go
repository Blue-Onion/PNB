package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Router() *chi.Mux {
	router := chi.NewMux()
	router.Get("/check/{hexcode}", CheckCode())
	return router
}
func CheckCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hexcode := chi.URLParam(r, "hexcode")
		res, err := checkCode(hexcode)
		if err != nil {
			RespondWithJson(w, 400, err.Error())
			return
		}
		RespondWithJson(w, 200, json.RawMessage(res))
	}
}

type res struct {
	Success bool
	Data    any
}

func RespondWithJson(w http.ResponseWriter, code int, payload any) {
	res := res{}
	if code < 299 {
		res.Success = true
	} else {
		res.Success = false
	}
	res.Data = payload
	data, err := json.Marshal(res)
	if err != nil {
		log.Println("Failed to parse Json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-type", "Application/Json")
	w.WriteHeader(code)
	w.Write(data)

}
