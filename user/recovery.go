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

var (
	Rec_ac string
	user Type.Rec_user
)

func Recovery(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/recovery.html", "templates/foot.html" )
	pkg.ForError(err)
	t.ExecuteTemplate(w, "recovery", nil)
}

func Recovery_now(w http.ResponseWriter, r *http.Request) {
	Rec_ac = r.FormValue("rec_ac")

	if Rec_ac == "" {
		fmt.Fprintf(w, "Не все данные заполнены")
	} else {
		db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
		pkg.ForError(err)
		defer db.Close()
		
		search, err := db.Query(fmt.Sprintf("SELECT login, email, password FROM user WHERE login='%s' OR email='%s'", Rec_ac, Rec_ac))
		pkg.ForError(err)
		defer search.Close()
		
		for search.Next() {
			err = search.Scan(&user.Login, &user.Email, &user.Pass)
			pkg.ForError(err)
		}
		
		if Rec_ac == user.Login || Rec_ac == user.Email {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			pkg.SendPassword(user)
			} else {
				fmt.Println("Пользователь с таким адресом почты не зарегистрировался")
		}
	}

	fmt.Println("Server is listing...")
}