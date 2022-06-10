package repo

import (
	"log"
	"os"
)

func GetData() []string {
	db, err := ConnToDb(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	var res []string
	sql, err := db.Query("SELECT fileName, fullFilePath, hashSum, algorithm FROM hashfiles;")
	var fil, check, path, algo string
	if err != nil {
		log.Fatalln(err)
	}
	for sql.Next() {
		err = sql.Scan(&fil, &path, &check, &algo)
		res = append(res, fil, path, check, algo)
	}
	return res
}

func PutTable() string {
	db, err := ConnToDb(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	_, err1 := db.Exec("INSERT INTO hashfiles (fileName, fullFilePath, hashSum, algorithm) VALUES ('5.txt','5/5.txt','4566','sha256');")
	if err1 != nil {
		log.Fatalf("%v", err)
	}
	res := "Succesful creation of table"
	return res
}
