package user

import (
	"fmt"
	"net/http"
	"html/template"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	Type "regist/types"
	pkg "regist/pkg"
)

func Verification(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/verification.html", "templates/foot.html")
	pkg.ForError(err)
	t.ExecuteTemplate(w, "verification", nil)
}

func Save_ver(w http.ResponseWriter, r *http.Request) {
	Password = r.FormValue("password3")
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	pkg.ForError(err)
	defer db.Close()
	
	search, err := db.Query(fmt.Sprintf("SELECT name, pass FROM verification WHERE name='%s'", Login))
	pkg.ForError(err)
	defer search.Close()
	
	var user Type.Ver_user
	for search.Next() {
		err = search.Scan(&user.Name, &user.Pass)
		pkg.ForError(err)
	}
	
	if Login == user.Name && Password == user.Pass {
		http.Redirect(w, r, "/main_site", http.StatusSeeOther)
		insert, err := db.Query(fmt.Sprintf("INSERT INTO user (name, surname, login, email, password) VALUES('%s', '%s', '%s', '%s', '%s')", User_name, User_surname, Login, Email, Password))
		pkg.ForError(err)
		defer insert.Close()
		update, err := db.Query(fmt.Sprintf("UPDATE active_user SET name='%s' WHERE id=1", User_name))
		pkg.ForError(err)
		defer update.Close()
		} else {
			fmt.Println("Неверный пароль")
		}		
}