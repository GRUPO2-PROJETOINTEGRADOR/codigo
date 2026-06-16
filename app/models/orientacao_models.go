package models

import (
	"time"
)

type OrientacaoEducativa struct {
	ID                  int        `json:"id"`
	LojaID              string     `json:"loja_id"`
	NomeLoja            string     `json:"nome_loja"`
	ResponsavelPresente string     `json:"responsavel_presente"`
	FuncaoResponsavel   string     `json:"funcao_responsavel"`
	DataOrientacao      time.Time  `json:"data_orientacao"`
	Observacoes         string     `json:"observacoes"`
	Signatario          string     `json:"signatario"`
	DataAssinatura      *time.Time `json:"data_assinatura"`
}

type ContextoOrientacao struct {
	Orientacoes []OrientacaoEducativa
	Lojas       []Loja
}
