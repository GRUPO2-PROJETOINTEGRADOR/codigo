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

type Residuo struct {
	ID          int
	LojaNome    string
	DataColeta  time.Time
	PesoKG      float64
	Aproveitado bool
}

type Kit struct {
	ID             int
	LojaNome       string
	DataEntregaKit time.Time
	QntKit         int
}

type PontoLojas struct {
	Mes      string
	Total    int
	Entradas int
}

type PontoKits struct {
	Periodo  string
	Total    int
	Entregas int
}

type EcoFlamboyantPageData struct {
	Participantes          []Participante
	Lojas                  []Loja
	Residuos               []Residuo
	TotalResiduos          int
	Kits                   []Kit
	TotalKits              int
	TotalLojasParticipantes int
	CrescimentoLojas       []PontoLojas
	FluxoKits              []PontoKits
}
