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



type UserCreateRequest struct{
    Name string `json:"name"`
}

type UserCreateReponse struct{
    Token string  `json:"token"`
}

type UserGetReponse struct{
    Name string `json:"name"`
}

type UserUpdateRequest struct{
    Name string `json:"name"`
}

type GachaDrawRequest struct{
    Times int `json:"times"`
}

type Gacharesult struct{
    CharacterID string `json:"characterID"`
    Name string `json:"name"`
}


type GachaDrawResponse struct{
    Results []Gacharesult `json:"results"`
}

type UserCharacter struct{
    UserCharacterID string `json:"userCharacterID"`
    CharacterID string `json:"characterID"`
    Name string `json:"name"`

}

type CharacterListResponse struct{
    Characters []UserCharacter `json:"characters"`
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

    var person UserGetReponse 
    for rows.Next() {
        err :=rows.Scan(&person.Name)
        
    if err != nil {
            panic(err.Error())
                }
                 
    }
    
    res, err := json.Marshal(person)
    
    w.Header().Set("Content-Type", "application/json")
    w.Write(res)
}

func userCreate(w http.ResponseWriter, r *http.Request){



    //リクエストボディからデータを取得
      var req UserCreateRequest
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
    w.Write(res)
         
}

func userUpdate(w http.ResponseWriter, r *http.Request){

  xtoken := r.Header.Get("x-token")

  var req UserUpdateRequest
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

func gachaDraw(w http.ResponseWriter, r *http.Request){
    xtoken := r.Header.Get("x-token")

    var req GachaDrawRequest
    error := json.NewDecoder(r.Body).Decode(&req)
    if error != nil {
    fmt.Println(error)
      return
    }
    times:=req.Times

    db, err := sql.Open("mysql", "root@/lesson1")
    log.Println("Connected to mysql.")
    //接続でエラーが発生した場合の処理
    if err != nil {
        log.Fatal(err)
    }

    rows, err := db.Query("SELECT * FROM gachatable")
    defer rows.Close()
    if err != nil {
        panic(err.Error())
    }

    var co int
    for rows.Next() {
        co +=1
    }
    var gacharesponse GachaDrawResponse
    var gachatableid int
    rand.Seed(time.Now().UnixNano())
    for i :=0;i<times;i++{
        gachatableid=rand.Intn(co)+1
        print(gachatableid)
        
        row1, err := db.Query(fmt.Sprintf("SELECT characterid FROM gachatable WHERE id = %d ; ",gachatableid))
        if err != nil {
            log.Fatal(err)
         }
        var characterid string
        for row1.Next() {        
            error :=row1.Scan(&characterid)
            if error != nil {
                panic(error.Error())
            }
        }
        
        row2,err:=db.Query(fmt.Sprintf("SELECT * FROM characters WHERE characterid = '%s' ; ",characterid))
        var result Gacharesult
        for row2.Next() {     
            error :=row2.Scan(&result.CharacterID,&result.Name)
            if error != nil {
             panic(error.Error())
            }
        }
        
        var newusercharacterid string
        newusercharacterid =RandString(8)


        ins, err := db.Prepare(fmt.Sprintf("INSERT INTO usercharacter(usercharacterid,characterid,usertoken) VALUES(?,?,?)"))
        if err != nil {
          log.Fatal(err)
         }
      
        ins.Exec(newusercharacterid,result.CharacterID,xtoken)



        gacharesponse.Results=append(gacharesponse.Results,result)

    }
    defer db.Close()

    res, err := json.Marshal(gacharesponse)
    
    w.Header().Set("Content-Type", "application/json")
    w.Write(res)

}

func characterList(w http.ResponseWriter, r *http.Request){
    xtoken := r.Header.Get("x-token")

    db, err := sql.Open("mysql", "root@/lesson1")
    log.Println("Connected to mysql.")
    //接続でエラーが発生した場合の処理
    if err != nil {
        log.Fatal(err)
    }

    var characterlistresponse CharacterListResponse

    row1, err := db.Query(fmt.Sprintf("SELECT usercharacterid , characterid FROM usercharacter WHERE usertoken = '%s' ; ",xtoken))
    if err != nil {
        log.Fatal(err)
        }

    for row1.Next() {
        var usercharacter UserCharacter
        error2 :=row1.Scan(&usercharacter.UserCharacterID,&usercharacter.CharacterID)
        if error2 != nil {
            panic(error2.Error())
        }
        row2, err := db.Query(fmt.Sprintf("SELECT name FROM characters WHERE characterid = '%s' ; ",usercharacter.CharacterID))
        if err != nil {
            log.Fatal(err)
        }
        for row2.Next() {     
            error3 :=row2.Scan(&usercharacter.Name)
            if error3 != nil {
             panic(error3.Error())
            }
        }

        characterlistresponse.Characters=append(characterlistresponse.Characters,usercharacter)

    }
    defer db.Close()

    res, err := json.Marshal(characterlistresponse)
    
    w.Header().Set("Content-Type", "application/json")
    w.Write(res)

}

func main() {    
    //req, _ := http.NewRequest("GET","localhost:8080" , nil)
    //req.Header.Set("Content-Type", "application/json")
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
