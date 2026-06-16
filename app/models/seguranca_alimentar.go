package models

type SegurancaAlimentar struct {
	ID               int    `json:"id"`
	LojaID           string `json:"loja_id"`
	DataAuditoria    string `json:"data_auditoria"`
	ResponsavelLoja  string `json:"responsavel_loja"`
	CargoResponsavel string `json:"cargo_responsavel"`
	Nota             int    `json:"nota"`
	AnexoTiller      string `json:"anexo_tiller"`
	Classificacao    string `json:"classificacao"`
	TipoInspecao     string `json:"tipo_inspecao"`
	NCGrave          bool   `json:"nc_grave"`
	PDFNome          string `json:"pdf_nome"`
	PDFTipo          string `json:"pdf_tipo"`
	PDFTamanho       int64  `json:"pdf_tamanho"`
	PDFArquivo       []byte `json:"-"`
}
