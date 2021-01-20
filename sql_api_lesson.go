package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type UserCreateRequest struct {
	Name string `json:"name"`
}

type UserCreateReponse struct {
	Token string `json:"token"`
}

type UserGetReponse struct {
	Name string `json:"name"`
}

type UserUpdateRequest struct {
	Name string `json:"name"`
}

type GachaDrawRequest struct {
	Times int `json:"times"`
}

type Gacharesult struct {
	CharacterID string `json:"characterID"`
	Name        string `json:"name"`
}

type GachaDrawResponse struct {
	Results []Gacharesult `json:"results"`
}

type UserCharacter struct {
	UserCharacterID string `json:"userCharacterID"`
	CharacterID     string `json:"characterID"`
	Name            string `json:"name"`
}

type CharacterListResponse struct {
	Characters []UserCharacter `json:"characters"`
}

type ErrorMeessage struct {
	Error string `json:"errors"`
}

type Rate struct{
	CharacterID string
	Ratesum float32 
}

var db *sql.DB

//ランダムな文字列生成のための関数
func init() {
	rand.Seed(time.Now().UnixNano())
	var err error
	db, err = sql.Open("mysql", "root@/lesson1")
	log.Println("Connected to mysql.")
	//接続でエラーが発生した場合の処理
	if err != nil {
		log.Println("データベースに接続できませんでした")
	}

}

var rs1Letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = rs1Letters[rand.Intn(len(rs1Letters))]
	}
	return string(b)
}

func userGet(w http.ResponseWriter, r *http.Request) {

	xtoken := r.Header.Get("x-token")

	rows, err := db.Query(fmt.Sprintf("SELECT name from users WHERE token = '%s' ;", xtoken))
	if err != nil {
		log.Printf("xtoken:%sのユーザーを取得できませんでした", xtoken)
		var errormessage ErrorMeessage
		errormessage.Error = "xtoken:" + xtoken + "のユーザーを取得できませんでした"
		res, _ := json.Marshal(errormessage)
		w.Write(res)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	var person UserGetReponse
	for rows.Next() {
		err := rows.Scan(&person.Name)

		if err != nil {
			log.Printf("xtoken:%sのユーザーを読み取れませんでした", xtoken)
			var errormessage ErrorMeessage
			errormessage.Error = "xtoken:" + xtoken + "のユーザーを取得できませんでした"
			res, _ := json.Marshal(errormessage)
			w.Write(res)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}

	res, _ := json.Marshal(person)

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func userCreate(w http.ResponseWriter, r *http.Request) {

	//リクエストボディからデータを取得
	var req UserCreateRequest
	error := json.NewDecoder(r.Body).Decode(&req)
	if error != nil {
		fmt.Println(error)
		log.Printf("リクエストを受け取れませんでした")
		var errormessage ErrorMeessage
		errormessage.Error = "リクエストを受け取れませんでした"
		res, _ := json.Marshal(errormessage)
		w.Write(res)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	fmt.Printf("%sを登録します。", req.Name)

	var newid string
	var newtoken string
	newid = RandString(8)
	newtoken = RandString(8)

	fmt.Printf("newid:%s,newtoken:%s", newid, newtoken)

	ins, err := db.Prepare("INSERT INTO users(id,token,name) VALUES(?,?,?)")
	if err != nil {
		fmt.Println(error)
		log.Printf("%sをデータベースに追加できませんでした", req.Name)
		var errormessage ErrorMeessage
		errormessage.Error = req.Name + "をデータベースに追加できませんでした"
		res, _ := json.Marshal(errormessage)
		w.Write(res)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	ins.Exec(newid, newtoken, req.Name)

	var person UserCreateReponse
	person.Token = newtoken
	res, _ := json.Marshal(person)

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

}

func userUpdate(w http.ResponseWriter, r *http.Request) {

	xtoken := r.Header.Get("x-token")

	var req UserUpdateRequest
	error := json.NewDecoder(r.Body).Decode(&req)
	if error != nil {
		fmt.Println(error)
		log.Printf("リクエストをうけとれませんでした")
		var errormessage ErrorMeessage
		errormessage.Error = "リクエストをうけとれませんでした"
		res, _ := json.Marshal(errormessage)
		w.Write(res)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	ins, err := db.Prepare(fmt.Sprintf("UPDATE users SET name = '%s' WHERE token = '%s' ;", req.Name, xtoken))
	if err != nil {
		fmt.Println(err)
		log.Printf("更新に失敗しました")
		var errormessage ErrorMeessage
		errormessage.Error = "更新に失敗しました"
		res, _ := json.Marshal(errormessage)
		w.Write(res)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ins.Exec()

}

func gachaDraw(w http.ResponseWriter, r *http.Request) {
	xtoken := r.Header.Get("x-token")


	row, err := db.Query(fmt.Sprintf("SELECT id from users WHERE token = '%s' ;", xtoken))
	if err != nil {
		log.Printf("xtoken:%sのユーザーを取得できませんでした", xtoken)
		var errormessage ErrorMeessage
		errormessage.Error = "xtoken:" + xtoken + "のユーザーを取得できませんでした"
		res, _ := json.Marshal(errormessage)
		w.Write(res)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	var userid string
	for row.Next() {
		err := row.Scan(&userid)

		if err != nil {
			log.Printf("xtoken:%sのユーザーを読み取れませんでした", xtoken)
			var errormessage ErrorMeessage
			errormessage.Error = "xtoken:" + xtoken + "のユーザーを読み取れませんでした"
			res, _ := json.Marshal(errormessage)
			w.Write(res)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}

	var req GachaDrawRequest
	error := json.NewDecoder(r.Body).Decode(&req)
	if error != nil {
		fmt.Println(error)
		log.Printf("リクエストをうけとれませんでした")
		var errormessage ErrorMeessage
		errormessage.Error = "リクエストをうけとれませんでした"
		res, _ := json.Marshal(errormessage)
		w.Write(res)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	times := req.Times

	rows, err := db.Query("SELECT * FROM gachatable")
	if err != nil {
		fmt.Println(err)
		log.Printf("排出率テーブルから取得できませんでした")
		var errormessage ErrorMeessage
		errormessage.Error = "排出率テーブルから取得できませんでした"
		res, _ := json.Marshal(errormessage)
		w.Write(res)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var rates []Rate
	var ratesum float32
	for rows.Next() {
		var rate Rate
		var a float32
		error:=rows.Scan(&rate.CharacterID,&a)
		if error != nil {
			fmt.Println(error)
			log.Printf("排出率テーブルを読み取れませんでした")
			var errormessage ErrorMeessage
			errormessage.Error = "排出率テーブルを読みとれませんでした"
			res, _ := json.Marshal(errormessage)
			w.Write(res)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ratesum = ratesum+a
		rate.Ratesum=ratesum
		rates=append(rates,rate)

	}

	var gacharesponse GachaDrawResponse
	var gachatablerate float32
	var characterid string 
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < times; i++ {
		gachatablerate = rand.Float32()/ratesum
		for i := 0; i < len(rates); i++{
			if (gachatablerate < rates[i].Ratesum){
				characterid = rates[i].CharacterID
				break
			}
		}

		row2, err := db.Query(fmt.Sprintf("SELECT * FROM characters WHERE characterid = '%s' ; ", characterid))
		if err != nil {
			fmt.Println(err)
			log.Printf("キャラクターテーブルを取得できませんでした")
			var errormessage ErrorMeessage
			errormessage.Error = "キャラクターテーブルを取得できませんでした"
			res, _ := json.Marshal(errormessage)
			w.Write(res)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var result Gacharesult
		for row2.Next() {
			error := row2.Scan(&result.CharacterID, &result.Name)
			if error != nil {
				fmt.Println(error)
				log.Printf("キャラクターテーブルを読み取れませんでした")
				var errormessage ErrorMeessage
				errormessage.Error = "キャラクターテーブルを読みとれませんでした"
				res, _ := json.Marshal(errormessage)
				w.Write(res)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		var newusercharacterid string
		newusercharacterid = RandString(8)

		ins, err := db.Prepare(fmt.Sprintf("INSERT INTO usercharacter(usercharacterid,characterid,userid) VALUES(?,?,?)"))
		if err != nil {
			fmt.Println(err)
			log.Printf("ユーザーキャラクターテーブルを更新できませんでした")
			var errormessage ErrorMeessage
			errormessage.Error = "ユーザーキャラクターテーブルを更新できませんでした"
			res, _ := json.Marshal(errormessage)
			w.Write(res)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ins.Exec(newusercharacterid, result.CharacterID, userid)

		gacharesponse.Results = append(gacharesponse.Results, result)

	}

	res, _ := json.Marshal(gacharesponse)

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

}

func characterList(w http.ResponseWriter, r *http.Request) {
	xtoken := r.Header.Get("x-token")

	row, err := db.Query(fmt.Sprintf("SELECT id from users WHERE token = '%s' ;", xtoken))
	if err != nil {
		log.Printf("xtoken:%sのユーザーを取得できませんでした", xtoken)
		var errormessage ErrorMeessage
		errormessage.Error = "xtoken:" + xtoken + "のユーザーを取得できませんでした"
		res, _ := json.Marshal(errormessage)
		w.Write(res)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	var userid string
	for row.Next() {
		err := row.Scan(&userid)

		if err != nil {
			log.Printf("xtoken:%sのユーザーを読み取れませんでした", xtoken)
			var errormessage ErrorMeessage
			errormessage.Error = "xtoken:" + xtoken + "のユーザーを読み取れませんでした"
			res, _ := json.Marshal(errormessage)
			w.Write(res)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}

	var characterlistresponse CharacterListResponse

	row1, err := db.Query(fmt.Sprintf("SELECT usercharacter.usercharacterid , usercharacter.characterid,characters.name FROM usercharacter JOIN characters ON usercharacter.characterid = characters.characterid WHERE usercharacter.userid = '%s' ; ", userid))
	if err != nil {
		fmt.Println(err)
		log.Printf("ユーザーキャラクターテーブルを取得きませんでした")
		var errormessage ErrorMeessage
		errormessage.Error = "ユーザーキャラクターテーブルを取得きませんでした"
		res, _ := json.Marshal(errormessage)
		w.Write(res)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for row1.Next() {
		var usercharacter UserCharacter
		error2 := row1.Scan(&usercharacter.UserCharacterID, &usercharacter.CharacterID,&usercharacter.Name)
		if error2 != nil {
			fmt.Println(err)
			log.Printf("ユーザーキャラクターテーブルを読み込めませんでした")
			var errormessage ErrorMeessage
			errormessage.Error = "ユーザーキャラクターテーブルを取得きませんでした"
			res, _ := json.Marshal(errormessage)
			w.Write(res)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		characterlistresponse.Characters = append(characterlistresponse.Characters, usercharacter)

	}

	res, _ := json.Marshal(characterlistresponse)

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

}

func main() {

	http.HandleFunc("/user/get", userGet)
	http.HandleFunc("/user/create", userCreate)
	http.HandleFunc("/user/update", userUpdate)
	http.HandleFunc("/gacha/draw", gachaDraw)
	http.HandleFunc("/character/list", characterList)

	log.Println("Server running...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Listen and serve failed. %+v", err)
	}

}

//参考にしたサイト
//https://qiita.com/kkam0907/items/92d3d31c84c596eacaee
//https://qiita.com/rechtburg/items/b5eed25719582d7f490d
//https://qiita.com/ShinyaIshikawa/items/fede44cee7c71721247a
//https://ota42y.com/blog/2014/10/04/go-mysql/
