package pkg

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	Type "regist/types"	
)

var (
	viv_user = []Type.User{}
)

var (
	DB *sql.DB
	err error
)

func ConnectDB() {
	DB, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	ForError(err)
	fmt.Println("DB connected!")
}

func SelectBlockUser() {
	var rows *sql.Rows
	rows, err := DB.Query("SELECT * FROM user WHERE id > 0 ORDER BY id")
	ForError(err)
	rows.Close()

	viv_user = []Type.User{}

	for rows.Next(){
        p := Type.User{}
        err = rows.Scan(&p.Id, &p.User_name, &p.User_surname, &p.Login, &p.Email, &p.Password)
        if err != nil{
			fmt.Println(err)
            continue
        }
        viv_user = append(viv_user, p)
    }
	DB.Close()
}

func SelectFunBlock(Input string) {
	search, err := DB.Query(fmt.Sprintf("SELECT id, login FROM user WHERE login='%s' AND id>0", Input))
	ForError(err)
	defer search.Close()

	var user Type.Block_user
	for search.Next() {
		err = search.Scan(&user.Id, &user.Login)
		ForError(err)
	}
}