package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dsn = "root:password@tcp(0.0.0.0:3306)/todolist?charset=utf8&parseTime=True&loc=Local"
var db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})

type TodoItemModel struct {
	Id          int `gorm:"primary_key"`
	Task        string
	LimitedTime bool
	Deadline    time.Time
	Finish      bool
	Overdue     bool
	Urgent      bool
}

func main() {
	db.Debug().AutoMigrate(&TodoItemModel{})

	log.Info("Starting To do list API Server")
	router := mux.NewRouter()
	router.HandleFunc("/healthz", Healthz).Methods("GET")
	router.HandleFunc("/create", createTask).Methods("POST")
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

func createTask(w http.ResponseWriter, r *http.Request) {
	var todoReq TodoItemModel
	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, &todoReq)

	log.WithFields(log.Fields{"task": todoReq.Task}).Info("Add new Task. Saving to database")
	deadline := todoReq.Deadline
	limitedTime := true
	if deadline.IsZero() {
		deadline = time.Now()
		limitedTime = false
	}
	todo := &TodoItemModel{Task: todoReq.Task, LimitedTime: limitedTime, Deadline: deadline, Finish: false, Overdue: false, Urgent: todoReq.Urgent}
	db.Create(&todo)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&todo)
}
