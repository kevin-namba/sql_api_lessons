package main

import (
  "database/sql"
    "encoding/json"
    "fmt"
    "log"
  //  "math/rand"
    "net/http"
    _ "github.com/go-sql-driver/mysql"
)

type ResponseData struct {
    ID int    `json:"ID"`
    Name   string `json:"name"`
}



func userHandler(w http.ResponseWriter, r *http.Request) {

  db, err := sql.Open("mysql", "root@/lesson1")
  log.Println("Connected to mysql.")
  //接続でエラーが発生した場合の処理
  if err != nil {
      log.Fatal(err)
  }
  defer db.Close()

    //データベースへクエリを送信。引っ張ってきたデータがrowsに入る。
    rows, err := db.Query("SELECT * FROM users")
    defer rows.Close()
    if err != nil {
        panic(err.Error())
    }


    type ResponseDatas []ResponseData

    var responsedatas ResponseDatas

    //レコード一件一件をあらかじめ用意しておいた構造体に当てはめていく。
    for rows.Next() {
        var person ResponseData //構造体Responsedata型の変数personを定義
        err := rows.Scan(&person.ID, &person.Name)

        if err != nil {
            panic(err.Error())
        }

        responsedatas = append(responsedatas,person)

    }

    res, err := json.Marshal(responsedatas)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    fmt.Fprint(w, string(res))
}



func main() {


    http.HandleFunc("/user", userHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
//参考にしたサイト
//https://qiita.com/kkam0907/items/92d3d31c84c596eacaee
//https://qiita.com/rechtburg/items/b5eed25719582d7f490d
