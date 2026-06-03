package models

import "time"

type Participante struct {
	LojaID      string     `json:"loja_id"`
	LojaName    string     `json:"loja_nome"`
	Status      bool       `json:"status_participacao"`
	DataEntrada time.Time  `json:"data_entrada"`
	DataSaida   *time.Time `json:"data_saida,omitempty"`
	AnexoEco    string     `json:"anexo_eco"`
}

type EcoFlamboyantPageData struct {
	Participantes []Participante
	Lojas         []Loja
}
