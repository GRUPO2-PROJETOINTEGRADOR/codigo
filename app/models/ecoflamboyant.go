package models

import "time"

type RegistroAuditoria struct {
	ID         int
	LojaNome   string
	Entidade   string
	Acao       string
	DataEvento time.Time
}

type Participante struct {
	LojaID      string     `json:"loja_id"`
	LojaName    string     `json:"loja_nome"`
	Status      bool       `json:"status_participacao"`
	DataEntrada time.Time  `json:"data_entrada"`
	DataSaida   *time.Time `json:"data_saida,omitempty"`
	AnexoEcoNome string    `json:"anexo_eco_nome"`
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

type PontoResiduos struct {
	Periodo      string
	PesoAdubo    float64
	PesoDescarte float64
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
	VolumeTotalGeral       float64
	TotalAdubo             float64
	TotalDescartado        float64
	TaxaAproveitamento     float64
	FluxoResiduos          []PontoResiduos
	Registros              []RegistroAuditoria
	AbaAtiva               string
	TodasLojas             []Loja
	FiltroDataInicio       string
	FiltroDataFim          string
	FiltroLojaID           string
	NomeLojaFiltrada       string
	HojeStr                string
	ErroForm               string
	PaginaAtual             int
	TotalPaginasLojas       int
	TotalPaginasResiduos    int
	TotalPaginasKits        int
	TotalPaginasRegistros   int
	ItensPorPagina          int
}
