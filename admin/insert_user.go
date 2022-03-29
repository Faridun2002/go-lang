package admin

import (
	"fmt"
	"net/http"
	"html/template"
	
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	Type "regist/types"
	pkg "regist/pkg"
)

func Insert_user(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/admin/insert_user.html", "templates/foot.html")
	pkg.ForError(err)
	t.ExecuteTemplate(w, "insert_user", nil)
}

func FunInsert(w http.ResponseWriter, r *http.Request) {
	Name_user := r.FormValue("name_user")
	Surname_user := r.FormValue("surname_user")
	Login_user := r.FormValue("login_user")
	Email_user := r.FormValue("email_user")
	Password_user := r.FormValue("password_user")

	if Name_user == "" || Surname_user == "" || Login_user == "" || Email_user == "" || Password_user == "" {
		fmt.Fprintf(w, "Не все данные заполнены")
	} else {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	pkg.ForError(err)
	defer db.Close()
	
	search, err := db.Query(fmt.Sprintf("SELECT login, email FROM user WHERE login='%s' AND email='%s'", Login_user, Email_user))
	pkg.ForError(err)
	defer search.Close()

	var user Type.User
	for search.Next() {
		err = search.Scan(&user.Login, &user.Email)
		pkg.ForError(err)
	}
	
	if Login_user == user.Login || Email_user == user.Email {
				fmt.Fprintf(w, "Пользователь с таким логином или почтой уже зарегистрирован")
	} else {
		insert, err := db.Query(fmt.Sprintf("INSERT INTO user (name, surname, login, email, password) VALUES('%s', '%s', '%s', '%s', '%s')", Name_user, Surname_user, Login_user, Email_user, Password_user))
		pkg.ForError(err)
		defer insert.Close()
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	fmt.Println("Server is listing...")
}
}