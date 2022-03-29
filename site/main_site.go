package site

import (
	"fmt"
	"html/template"
	"net/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	Type "regist/types"
	pkg "regist/pkg"
)

var (active_user = []Type.Active_user{})

func Index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/foot.html")
	pkg.ForError(err)
	t.ExecuteTemplate(w, "index", nil)
}

func Main_site(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/main_site.html", "templates/foot.html")
	pkg.ForError(err)
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	pkg.ForError(err)
	defer db.Close()
	
	rows, err := db.Query("SELECT name FROM golang.active_user")
    pkg.ForError(err)
    defer rows.Close()

	active_user = []Type.Active_user{}
	
	for rows.Next(){
        var p = Type.Active_user{}
        err := rows.Scan(&p.Name)
        if err != nil{
            fmt.Println(err)
            continue
        }
        active_user = append(active_user, p)
    }

	fmt.Println(active_user)
	t.ExecuteTemplate(w, "main_site", active_user)
}
