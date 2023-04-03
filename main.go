package main

import (
	"database/sql"
	"encoding/json"

	// "gc1/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/wijaysali/p2gc1/database"
	"github.com/wijaysali/p2gc1/models"
	// "github.com/your-username/rest-api/database"
	// "github.com/your-username/rest-api/models"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(db *sql.DB) {
	a.DB = db

	a.Router = mux.NewRouter()
	a.setRouters()
}

func (a *App) setRouters() {
	a.Router.HandleFunc("/users", a.getUsers).Methods("GET")
	a.Router.HandleFunc("/users/{id:[0-9]+}", a.getUserByID).Methods("GET")
	a.Router.HandleFunc("/users", a.createUser).Methods("POST")
	a.Router.HandleFunc("/users/{id:[0-9]+}", a.updateUser).Methods("PUT")
	a.Router.HandleFunc("/users/{id:[0-9]+}", a.deleteUser).Methods("DELETE")
}

func (a *App) getUsers(w http.ResponseWriter, r *http.Request) {
	repo := database.UserRepository{DB: a.DB}
	users, err := repo.GetAllUsers()
	if err != nil {
		log.Fatal(err)
	}

	responseJSON(w, users)
}

func (a *App) getUserByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatal(err)
	}

	repo := database.UserRepository{DB: a.DB}
	user, err := repo.GetUserByID(id)
	if err != nil {
		log.Fatal(err)
	}

	responseJSON(w, user)
}

func (a *App) createUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		responseError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	repo := database.UserRepository{DB: a.DB}
	id, err := repo.CreateUser(user)
	if err != nil {
		log.Fatal(err)
	}

	responseJSON(w, map[string]int{"id": id})
}

func (a *App) updateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatal(err)
	}

	var user models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		responseError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	user.ID = id

	repo := database.UserRepository{DB: a.DB}
	err = repo.UpdateUser(user)
	if err != nil {
		log.Fatal(err)
	}

	responseJSON(w, map[string]string{"message": "success"})
}

func (a *App) deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatal(err)
	}

	repo := database.UserRepository{DB: a.DB}
	err = repo.DeleteUser(id)
	if err != nil {
		log.Fatal(err)
	}

	responseJSON(w, map[string]string{"message": "success"})
}

func responseJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func responseError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Write([]byte(message))
}

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := App{}
	app.Initialize(db)
	log.Fatal(http.ListenAndServe(":8000", app.Router))
}
