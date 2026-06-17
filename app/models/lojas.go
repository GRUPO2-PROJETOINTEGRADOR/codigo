package models

type Loja struct {
	ID        string `json:"id"`
	Nome      string `json:"nome"`
	Categoria string `json:"categoria"`
}

type LojaBusca struct {
	ID   string `json:"id"`
	Nome string `json:"nome"`
	LUC  string `json:"luc"`
}
