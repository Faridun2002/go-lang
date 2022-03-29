package pkg

import (
	"net/http"
	User "regist/user"
	Admin "regist/admin"
	Main "regist/site"
)

func HandleFunction() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.Handle("/templates/excel/", http.StripPrefix("/templates/excel/",  http.FileServer(http.Dir("./templates/excel/"))))
	http.HandleFunc("/main_site", Main.Main_site)
	http.HandleFunc("/", Main.Index)
	http.HandleFunc("/create", User.Create)
	http.HandleFunc("/save_create", User.Save_create)
	http.HandleFunc("/verification", User.Verification)
	http.HandleFunc("/save_ver", User.Save_ver)
	http.HandleFunc("/login", User.Login_user)
	http.HandleFunc("/login_now", User.Login_now)
	http.HandleFunc("/recovery", User.Recovery)
	http.HandleFunc("/recovery_now", User.Recovery_now)
	http.HandleFunc("/admin", Admin.Admin)
	http.HandleFunc("/block_user", Admin.Block_user)
	http.HandleFunc("/funBlock", Admin.FunBlock)
	http.HandleFunc("/recovery_user", Admin.Recovery_user)
	http.HandleFunc("/funRecovery", Admin.FunRecovery)
	http.HandleFunc("/delete_user", Admin.Delete_user)
	http.HandleFunc("/funDelete", Admin.FunDelete)
	http.HandleFunc("/insert_user", Admin.Insert_user)
	http.HandleFunc("/funInsert", Admin.FunInsert)
	http.HandleFunc("/update", Admin.Update)
	http.HandleFunc("/funUpdate", Admin.FunUpdate)
	http.HandleFunc("/update_user", Admin.Update_user)
	http.HandleFunc("/funUpdateUser", Admin.FunUpdateUser)
	http.HandleFunc("/save_excel", Admin.Save_excel)
	http.HandleFunc("/send_excel", Admin.Send_excel)
	http.ListenAndServe(":8080", nil)
}