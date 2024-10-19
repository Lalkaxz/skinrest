package models

type Skin struct {
	Name string `json:"skinname" binding:"required"`
	Type string `json:"skintype" binding:"required"`
	Src  string `json:"skinsrc" binding:"required"`
}

type SkinData struct {
	Id   int
	Name string
	Type string
	Src  string
}
