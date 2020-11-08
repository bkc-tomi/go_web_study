package restapi

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	// DriverName ドライバ名(mysql固定)
	dName = "mysql"
	// DataSourceName user:password@tcp(container-name:port)/dbname
	dsName = "root:golang@tcp(mysql-user-container:3306)/golang_db?parseTime=true&loc=Asia%2FTokyo"
)

// User structure of user data
type User struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Profile     string `json:"profile"`
	DateOfBirth string `json:"date_of_birth"`
	CreateAt    string `json:"create_at"`
	UpdateAt    string `json:"update_at"`
}

func getUsersFromDB(query string) (users []User, err error) {
	db, err := sql.Open(dName, dsName)
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		return
	}
	defer rows.Close()

	values := make([]sql.RawBytes, 6)
	scanArgs := make([]interface{}, 6)
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		var user User
		err = rows.Scan(scanArgs...)
		if err != nil {
			return
		}
		user.ID, _ = strconv.Atoi(string(values[0]))
		user.Name = string(values[1])
		user.Profile = string(values[2])
		user.DateOfBirth = string(values[3])
		user.CreateAt = string(values[4])
		user.UpdateAt = string(values[5])
		users = append(users, user)
	}
	return
}

func setMsgToResponse(w http.ResponseWriter, status int, msgType string, msg string) {
	w.WriteHeader(status)
	errMsg := `{"` + msgType + `": "` + msg + `"}`
	w.Write([]byte(errMsg))
}

func parseRequestBody(r *http.Request) (jsonBody map[string]interface{}, err error) {
	//To allocate slice for request body
	length, err := strconv.Atoi(r.Header.Get("Content-Length"))
	if err != nil {
		return
	}

	//Read body data to parse json
	body := make([]byte, length)
	length, err = r.Body.Read(body)
	if err != nil && err != io.EOF {
		return
	}
	//parse json
	err = json.Unmarshal(body[:length], &jsonBody)
	if err != nil {
		return
	}
	return
}

// GetUsers get all user data from database
func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users, err := getUsersFromDB("SELECT * FROM users")
	if err != nil {
		setMsgToResponse(w, http.StatusInternalServerError, "error", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// GetUser get user data from db selected by id
func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	users, err := getUsersFromDB("SELECT * FROM users WHERE id=" + vars["id"])
	if err != nil {
		setMsgToResponse(w, http.StatusInternalServerError, "error", err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users[0])
}

// SearchUser search user data on keyword selected by query
func SearchUser(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	users, err := getUsersFromDB(`SELECT * FROM users WHERE name LIKE "%` + keyword + `%"`)
	if err != nil {
		setMsgToResponse(w, http.StatusInternalServerError, "error", err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// PostUser post user to database
func PostUser(w http.ResponseWriter, r *http.Request) {
	// parse request to json
	jsonBody, err := parseRequestBody(r)
	if err != nil {
		setMsgToResponse(w, http.StatusInternalServerError, "error", err.Error())
		return
	}
	if jsonBody["name"] == nil || jsonBody["name"] == "" {
		setMsgToResponse(w, http.StatusBadRequest, "error", "名前が入力されていません。")
		return
	}
	// write db
	db, err := sql.Open(dName, dsName)
	if err != nil {
		setMsgToResponse(w, http.StatusInternalServerError, "error", err.Error())
		return
	}
	defer db.Close()
	ins, err := db.Prepare("INSERT INTO users(name, profile, date_of_birth, create_at, update_at) VALUES(?,?,?, NOW(), NOW())")
	if err != nil {
		setMsgToResponse(w, http.StatusInternalServerError, "error", err.Error())
		return
	}
	ins.Exec(jsonBody["name"], jsonBody["profile"], jsonBody["date_of_birth"])
	// create response
	w.Header().Set("Content-Type", "application/json")
	setMsgToResponse(w, http.StatusCreated, "message", "登録しました。")
	return
}

// PutUser rewrite user data selected by id
func PutUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	users, err := getUsersFromDB(`SELECT * FROM users WHERE id=` + vars["id"])
	if err != nil {
		setMsgToResponse(w, http.StatusInternalServerError, "error", err.Error())
		return
	}
	if len(users) == 0 {
		setMsgToResponse(w, http.StatusNotFound, "error", "ユーザが見つかりませんでした。")
		return
	}
	// parse request to json
	jsonBody, err := parseRequestBody(r)
	if err != nil {
		setMsgToResponse(w, http.StatusInternalServerError, "error", err.Error())
		return
	}
	if jsonBody["name"] == nil || jsonBody["name"] == "" {
		setMsgToResponse(w, http.StatusBadRequest, "error", "名前が入力されていません。")
		return
	}
	// update db data
	db, err := sql.Open(dName, dsName)
	if err != nil {
		setMsgToResponse(w, http.StatusInternalServerError, "error", err.Error())
		return
	}
	defer db.Close()
	ins, err := db.Prepare("UPDATE users SET name=?, profile=?, date_of_birth=?, update_at=NOW() WHERE id = ?")
	if err != nil {
		setMsgToResponse(w, http.StatusInternalServerError, "error", err.Error())
		return
	}
	ins.Exec(jsonBody["name"], jsonBody["profile"], jsonBody["date_of_birth"], vars["id"])
	// response message
	w.Header().Set("Content-Type", "application/json")
	setMsgToResponse(w, http.StatusCreated, "message", "ユーザ情報を更新しました。")
	return
}

// DeleteUser delete user data
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	users, err := getUsersFromDB(`SELECT * FROM users WHERE id=` + vars["id"])
	if err != nil {
		setMsgToResponse(w, http.StatusInternalServerError, "error", err.Error())
		return
	}
	if len(users) == 0 {
		setMsgToResponse(w, http.StatusNotFound, "error", "ユーザが見つかりませんでした。")
		return
	}
	// delete db data
	db, err := sql.Open(dName, dsName)
	if err != nil {
		setMsgToResponse(w, http.StatusInternalServerError, "error", err.Error())
		return
	}
	defer db.Close()
	ins, err := db.Prepare("DELETE FROM users WHERE id=?")
	if err != nil {
		setMsgToResponse(w, http.StatusInternalServerError, "error", err.Error())
		return
	}
	ins.Exec(vars["id"])
	// response message
	w.Header().Set("Content-Type", "application/json")
	setMsgToResponse(w, http.StatusCreated, "message", "ユーザを削除しました。")
	return
}
