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

func Recovery_user(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/admin/recovery_user.html", "templates/foot.html")
	pkg.ForError(err)
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	pkg.ForError(err)
	defer db.Close()
	
	rows, err := db.Query("SELECT * FROM user WHERE id < 0 ORDER BY id")
    pkg.ForError(err)
    defer rows.Close()

	viv_user = []Type.User{}
	
	for rows.Next(){
        var p = Type.User{}
        err := rows.Scan(&p.Id, &p.User_name, &p.User_surname, &p.Login, &p.Email, &p.Password)
        if err != nil{
            fmt.Println(err)
            continue
        }
        viv_user = append(viv_user, p)
    }

	t.ExecuteTemplate(w, "recovery_user", viv_user)
}

func FunRecovery(w http.ResponseWriter, r *http.Request) {
	Input = r.FormValue("recovery_input")
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	pkg.ForError(err)
	defer db.Close()

	search, err := db.Query(fmt.Sprintf("SELECT id, login FROM user WHERE login='%s' AND id<0", Input))
	pkg.ForError(err)
	defer search.Close()

	var user Type.Block_user
	for search.Next() {
		err = search.Scan(&user.Id, &user.Login)
		pkg.ForError(err)
	}

	ID := -user.Id 
	update, err := db.Query(fmt.Sprintf("UPDATE user SET id='%d' WHERE login='%s' AND id<0", ID, Input))
	pkg.ForError(err)
	defer update.Close()

	fmt.Println(ID)

	http.Redirect(w, r, "/recovery_user", http.StatusSeeOther)
}
