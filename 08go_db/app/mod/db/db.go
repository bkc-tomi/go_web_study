package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Talks struct {
	ID       int
	talk     string
	createAt time.Time
	updateAt time.Time
	deleteAt time.Time
}

const (
	// DriverName ドライバ名(mysql固定)
	dName = "mysql"
	// DataSourceName user:password@tcp(container-name:port)/dbname
	dsName = "root:golang@tcp(mysql-container:3306)/golang_db"
)

/*
todo: CRUD操作ごとにDBに接続するようにする。
*/

// GetTalks get all talks from db
func GetTalks() (talks []Talks, err error) {
	db, err := sql.Open(dName, dsName)
	if err != nil {
		return
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM talks")
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var talk Talks
		if err := rows.Scan(
			&talk.ID,
			&talk.talk,
			&talk.createAt,
			&talk.updateAt,
			&talk.deleteAt,
		); err != nil {
			log.Print(err)
		}
		talks = append(talks, talk)
	}
	log.Print("get from db:", talks)
	return
}

func PostTalk() {

}

func UpdateTalk() {

}

func DeleteTalk() {

}
