package db

import (
	"fmt"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	TimeFormat = "2006/01/02 15:04:05"
)

// passed
func TestGetTalks(t *testing.T) {
	// テスト用のレコードを追加
	posts := []struct {
		talk   string
		create time.Time
		update time.Time
		delete interface{}
	}{
		{"こんにちは", time.Now(), time.Now(), nil},
		{"消すけど消えちゃダメなもの", time.Now(), time.Now(), time.Now()},
		{"消さないもの", time.Now(), time.Now(), nil},
		{"消さないもの", time.Now(), time.Now(), time.Now()},
	}
	for _, post := range posts {
		err := PostTalkSelectTime(post.talk, post.create, post.update, post.delete)
		if err != nil {
			fmt.Printf("error: %v", err)
			return
		}
	}
	// テスト
	talks, err := GetTalks()
	if err != nil {
		t.Errorf("error: %v", err)
	}
	count := 0
	for _, talk := range talks {
		fmt.Println(talk.ID, talk.Talk, talk.CreateAt.Format(TimeFormat), talk.UpdateAt.Format(TimeFormat), talk.DeleteAt.Format(TimeFormat))
		count++
	}
	if count != 2 {
		t.Error("削除されたレコードも取得している可能性があります。")
	}
	// 追加したレコードの削除
	DeleteAllRecord()
}

// passed
func TestPostTalk(t *testing.T) {
	// テスト
	type args struct {
		t string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "書き込み1", args: args{t: "こんにちは"}, wantErr: false},
		{name: "書き込み1", args: args{t: "こんばんは"}, wantErr: false},
		{name: "書き込み2", args: args{t: "ああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああああ"}, wantErr: true},
		{name: "書き込み3", args: args{t: ""}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PostTalk(tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("PostTalk() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	// 追加したレコードの削除
	DeleteAllRecord()
}

func TestUpdateTalk(t *testing.T) {
	// テスト用のレコード追加
	posts := []struct {
		talk   string
		create time.Time
		update time.Time
		delete interface{}
	}{
		{"変更前１", time.Now(), time.Now(), nil},
		{"変更前２", time.Now(), time.Now(), nil},
		{"変更前３", time.Now(), time.Now(), nil},
	}
	for _, post := range posts {
		err := PostTalkSelectTime(post.talk, post.create, post.update, post.delete)
		if err != nil {
			fmt.Printf("PostTalkSelectTime() error: %v", err)
			return
		}
	}
	// テスト
	talk, err := GetTalk(3)
	if err != nil {
		t.Errorf("GetTalk() error: %v", err)
	}
	talk.Talk = "変更後"
	time.Sleep(time.Second * 2)
	if err := UpdateTalk(talk); err != nil {
		t.Errorf("UpdateTark() error: %v", err)
	}
	// テスト結果の確認
	afterTalk, err := GetTalk(3)
	if err != nil {
		t.Errorf("GetTalk() error: %v", err)
	}
	if afterTalk.Talk != "変更後" {
		t.Errorf("want: 変更後 result: %v", afterTalk.Talk)
	}
	afters, err := GetTalks()
	for _, after := range afters {
		fmt.Println(after.ID, after.Talk, after.CreateAt.Format(TimeFormat), after.UpdateAt.Format(TimeFormat))
	}
	DeleteAllRecord()
}

type DeleteTest struct {
	id     int
	result bool // エラーあり：true, エラーなし：false
}

// passed
func TestDeleteTalk(t *testing.T) {
	// テスト用のレコード追加
	posts := []struct {
		talk   string
		create time.Time
		update time.Time
		delete interface{}
	}{
		{"消したいもの", time.Now(), time.Now(), nil},
		{"消すけど消えちゃダメなもの", time.Now(), time.Now(), time.Now()},
		{"消さないもの", time.Now(), time.Now(), nil},
	}
	for _, post := range posts {
		err := PostTalkSelectTime(post.talk, post.create, post.update, post.delete)
		if err != nil {
			fmt.Printf("PostTalkSelectTime() error: %v", err)
			return
		}
	}
	// テスト
	tests := []DeleteTest{
		{1, false}, {2, true},
	}
	for _, test := range tests {
		err := DeleteTalk(test.id)
		if (err != nil) != test.result {
			t.Errorf("DeleteTalk() err: %v, want: %v", err, test.result)
		}
	}
	// テストの結果を表示 id=３だけが表示されるはず。
	ts, err := GetTalks()
	if err != nil {
		fmt.Printf("GetTalks() error: %v", err)
		return
	}
	for _, t := range ts {
		fmt.Println(t)
	}
	// 追加したレコードの削除
	DeleteAllRecord()
}
