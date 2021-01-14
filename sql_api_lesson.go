package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "time"
    "net/http"
    _ "github.com/go-sql-driver/mysql"
)

type ResponseData struct {
    ID string    `json:"ID"`
    Token string  `json:"token"`
    Name   string `json:"name"`
}

type RequestData struct{
    Name   string `json:"name"`
}

type UserGetReponse struct{
    Name   string `json:"name"`
}

type UserCreateReponse struct{
    Token string  `json:"token"`
}

type UserUpdateReponse struct{
    Name   string `json:"name"`
}




//ランダムな文字列生成のための関数
func init() {
    rand.Seed(time.Now().UnixNano())
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
  //r.Header.Set("Content-Type","application/json")
  xtoken := r.Header.Get("x-token")
  //fmt.Printf("xtoken:%s", xtoken)
  

  db, err := sql.Open("mysql", "root@/lesson1")
  log.Println("Connected to mysql.")
  //接続でエラーが発生した場合の処理
  if err != nil {
      log.Fatal(err)
  }
  defer db.Close()

    //データベースへクエリを送信。引っ張ってきたデータがrowsに入る。
    rows, err := db.Query(fmt.Sprintf("SELECT name from users WHERE token = '%s' ;",xtoken))
    if err != nil {
       log.Fatal(err)
       
    }

    
    defer db.Close()
    var person UserGetReponse //構造体Responsedata型の変数personを定義
    for rows.Next() {
        
    err :=rows.Scan(&person.Name)
        
    if err != nil {
            panic(err.Error())
                }
                 
    }
    
    res, err := json.Marshal(person)
    
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprint(w, string(res))
}

func userCreate(w http.ResponseWriter, r *http.Request){



//リクエストボディからデータを取得
  var req RequestData
  error := json.NewDecoder(r.Body).Decode(&req)
	if error != nil {
	fmt.Println(error)
			return
	}


  fmt.Printf("%sを登録します。",req.Name)

//w.Header().Set("Content-Type", "application/json")

  // db接続
  db, err := sql.Open("mysql", "root@/lesson1")
  log.Println("Connected to mysql.")
  //接続でエラーが発生した場合の処理
  if err != nil {
      log.Fatal(err)
  }

  var  newid string
  var newtoken string
  newid =RandString(8)
  newtoken =RandString(8)

  fmt.Printf("newid:%s,newtoken:%s", newid,newtoken)
  //データを挿入する

ins, err := db.Prepare("INSERT INTO users(id,token,name) VALUES(?,?,?)")
if err != nil {
    log.Fatal(err)
   }

ins.Exec(newid, newtoken,req.Name)
defer db.Close()

var person UserCreateReponse 
person.Token=newtoken
res, err := json.Marshal(person)
    
w.Header().Set("Content-Type", "application/json")
fmt.Fprint(w, string(res))
     
}

func userUpdate(w http.ResponseWriter, r *http.Request){

  xtoken := r.Header.Get("x-token")

  var req RequestData
  error := json.NewDecoder(r.Body).Decode(&req)
  if error != nil {
  fmt.Println(error)
      return
  }


  db, err := sql.Open("mysql", "root@/lesson1")
  log.Println("Connected to mysql.")
  //接続でエラーが発生した場合の処理
  if err != nil {
      log.Fatal(err)
  }
  ins, err := db.Prepare(fmt.Sprintf("UPDATE users SET name = '%s' WHERE token = '%s' ;",req.Name,xtoken))
  if err != nil {
    log.Fatal(err)
   }

ins.Exec()
defer db.Close()
}

func main() {    
    //req, _ := http.NewRequest("GET","localhost:8080" , nil)
    //req.Header.Set("Content-Type", "application/json")
    http.HandleFunc("/user/get", userGet)
    http.HandleFunc("/user/create", userCreate)
    http.HandleFunc("/user/update", userUpdate)
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
