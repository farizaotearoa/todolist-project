package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
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
	router.HandleFunc("/todo/create", CreateTask).Methods("POST")
	router.HandleFunc("/todo/update/{id}", UpdateTask).Methods("PUT")
	router.HandleFunc("/todo/delete/{id}", DeleteTask).Methods("DELETE")
	http.ListenAndServe(":8080", router)

	handler := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "DELETE"},
	}).Handler(router)

	http.ListenAndServe(":8000", handler)
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

func CreateTask(w http.ResponseWriter, r *http.Request) {
	todoReq := &TodoItemModel{}
	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, &todoReq)

	deadline := todoReq.Deadline
	limitedTime := true
	if deadline.IsZero() {
		deadline = time.Now()
		limitedTime = false
	}
	todo := &TodoItemModel{Task: todoReq.Task, LimitedTime: limitedTime, Deadline: deadline, Finish: false, Overdue: false, Urgent: todoReq.Urgent}
	db.Create(&todo)

	log.WithFields(log.Fields{"task": todoReq.Task}).Info("Add new Task")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&todo)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	todo := &TodoItemModel{}
	result := db.First(&todo, id)
	if result.Error != nil {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"error": "Record Not Found"}`)
	} else {
		finish, _ := strconv.ParseBool(r.FormValue("finish"))
		log.WithFields(log.Fields{"Id": id, "Task": todo.Task, "Finish": finish}).Info("Updating Task")
		todo := &TodoItemModel{}
		db.First(&todo, id)
		todo.Finish = finish
		db.Save(&todo)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&todo)
	}
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	todo := &TodoItemModel{}
	result := db.First(&todo, id)
	if result.Error != nil {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"error": "Record Not Found"}`)
	} else {
		log.WithFields(log.Fields{"Id": id, "Task": todo.Task}).Info("Deleting Task")
		todo := &TodoItemModel{}
		db.First(&todo, id)
		db.Delete(&todo)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{}`)
	}
}
