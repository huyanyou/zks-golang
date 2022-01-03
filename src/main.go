package main

import (
	l "HeDa/src/service/login"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	r := httprouter.New()
	r.POST("/", l.Login)
	h := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	h.ListenAndServe()
}
