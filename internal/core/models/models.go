package models

import "k8s.io/client-go/kubernetes"

type HashDataFromDB struct {
	ID           int
	Hash         string
	FileName     string
	FullFilePath string
	Algorithm    string
}

type DeletedHashes struct {
	FileName    string
	OldChecksum string
	FilePath    string
	Algorithm   string
}

type ConnectionDB struct {
	Dbdriver   string
	DbUser     string
	DbPassword string
	DbPort     string
	DbHost     string
	DbName     string
}

type KuberData struct {
	Clientset  *kubernetes.Clientset
	Namespace  string
	TargetName string
	TargetType string
}
