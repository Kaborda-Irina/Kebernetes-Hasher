package controllers

import (
	"fmt"
	"github.com/Kaborda-Irina/Kubernetes-Hasher/internal/core/server/repo"
	"net/http"
)

func GetData(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, repo.GetData())
}
