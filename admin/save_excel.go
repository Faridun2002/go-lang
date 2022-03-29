package admin

import (
	"fmt"
	"time"
	"math/rand"
	"strconv"
	"html/template"
	"net/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	Type "regist/types"
	pkg "regist/pkg"
	"github.com/xuri/excelize/v2"
)

var (
	randNum int
	num int = 1
	E string
	sNum string
)

func Save_excel(w http.ResponseWriter, r *http.Request) {
	Time_now = time.Now()
	fmt.Println(Time_now.Format("2006.01.02"))
	rand.Seed(time.Now().UTC().UnixNano())
	randNum = rand.Intn(9999 - 1000)

	E = strconv.Itoa(randNum)
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	pkg.ForError(err)
	defer db.Close()

	rows, err := db.Query("SELECT * FROM user WHERE id != 0 ORDER BY id")
	pkg.ForError(err)
	defer rows.Close()

	f := excelize.NewFile()
	f.SetCellValue("Sheet1", "A1", "ID")
	f.SetCellValue("Sheet1", "B1", "Name")
	f.SetCellValue("Sheet1", "C1", "Surname")
	f.SetCellValue("Sheet1", "D1", "Login")
	f.SetCellValue("Sheet1", "E1", "Email")
	f.SetCellValue("Sheet1", "F1", "Password")

	for rows.Next(){
		var p = Type.User{}
		num += 1
		sNum = strconv.Itoa(num)
		err := rows.Scan(&p.Id, &p.User_name, &p.User_surname, &p.Login, &p.Email, &p.Password)
		if err != nil{
			fmt.Println(err)
			continue
		}
		f.SetCellValue("Sheet1", "A" + sNum, p.Id)
		f.SetCellValue("Sheet1", "B" + sNum, p.User_name)
		f.SetCellValue("Sheet1", "C" + sNum, p.User_surname)
		f.SetCellValue("Sheet1", "D" + sNum, p.Login)
		f.SetCellValue("Sheet1", "E" + sNum, p.Email)
		f.SetCellValue("Sheet1", "F" + sNum, p.Password)
	}
	if err := f.SaveAs("templates/excel/" + Time_now.Format("2006.01.02") + ".xlsx"); err != nil {
		fmt.Println("Не удалось сохранить файл")
	} else {
		fmt.Println("Отчет сохранен")
	}
	Send()
	http.Redirect(w, r, "/send_excel", http.StatusSeeOther)
}

func Send_excel(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/admin/send_excel.html", "templates/foot.html")
	pkg.ForError(err)
	t.ExecuteTemplate(w, "send_excel", nil)
}