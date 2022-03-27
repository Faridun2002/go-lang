package main

import (
	"bytes"
	// "log"
	"fmt"
	"html/template"
	"net/http"
	"math/rand"
	"strconv"
	"time"
	"io/ioutil"
	"path/filepath"
	"encoding/base64"
	"strings"
	"net/smtp"
	"mime/multipart"
	
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	_ "github.com/golang/protobuf/ptypes/struct"
	_ "github.com/gorilla/sessions"
	"github.com/xuri/excelize/v2"
	gomail "gopkg.in/mail.v2"
	// Type "./types"
	// pkg "./pkg"
)

type User struct {
	Id 			int16
	User_name 	string
	User_surname string 
	Login		string
	Email		string
	Password string
}
type Active_user struct {
	Name 		string
}
type Ver_user struct {
	Name		string
	Pass 		string
}
type Block_user struct {
	Id 			int16
	Login 		string
}
type Rec_user struct {
	Login 		string
	Email 		string
	Pass 		string
}
type Sender struct {
	Auth smtp.Auth
}
type Message struct {
	To          []string
	Subject     string
	Body        string
	Attachments map[string][]byte
}

var (
	viv_user = []User{}
	active_user = []Active_user{}
	User_name string
	User_surname string
	Login string
	Email string
	Password string
	Bytes int
	Rec_ac string
	user Rec_user
	Input string
	InputCheck []string
	Block string 
	Insert string 
	Delete string
	UpdateUserName string 
	num int = 1
	sNum string
	randNum int
	E string
	time_now time.Time
	host 		= "smtp.gmail.com"
	username 	= "faridunjalolov1@gmail.com"
	password 	= "farn2212"
	portNumber 	= "587"
)

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/foot.html")
	if err != nil {
		panic(err)
	}
	t.ExecuteTemplate(w, "index", nil)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html", "templates/foot.html")
	if err != nil {
		panic(err)
	}
	t.ExecuteTemplate(w, "create", nil)
}

func save_create(w http.ResponseWriter, r *http.Request) {
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
			if err != nil {
				panic("Не удалось подключиться к БД")
			}
			fmt.Println("DB connected!")
			defer db.Close()
			
			search, err := db.Query(fmt.Sprintf("SELECT login, email FROM user WHERE login='%s' AND email='%s'", Login, Email))
			if err != nil {
				panic(err)
			}
			defer search.Close()
	
	var user User
	for search.Next() {
		err = search.Scan(&user.Login, &user.Email)
		if err != nil {
			panic(err)
		}
	}
	
	if user.Login == Login || user.Email == Email {
		fmt.Fprintf(w, "Пользователь с таким логином или почтой уже зарегистрирован")
		} else {
			insert_ver, err := db.Query(fmt.Sprintf("INSERT INTO verification (name, pass) VALUES('%s', '%d')", Login, Bytes))
			if err != nil {
				panic(err)
			}
			defer insert_ver.Close()
			mail()
			fmt.Println("Пароль отправлен на почту")
			http.Redirect(w, r, "/verification", http.StatusSeeOther)
		}
		
		http.Redirect(w, r, "/", http.StatusSeeOther)
	fmt.Println("Server is listing...")
}
}

func mail() {
	B := strconv.Itoa(Bytes)
	mail := gomail.NewMessage()
	mail.SetHeader("From", "faridunjalolov1@gmail.com")
	mail.SetHeader("To", Email)
	mail.SetHeader("Subject", "Пароль для подтверждения регистрации")
	mail.SetBody("text/html", "<body>Введите " + B + " этот пароль для подтверждение регистрации</body>")
	a := gomail.NewDialer("smtp.gmail.com",587,"faridunjalolov1@gmail.com","farn2212")
	if err := a.DialAndSend(mail); err != nil {
		fmt.Println()
		panic(err)
	}
	fmt.Println(Bytes)
}

func verification(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/verification.html", "templates/foot.html")
	if err != nil {
		panic(err)
	}
	t.ExecuteTemplate(w, "verification", nil)
}

func save_ver(w http.ResponseWriter, r *http.Request) {
	Password = r.FormValue("password3")
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic("Не удалось подключиться к БД")
	}
	fmt.Println("DB connected!")
	defer db.Close()
	
	search, err := db.Query(fmt.Sprintf("SELECT name, pass FROM verification WHERE name='%s'", Login))
	if err != nil {
		panic(err)
	}
	defer search.Close()

	var user Ver_user
	for search.Next() {
		err = search.Scan(&user.Name, &user.Pass)
		if err != nil {
			panic(err)
		} 
	}

	if Login == user.Name && Password == user.Pass {
		http.Redirect(w, r, "/main_site", http.StatusSeeOther)
		insert, err := db.Query(fmt.Sprintf("INSERT INTO user (name, surname, login, email, password) VALUES('%s', '%s', '%s', '%s', '%s')", User_name, User_surname, Login, Email, Password))
		if err != nil {
		panic(err)
		}
		defer insert.Close()
		update, err := db.Query(fmt.Sprintf("UPDATE active_user SET name='%s' WHERE id=1", User_name))
		if err != nil {
			panic(err)
		}
		defer update.Close()
		} else {
		fmt.Println("Неверный пароль")
	}

}

func login(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/login.html", "templates/foot.html")
	if err != nil {
		panic(err)
	}
	t.ExecuteTemplate(w, "login", nil)
}

func login_now(w http.ResponseWriter, r *http.Request) {
	Login = r.FormValue("login2")
	Password = r.FormValue("password2")
	
	if Login == "" || Password == "" {
		fmt.Fprintf(w, "Логин или пароль не введён")
	} else {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic("Не удалось подключиться к БД")
	}
	fmt.Println("DB connected!")
	defer db.Close()
	
	search, err := db.Query(fmt.Sprintf("SELECT id, name, login, password FROM user WHERE login='%s' AND password='%s'", Login, Password))
	if err != nil {
		panic(err)
	}
	defer search.Close()

	var user User
	for search.Next() {
		err = search.Scan(&user.Id ,&user.User_name, &user.Login, &user.Password)
		if err != nil {
			panic(err)
		} 
	}
	fmt.Println(user)
	
	if Login == user.Login && Password == user.Password {
		if user.Id < 0 {
			fmt.Fprintf(w, "Ваш аккаунт заблокирован. Вы можете обратиться к администратору по электронной почте helpAdmin@gmail.com")
			} else {
			update, err := db.Query(fmt.Sprintf("UPDATE active_user SET name='%s' WHERE id=1", user.User_name))
			if err != nil {
				panic(err)
			}
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

func admin(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/admin/index_admin.html", "templates/foot.html")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic("Не удалось подключиться к БД")
	}
	fmt.Println("DB connected!")
	defer db.Close()
	
	rows, err := db.Query("SELECT * FROM user WHERE id > 0 ORDER BY id")
    if err != nil {
        panic(err)
    }
    defer rows.Close()

	viv_user = []User{}
	
	for rows.Next(){
        var p = User{}
        err := rows.Scan(&p.Id, &p.User_name, &p.User_surname, &p.Login, &p.Email, &p.Password)
        if err != nil{
            fmt.Println(err)
            continue
        }
        viv_user = append(viv_user, p)
    }

	t.ExecuteTemplate(w, "admin", viv_user)
}

func block_user(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/admin/block_user.html", "templates/foot.html")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	
	rows, err := db.Query("SELECT * FROM user WHERE id > 0 ORDER BY id")
    if err != nil {
        panic(err)
    }
    defer rows.Close()

	viv_user = []User{}
	
	for rows.Next(){
        var p = User{}
        err := rows.Scan(&p.Id, &p.User_name, &p.User_surname, &p.Login, &p.Email, &p.Password)
        if err != nil{
            fmt.Println(err)
            continue
        }
        viv_user = append(viv_user, p)
    }

	t.ExecuteTemplate(w, "block_user", viv_user)
}

// func contains(slice []string, item string) bool {
//     set := make(map[string]struct{}, len(slice))
//     for _, s := range slice {
//         set[s] = struct{}{}
//     }
//     _, ok := set[item]
//     return ok
// }

// func ProcessCheckboxes(w http.ResponseWriter, r *http.Request) {
//     r.ParseForm()
//     fmt.Printf("%+v\n", r.Form)
//     productsSelected := r.Form["product_image_id"]
//     log.Println(contains(productsSelected, "Grape"))
// }
func funBlock(w http.ResponseWriter, r *http.Request) {
	Input = r.FormValue("block_input")
	// r.ParseForm()
	// InputCheck = r.Form["checkbox_login"]
	// for _, i := range InputCheck {
	// 	fmt.Println(i)
	// }

    db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	search, err := db.Query(fmt.Sprintf("SELECT id, login FROM user WHERE login='%s' AND id>0", Input))
	if err != nil {
		panic(err)
	}
	defer search.Close()

	var user Block_user
	for search.Next() {
		err = search.Scan(&user.Id, &user.Login)
		if err != nil {
			panic(err)
		} 
	}

	ID := -user.Id 
	update, err := db.Query(fmt.Sprintf("UPDATE user SET id='%d' WHERE login='%s' AND id>0", ID, Input))
		if err != nil {
			panic(err)
		}
	defer update.Close()

	http.Redirect(w, r, "/block_user", http.StatusSeeOther)
}

// func funBlock(w http.ResponseWriter, r *http.Request) {
// 	Input = r.FormValue("block_input")
// 	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer db.Close()
//
// 	search, err := db.Query(fmt.Sprintf("SELECT id, login FROM user WHERE login='%s' AND id>0", Input))
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer search.Close()
//
// 	var user Block_user
// 	for search.Next() {
// 		err = search.Scan(&user.Id, &user.Login)
// 		if err != nil {
// 			panic(err)
// 		} 
// 	}
//
// 	ID := -user.Id 
// 	update, err := db.Query(fmt.Sprintf("UPDATE user SET id='%d' WHERE login='%s' AND id>0", ID, Input))
// 		if err != nil {
// 			panic(err)
// 		}
// 	defer update.Close()
//
// 	fmt.Println(ID)
//
// 	http.Redirect(w, r, "/block_user", http.StatusSeeOther)
// }

func recovery_user(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/admin/recovery_user.html", "templates/foot.html")
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	
	rows, err := db.Query("SELECT * FROM user WHERE id < 0 ORDER BY id")
    if err != nil {
        panic(err)
    }
    defer rows.Close()

	viv_user = []User{}
	
	for rows.Next(){
        var p = User{}
        err := rows.Scan(&p.Id, &p.User_name, &p.User_surname, &p.Login, &p.Email, &p.Password)
        if err != nil{
            fmt.Println(err)
            continue
        }
        viv_user = append(viv_user, p)
    }

	// fmt.Println(viv_user)
	t.ExecuteTemplate(w, "recovery_user", viv_user)
}

func funRecovery(w http.ResponseWriter, r *http.Request) {
	Input = r.FormValue("recovery_input")
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	search, err := db.Query(fmt.Sprintf("SELECT id, login FROM user WHERE login='%s' AND id<0", Input))
	if err != nil {
		panic(err)
	}
	defer search.Close()

	var user Block_user
	for search.Next() {
		err = search.Scan(&user.Id, &user.Login)
		if err != nil {
			panic(err)
		} 
	}

	ID := -user.Id 
	update, err := db.Query(fmt.Sprintf("UPDATE user SET id='%d' WHERE login='%s' AND id<0", ID, Input))
		if err != nil {
			panic(err)
		}
	defer update.Close()

	fmt.Println(ID)

	http.Redirect(w, r, "/recovery_user", http.StatusSeeOther)
}

func delete_user(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/admin/delete_user.html", "templates/foot.html")
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	
	rows, err := db.Query("SELECT * FROM user WHERE id!=0 ORDER BY id")
    if err != nil {
        panic(err)
    }
    defer rows.Close()

	viv_user = []User{}
	
	for rows.Next(){
        var p = User{}
        err := rows.Scan(&p.Id, &p.User_name, &p.User_surname, &p.Login, &p.Email, &p.Password)
        if err != nil{
            fmt.Println(err)
            continue
        }
        viv_user = append(viv_user, p)
    }

	// fmt.Println(viv_user)
	t.ExecuteTemplate(w, "delete_user", viv_user)
}

func funDelete(w http.ResponseWriter, r *http.Request) {
	Input = r.FormValue("delete_input")
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	search, err := db.Query(fmt.Sprintf("SELECT id, login FROM user WHERE login='%s' AND id!=0", Input))
	if err != nil {
		panic(err)
	}
	defer search.Close()

	var user Block_user
	for search.Next() {
		err = search.Scan(&user.Id, &user.Login)
		if err != nil {
			panic(err)
		} 
	}

	delete, err := db.Query(fmt.Sprintf("DELETE FROM user WHERE login='%s' AND id!=0", user.Login))
		if err != nil {
			panic(err)
		}
	defer delete.Close()

	http.Redirect(w, r, "/delete_user", http.StatusSeeOther)
}

func insert_user(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/admin/insert_user.html", "templates/foot.html")
	if err != nil {
		panic(err)
	}
	t.ExecuteTemplate(w, "insert_user", nil)
}

func funInsert(w http.ResponseWriter, r *http.Request) {
	Name_user := r.FormValue("name_user")
	Surname_user := r.FormValue("surname_user")
	Login_user := r.FormValue("login_user")
	Email_user := r.FormValue("email_user")
	Password_user := r.FormValue("password_user")

	if Name_user == "" || Surname_user == "" || Login_user == "" || Email_user == "" || Password_user == "" {
		fmt.Fprintf(w, "Не все данные заполнены")
	} else {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	
	search, err := db.Query(fmt.Sprintf("SELECT login, email FROM user WHERE login='%s' AND email='%s'", Login_user, Email_user))
	if err != nil {
		panic(err)
	}
	defer search.Close()

	var user User
	for search.Next() {
		err = search.Scan(&user.Login, &user.Email)
		if err != nil {
			panic(err)
		}
	}
	
	if Login_user == user.Login || Email_user == user.Email {
				fmt.Fprintf(w, "Пользователь с таким логином или почтой уже зарегистрирован")
	} else {
		insert, err := db.Query(fmt.Sprintf("INSERT INTO user (name, surname, login, email, password) VALUES('%s', '%s', '%s', '%s', '%s')", Name_user, Surname_user, Login_user, Email_user, Password_user))
		if err != nil {
		panic(err)
		}
		defer insert.Close()
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	fmt.Println("Server is listing...")
}
}

func update(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/admin/update.html", "templates/foot.html")
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	
	rows, err := db.Query("SELECT * FROM user WHERE id!=0 ORDER BY id")
    if err != nil {
        panic(err)
    }
    defer rows.Close()

	viv_user = []User{}
	
	for rows.Next(){
        var p = User{}
        err := rows.Scan(&p.Id, &p.User_name, &p.User_surname, &p.Login, &p.Email, &p.Password)
        if err != nil{
            fmt.Println(err)
            continue
        }
        viv_user = append(viv_user, p)
    }

	// fmt.Println(viv_user)
	t.ExecuteTemplate(w, "update", viv_user)
}

func funUpdate(w http.ResponseWriter, r *http.Request) {
	Input = r.FormValue("update_input")
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	search, err := db.Query(fmt.Sprintf("SELECT id, login FROM user WHERE login='%s' AND id!=0", Input))
	if err != nil {
		panic(err)
	}
	defer search.Close()

	var user Block_user
	for search.Next() {
		err = search.Scan(&user.Id, &user.Login)
		if err != nil {
			panic(err)
		} 
	}

	UpdateUserName = user.Login
	http.Redirect(w, r, "/update_user", http.StatusSeeOther)
}

func update_user(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/admin/update_user.html", "templates/foot.html")
	if err != nil {
		panic(err)
	}
	
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	
	rows, err := db.Query(fmt.Sprintf("SELECT name, surname, login, email, password FROM user WHERE login='%s' AND id!=0 LIMIT 1", UpdateUserName))
    if err != nil {
        panic(err)
    }
    defer rows.Close()

	viv_user = []User{}
	
	for rows.Next(){
        var p = User{}
        err := rows.Scan(&p.User_name, &p.User_surname, &p.Login, &p.Email, &p.Password)
        if err != nil{
            fmt.Println(err)
            continue
        }
        viv_user = append(viv_user, p)
    }
	t.ExecuteTemplate(w, "update_user", viv_user)
}

func funUpdateUser(w http.ResponseWriter, r *http.Request) {
	
	Name_user := r.FormValue("name_user_update")
	Surname_user := r.FormValue("surname_user_update")
	Login_user := r.FormValue("login_user_update")
	Email_user := r.FormValue("email_user_update")
	Password_user := r.FormValue("password_user_update")

	if Name_user == "" || Surname_user == "" || Login_user == "" || Email_user == "" || Password_user == "" {
		fmt.Fprintf(w, "Не все данные заполнены")
	} else {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	
	search, err := db.Query(fmt.Sprintf("SELECT login, email FROM user WHERE login='%s' AND email='%s' LIMIT 1", Login_user, Email_user))
	if err != nil {
		panic(err)
	}
	defer search.Close()

	var user User
	for search.Next() {
		err = search.Scan(&user.Login, &user.Email)
		if err != nil {
			panic(err)
		}
	}
	
		insert, err := db.Query(fmt.Sprintf("UPDATE user SET name='%s', surname='%s', login='%s', email='%s', password='%s' WHERE id!=0 AND login='%s' LIMIT 1", Name_user, Surname_user, Login_user, Email_user, Password_user, UpdateUserName))
		if err != nil {
		panic(err)
		}
		defer insert.Close()
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	

	http.Redirect(w, r, "/", http.StatusSeeOther)
	fmt.Println("Server is listing...")
}
}

func recovery(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/recovery.html", "templates/foot.html" )
	if err != nil {
		panic(err)
	}
	t.ExecuteTemplate(w, "recovery", nil)
}

func recovery_now(w http.ResponseWriter, r *http.Request) {
	Rec_ac = r.FormValue("rec_ac")

	if Rec_ac == "" {
		fmt.Fprintf(w, "Не все данные заполнены")
	} else {
		db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
		if err != nil {
			panic(err)
		}
		defer db.Close()
		
		search, err := db.Query(fmt.Sprintf("SELECT login, email, password FROM user WHERE login='%s' OR email='%s'", Rec_ac, Rec_ac))
		if err != nil {
			panic(err)
		}
		defer search.Close()
		
		for search.Next() {
			err = search.Scan(&user.Login, &user.Email, &user.Pass)
			if err != nil {
				panic(err)
			}
		}
		
		if Rec_ac == user.Login || Rec_ac == user.Email {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			sendPassword()
			} else {
				fmt.Println("Пользователь с таким адресом почты не зарегистрировался")
		}
	}

	fmt.Println("Server is listing...")
}

func sendPassword() {
	mail := gomail.NewMessage()
	mail.SetHeader("From", "faridunjalolov1@gmail.com")
	mail.SetHeader("To", user.Email)
	mail.SetHeader("Subject", "Пароль для входа")
	mail.SetBody("text/html", "<html><body><b>Ваш логин:</b> " + user.Login + "<br><b>Ваш пароль:</b> " + user.Pass + "</body></html>")
	a := gomail.NewDialer("smtp.gmail.com",587,"faridunjalolov1@gmail.com","farn2212")
	if err := a.DialAndSend(mail); err != nil {
		fmt.Println()
		panic(err)
	}
}

func New() *Sender {
	Auth := smtp.PlainAuth("", username, password, host)
	return &Sender{Auth}
}

func (s *Sender) Send(m *Message) error {
	return smtp.SendMail(fmt.Sprintf("%s:%s", host, portNumber), s.Auth, username, m.To, m.ToBytes())
}

func NewMessage(s, b string) *Message {
	return &Message{Subject: s, Body: b, Attachments: make(map[string][]byte)}
}

func (m *Message) AttachFile(src string) error {
	b, err := ioutil.ReadFile(src)
	if err != nil {
		panic(err)
	}
	_, fileName := filepath.Split(src)
	m.Attachments[fileName] = b
	return nil
}

func (m *Message) ToBytes() []byte {
	buf := bytes.NewBuffer(nil)
	withAttachments := len(m.Attachments) > 0
	buf.WriteString(fmt.Sprintf("Subject: %s\n", m.Subject))
	buf.WriteString(fmt.Sprintf("To: %s\n", strings.Join(m.To, ",")))
	buf.WriteString("MIME-Version: 1.0\n")
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()
	if withAttachments {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))
		buf.WriteString(fmt.Sprintf("--%s\n", boundary))
	} else {
		buf.WriteString("Content-Type: text/plain; charset=utf-8\n")
	}
	buf.WriteString(m.Body)
	if withAttachments {
		for k, v := range m.Attachments {
			buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType(v)))
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", k))

			b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)
			buf.Write(b)
			buf.WriteString(fmt.Sprintf("\n--%s", boundary))
		}

		buf.WriteString("--")
	}

	return buf.Bytes()
}

func Send() {
	sender := New()
	m := NewMessage("Test", "Body message.")
	m.To = []string{"faridunjalolov0@gmail.com"}
	m.AttachFile("templates/excel/" + time_now.Format("2006.01.02") + ".xlsx")
	fmt.Println(sender.Send(m))
}

func save_excel(w http.ResponseWriter, r *http.Request) {
time_now = time.Now()
fmt.Println(time_now.Format("2006.01.02"))
rand.Seed(time.Now().UTC().UnixNano())
randNum = rand.Intn(9999 - 1000)

E = strconv.Itoa(randNum)
db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
if err != nil {
panic(err)
}
defer db.Close()

rows, err := db.Query("SELECT * FROM user WHERE id != 0 ORDER BY id")
if err != nil {
panic(err)
}
defer rows.Close()

f := excelize.NewFile()
f.SetCellValue("Sheet1", "A1", "ID")
f.SetCellValue("Sheet1", "B1", "Name")
f.SetCellValue("Sheet1", "C1", "Surname")
f.SetCellValue("Sheet1", "D1", "Login")
f.SetCellValue("Sheet1", "E1", "Email")
f.SetCellValue("Sheet1", "F1", "Password")

for rows.Next(){
var p = User{}
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
if err := f.SaveAs("templates/excel/" + time_now.Format("2006.01.02") + ".xlsx"); err != nil {
fmt.Println(err)
}
fmt.Println("Отчет сохранен")
Send()
http.Redirect(w, r, "/send_excel", http.StatusSeeOther)
}

func send_excel(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/admin/send_excel.html", "templates/foot.html")
	if err != nil {
		panic(err)
	}
	E = strconv.Itoa(randNum)
	t.ExecuteTemplate(w, "send_excel", struct{E string}{E: E})
}

func main_site(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/main_site.html", "templates/foot.html")
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	
	rows, err := db.Query("SELECT name FROM golang.active_user")
    if err != nil {
        panic(err)
    }
    defer rows.Close()

	active_user = []Active_user{}
	
	for rows.Next(){
        var p = Active_user{}
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

func handleFunc() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.Handle("/templates/excel/", http.StripPrefix("/templates/excel/",  http.FileServer(http.Dir("./templates/excel/"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/create", create)
	http.HandleFunc("/save_create", save_create)
	http.HandleFunc("/verification", verification)
	http.HandleFunc("/save_ver", save_ver)
	http.HandleFunc("/login", login)
	http.HandleFunc("/login_now", login_now)
	http.HandleFunc("/admin", admin)
	http.HandleFunc("/block_user", block_user)
	http.HandleFunc("/funBlock", funBlock)
	http.HandleFunc("/recovery_user", recovery_user)
	http.HandleFunc("/funRecovery", funRecovery)
	http.HandleFunc("/delete_user", delete_user)
	http.HandleFunc("/funDelete", funDelete)
	http.HandleFunc("/insert_user", insert_user)
	http.HandleFunc("/funInsert", funInsert)
	http.HandleFunc("/update", update)
	http.HandleFunc("/funUpdate", funUpdate)
	http.HandleFunc("/update_user", update_user)
	http.HandleFunc("/funUpdateUser", funUpdateUser)
	http.HandleFunc("/recovery", recovery)
	http.HandleFunc("/recovery_now", recovery_now)
	http.HandleFunc("/save_excel", save_excel)
	http.HandleFunc("/send_excel", send_excel)
	http.HandleFunc("/main_site", main_site)
	http.ListenAndServe(":8080", nil)
}

func main() {
	handleFunc()
}
