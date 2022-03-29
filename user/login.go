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

func Login_user(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/login.html", "templates/foot.html")
	pkg.ForError(err)
	t.ExecuteTemplate(w, "login", nil)
}

func Login_now(w http.ResponseWriter, r *http.Request) {
	Login = r.FormValue("login2")
	Password = r.FormValue("password2")
	
	if Login == "" || Password == "" {
		fmt.Fprintf(w, "Логин или пароль не введён")
	} else {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	pkg.ForError(err)
	
	defer db.Close()
	
	search, err := db.Query(fmt.Sprintf("SELECT id, name, login, password FROM user WHERE login='%s' AND password='%s'", Login, Password))
	pkg.ForError(err)
	defer search.Close()

	var user Type.User
	for search.Next() {
		err = search.Scan(&user.Id ,&user.User_name, &user.Login, &user.Password)
		pkg.ForError(err)
	}
	fmt.Println(user)
	
	if Login == user.Login && Password == user.Password {
		if user.Id < 0 {
			fmt.Fprintf(w, "Ваш аккаунт заблокирован. Вы можете обратиться к администратору по электронной почте helpAdmin@gmail.com")
			} else {
			update, err := db.Query(fmt.Sprintf("UPDATE active_user SET name='%s' WHERE id=1", user.User_name))
			pkg.ForError(err)
			defer update.Close()
			if Login == "Admin" {
				http.Redirect(w, r, "/admin", http.StatusSeeOther)
			} else {
				http.Redirect(w, r, "/main_site", http.StatusSeeOther)
			}
		}
		} else {
			fmt.Fprintf(w, "Логин или пароль введён неверно")
	}
}
}