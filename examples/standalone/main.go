package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/felixgunawan/integrationtest/examples/standalone/config"
	"github.com/jmoiron/sqlx"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq" //using postgres
	"github.com/rs/cors"
)

var db *sqlx.DB

type Example struct {
	ID   string `json:"id",db:"id"`
	Name string `json:"name",db:"name"`
}

type ExampleInsert struct {
	Name string `json:"name"`
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", index).Methods("GET")
	r.HandleFunc("/add", insertExample).Methods("POST")
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})
	handler := c.Handler(r)
	fmt.Println("Server started at 55001")
	log.Fatal(http.ListenAndServe(":55001", handler))
}

type Message struct {
	Msg string `json:"message"`
}

func index(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Message{
		Msg: "hello",
	})
}

func insertExample(w http.ResponseWriter, r *http.Request) {
	var insertParam ExampleInsert
	err := json.NewDecoder(r.Body).Decode(&insertParam)
	if err == nil {
		json.NewEncoder(w).Encode(Message{
			Msg: err.Error(),
		})
		return
	}
	_, err = db.Exec(`INSERT INTO example (name) VALUES ($1)`, insertParam.Name)
	if err == nil {
		json.NewEncoder(w).Encode(Message{
			Msg: err.Error(),
		})
		return
	}
	json.NewEncoder(w).Encode(Message{
		Msg: "insert success",
	})
	return
}

func connectDb() (*sqlx.DB, error) {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.MainConfig.DB.Host, config.MainConfig.DB.Port, config.MainConfig.DB.User, config.MainConfig.DB.Pass, config.MainConfig.DB.Name)
	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		return db, err
	}
	err = db.Ping()
	return db, err
}
