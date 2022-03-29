package admin

import (
	"fmt"
	"html/template"
	"net/http"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	pkg "regist/pkg"
	Type "regist/types"
)

var (
	viv_user = []Type.User{}
	DB *sql.DB
)

func ConnectDB() {
	
}

func Admin(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/admin/index_admin.html", "templates/foot.html")
	pkg.ForError(err)

	DB, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	pkg.ForError(err)	
	
	rows, err := DB.Query("SELECT * FROM user WHERE id > 0 ORDER BY id")
    pkg.ForError(err)
    defer rows.Close()

	viv_user = []Type.User{}
	
	for rows.Next(){
        var p = Type.User{}
        err := rows.Scan(&p.Id, &p.User_name, &p.User_surname, &p.Login, &p.Email, &p.Password)
        if err != nil{
            fmt.Println("Цикл завершился ошибкой")
            continue
        }
        viv_user = append(viv_user, p)
    }
	DB.Close()
	t.ExecuteTemplate(w, "admin", viv_user)
}
