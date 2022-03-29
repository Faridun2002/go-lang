package types

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