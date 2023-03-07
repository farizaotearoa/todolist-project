package main

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dsn = "root:password@tcp(0.0.0.0:3306)/todolist?charset=utf8&parseTime=True&loc=Local"
var db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})

type TodoItemModel struct {
	Id     int `gorm:"primary_key"`
	Task   string
	Finish bool
}

func main() {
	db.Debug().AutoMigrate(&TodoItemModel{})

	log.Info("Starting To do list API Server")
	router := mux.NewRouter()
	router.HandleFunc("/healthz", Healthz).Methods("GET")
	http.ListenAndServe(":8080", router)
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	log.Info("API Health is OK")
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}
