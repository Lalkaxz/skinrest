package models

type User struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserData struct {
	Id       int
	Login    string
	Password string
	Token    string
}

type UserInfo struct {
	Login string
	Skins []SkinData
}
