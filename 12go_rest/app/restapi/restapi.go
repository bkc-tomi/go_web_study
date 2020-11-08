package restapi

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	// DriverName ドライバ名(mysql固定)
	dName = "mysql"
	// DataSourceName user:password@tcp(container-name:port)/dbname
	dsName = "root:golang@tcp(mysql-user-container:3306)/golang_db?parseTime=true&loc=Asia%2FTokyo"
)

// User structure of user data
type User struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Profile     string    `json:"profile"`
	DateOfBirth string    `json:"date_of_birth"`
	CreateAt    time.Time `json:"create_at"`
	UpdateAt    time.Time `json:"update_at"`
}

// GetUsers get all user data from database
func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var users []User
	db, err := sql.Open(dName, dsName)
	if err != nil {
		fmt.Println("db open error", err)
		return
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		fmt.Println("db query error:", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		rows.Scan(&user.ID, &user.Name, &user.Profile, &user.DateOfBirth, &user.CreateAt, &user.UpdateAt)
		users = append(users, user)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
