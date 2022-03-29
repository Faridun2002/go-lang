package user

import (
	"fmt"
	"net/http"
	"html/template"
	"math/rand"
	"time"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	Type "regist/types"
	pkg "regist/pkg"
)

var(
	User_name string
	User_surname string
	Login string
	Email string
	Password string
	Bytes int
)

func Create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html", "templates/foot.html")
	pkg.ForError(err)
	t.ExecuteTemplate(w, "create", nil)
}

func Save_create(w http.ResponseWriter, r *http.Request) {
	User_name = r.FormValue("user_name")
	User_surname = r.FormValue("user_surname")
	Login = r.FormValue("login")
	Email = r.FormValue("email")
	Password = r.FormValue("password")
	
	rand.Seed(time.Now().UTC().UnixNano())
	Bytes = rand.Intn(9999 - 1000)
	if Bytes < 1000 {
		Bytes += 1000
	}
	if User_name == "" || User_surname == "" || Login == "" || Email == "" || Password == "" {
		fmt.Fprintf(w, "Не все данные заполнены")
		} else {
			db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
			pkg.ForError(err)

			search, err := db.Query(fmt.Sprintf("SELECT login, email FROM user WHERE login='%s' AND email='%s'", Login, Email))
			pkg.ForError(err)
			defer search.Close()
	
			var user Type.User
	for search.Next() {
		err = search.Scan(&user.Login, &user.Email)
		pkg.ForError(err)
	}
	
	if user.Login == Login || user.Email == Email {
		fmt.Fprintf(w, "Пользователь с таким логином или почтой уже зарегистрирован")
		} else {
			insert_ver, err := db.Query(fmt.Sprintf("INSERT INTO verification (name, pass) VALUES('%s', '%d')", Login, Bytes))
			pkg.ForError(err)
			defer insert_ver.Close()
			pkg.Mail(Bytes, Email)
			fmt.Println("Пароль отправлен на почту")
			http.Redirect(w, r, "/verification", http.StatusSeeOther)
		}
		
		http.Redirect(w, r, "/", http.StatusSeeOther)
		fmt.Println("Server is listing...")
	}
}