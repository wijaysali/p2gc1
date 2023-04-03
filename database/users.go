package database

import (
	"database/sql"
	"log"
	"models/Downloads/golang/pkg/mod/github.com/wijaysali/p2gc1@v0.0.0-20230402080946-46c235460cff/models"
	// "models/Downloads/golang/pkg/mod/github.com/wijaysali/p2gc1@v0.0.0-20230402080946-46c235460cff/models"
	// // "github.com/wijaysali/p2gc1/models"
)

func GetUserByID(id int64) (*models.User, error) {
	db, err := InitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	user := &models.User{}
	row := db.QueryRow("SELECT id, first_name, last_name FROM users WHERE id=?", id)
	err = row.Scan(&user.ID, &user.FirstName, &user.LastName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			log.Println(err)
			return nil, err
		}
	}
	return user, nil
}

func GetAllUsers() ([]*models.User, error) {
	db, err := InitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, first_name, last_name FROM users")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	users := make([]*models.User, 0)
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	return users, nil
}

func CreateUser(user *models.User) (int64, error) {
	db, err := InitDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	result, err := db.Exec("INSERT INTO users (first_name, last_name) VALUES (?, ?)", user.FirstName, user.LastName)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return id, nil
}

func UpdateUser(user *models.User) error {
	db, err := InitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE users SET first_name=?, last_name=? WHERE id=?", user.FirstName, user.LastName, user.ID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func DeleteUser(id int64) error {
	db, err := InitDB()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
