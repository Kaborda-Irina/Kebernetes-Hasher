package controllers

import (
	"fmt"
	"github.com/Kaborda-Irina/Kubernetes-Hasher/internal/core/server/repo"
	"net/http"
)

func PutData(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, repo.PutTable())
}
