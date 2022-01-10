package main

import (
	l "HeDa/src/service/login"
	"HeDa/src/service/skeleton"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	r := httprouter.New()
	r.AddBerforeHandle(skeleton.GlobalMiddle)
	r.GET("/login", l.Login)
	r.POST("/logon", l.Logon)
	h := &http.Server{
		Addr:           ":9090",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	h.ListenAndServe()
}
