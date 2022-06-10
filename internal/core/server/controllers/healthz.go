package controllers

import (
	"fmt"
	"net/http"
	"time"
)

func Healthz(w http.ResponseWriter, r *http.Request) {
	started := time.Now()
	duration := time.Now().Sub(started)
	if duration.Seconds() > 10 {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("error: %v", duration.Seconds())))
	} else {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}
}
