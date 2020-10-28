package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// Talk golang_dbのtalksテーブルの構造
type Talk struct {
	ID       int
	Talk     string
	CreateAt time.Time
	UpdateAt time.Time
	DeleteAt time.Time
}

const (
	// DriverName ドライバ名(mysql固定)
	dName = "mysql"
	// DataSourceName user:password@tcp(container-name:port)/dbname
	dsName = "root:golang@tcp(mysql-container:3306)/golang_db?parseTime=true&loc=Asia%2FTokyo"
)

/*
todo: CRUD操作ごとにDBに接続するようにする。
*/

// GetTalks データベースから投稿されたトークを取得する。
// todo: deleteAtがあるかどうかで取得するかどうかを判断。
func GetTalks() (talks []Talk, err error) {
	db, err := sql.Open(dName, dsName)
	if err != nil {
		return
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM talks WHERE delete_at IS NULL ORDER BY create_at DESC")
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var talk Talk
		if err := rows.Scan(
			&talk.ID,
			&talk.Talk,
			&talk.CreateAt,
			&talk.UpdateAt,
			&talk.DeleteAt,
		); err != nil {
			// log.Print(err)
		}
		talks = append(talks, talk)
	}
	return
}

// GetTalk 指定されたIDのトークをデータベースから取得する。
func GetTalk(id int) (talk Talk, err error) {

	db, err := sql.Open(dName, dsName)
	if err != nil {
		return
	}
	defer db.Close()
	strID := fmt.Sprint(id)
	query := "SELECT * FROM talks WHERE id=" + strID
	rows, err := db.Query(query)
	if err != nil {
		return
	}
	defer rows.Close()
	rows.Next()
	if err := rows.Scan(
		&talk.ID,
		&talk.Talk,
		&talk.CreateAt,
		&talk.UpdateAt,
		&talk.DeleteAt,
	); err != nil {
		log.Print(err)
	}
	return
}

// StringLengthError 文字列の長さによるエラー
type StringLengthError struct {
	text string
}

func (e *StringLengthError) Error() string {
	return fmt.Sprintf("text %v length over 140 or 0.", e.text)
}

// PostTalk データベースにトークを投稿する。1~140文字の間
func PostTalk(t string) error {
	if len(t) <= 0 || len(t) > 140 {
		return &StringLengthError{t}
	}
	db, err := sql.Open(dName, dsName)
	if err != nil {
		return err
	}
	defer db.Close()
	now := time.Now()
	ins, err := db.Prepare("INSERT INTO talks(talk,create_at,update_at) VALUES(?,?,?)")
	if err != nil {
		return err
	}
	ins.Exec(t, now, now)
	return nil
}

// PostTalkSelectTime データベースにトークを時間指定で投稿する。1~140文字の間
func PostTalkSelectTime(t string, create time.Time, update time.Time, delete interface{}) error {
	if len(t) <= 0 || len(t) > 140 {
		return &StringLengthError{t}
	}
	db, err := sql.Open(dName, dsName)
	if err != nil {
		return err
	}
	defer db.Close()
	ins, err := db.Prepare("INSERT INTO talks(talk,create_at,update_at, delete_at) VALUES(?,?,?,?)")
	if err != nil {
		return err
	}
	ins.Exec(t, create, update, delete)
	return nil
}

// UpdateTalk IDで指定したトークを書き換える。
func UpdateTalk(talk Talk) (err error) {
	db, err := sql.Open(dName, dsName)
	if err != nil {
		return
	}
	defer db.Close()
	upd, err := db.Prepare("UPDATE talks SET talk=?, update_at=? WHERE id=?")
	if err != nil {
		return
	}
	result, err := upd.Exec(talk.Talk, time.Now(), talk.ID)
	if err != nil {
		return
	}
	log.Print(result)
	return
}

// ExistError 指定したidのレコードが存在しないことによるエラー
type ExistError struct {
	id string
}

func (e *ExistError) Error() string {
	return fmt.Sprintf("id: %v is already deleted.", e.id)
}

// DeleteTalk IDで指定したトークを削除する。deleteAtを追記。
func DeleteTalk(id int) (err error) {
	db, err := sql.Open(dName, dsName)
	if err != nil {
		return
	}
	defer db.Close()
	// すでに削除されていないかの確認
	strID := fmt.Sprint(id)
	exist, err := db.Query("SELECT delete_at FROM talks WHERE id=" + strID + " AND delete_at IS NOT NULL")
	if err != nil {
		return
	}

	count := 0
	for exist.Next() {
		if err := exist.Scan(
			&time.Time{},
		); err != nil {
			log.Print(err)
		}
		count++
		// log.Printf("get id=%v deleteAt:%v", strId, deleteAt.Format("2006/01/02 15:04:05"))
	}
	if count > 0 {
		err = &ExistError{strID}
		return
	}
	// 削除の実行
	del, err := db.Prepare("UPDATE talks SET delete_at=NOW() WHERE id=?")
	if err != nil {
		return
	}
	result, err := del.Exec(id)
	if err != nil {
		return
	}
	log.Print(result)
	return
}

// DeleteAllRecord talksテーブルのレコードを全て削除する。delete_atの書き込みではない。
func DeleteAllRecord() {
	db, err := sql.Open(dName, dsName)
	if err != nil {
		return
	}
	defer db.Close()

	dels, err := db.Prepare("TRUNCATE TABLE talks")
	if err != nil {
		log.Printf("error: %v", err)
	}
	dels.Exec()
}
