package models

import (
	"time"
)

type OrientacaoEducativa struct {
	ID                  int       `json:"id"`
	LojaID              string    `json:"loja_id"`
	ResponsavelPresente string    `json:"responsavel_presente"`
	FuncaoResponsavel   string    `json:"funcao_responsavel"`
	DataOrientacao      time.Time `json:"data_orientacao"`
	Observacoes         string    `json:"observacoes"`
}

type ContextoOrientacao struct {
	Orientacoes []OrientacaoEducativa
	Lojas       []Loja
}
