package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wijaysali/p2gc1/models"
)

type UserHandler struct {
	DB *sql.DB
}

func (u *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Query the database for all users
	rows, err := u.DB.Query("SELECT * FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Iterate over the rows and append to a slice of users
	users := []models.User{}
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshal the slice of users to JSON and write to the response writer
	jsonBytes, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (u *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the request URL parameters
	vars := mux.Vars(r)
	id := vars["id"]

	// Query the database for the user with the specified ID
	row := u.DB.QueryRow("SELECT * FROM users WHERE id = ?", id)

	// Scan the row into a user struct
	var user models.User
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshal the user to JSON and write to the response writer
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Decode the request body into a user struct
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert the user into the database
	result, err := u.DB.Exec("INSERT INTO users (first_name, last_name, email) VALUES (?, ?, ?)", user.FirstName, user.LastName, user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the ID of the newly inserted user
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the ID of the user struct and marshal to JSON
	user.ID = int(id)
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshall the user to JSON and write to the response writer
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)
}

func (u *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the request URL parameters
	vars := mux.Vars(r)
	id := vars["id"]

	// Decode the request body into a user struct
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update the user in the database
	_, err = u.DB.Exec("UPDATE users SET first_name = ?, last_name = ?, email = ? WHERE id = ?", user.FirstName, user.LastName, user.Email, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the ID of the user struct and marshal to JSON
	user.ID = int(id)
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)

}

func (u *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the request URL parameters
	vars := mux.Vars(r)
	id := vars["id"]

	// Delete the user from the database
	_, err := u.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
