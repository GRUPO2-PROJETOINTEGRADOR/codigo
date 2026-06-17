package services

import (
	"codigo/app/models"
	utils "codigo/app/repository"
	"database/sql"
	"errors"
	"strconv"
	"time"
)

func ListarLojas(db *sql.DB) ([]models.Loja, error) {
	return utils.ListarLojas(db)
}

func ObterLojasParticipantes(db *sql.DB) ([]models.Loja, error) {
	return utils.ListarLojasParticipantes(db)
}

func CriarParticipante(db *sql.DB, lojaID string, dataEntrada time.Time, dataSaida *time.Time, nomeAnexo string, dadosAnexo []byte) error {
	if lojaID == "" {
		return errors.New("loja obrigatória")
	}
	if dataEntrada.IsZero() {
		return errors.New("data de entrada obrigatória")
	}
	if nomeAnexo == "" {
		return errors.New("termo de aceite obrigatório")
	}
	if dataSaida != nil && !dataSaida.IsZero() && dataSaida.Before(dataEntrada) {
		return errors.New("A data de saída não pode ser anterior à data de entrada.")
	}

	if err := utils.CriarParticipante(db, lojaID, dataEntrada, dataSaida, nomeAnexo, dadosAnexo); err != nil {
		return err
	}
	return utils.InserirAuditoria(db, lojaID, "eco_participante", "cadastro")
}

func ListarParticipantes(db *sql.DB, dataInicio, dataFim string, limit, offset int) ([]models.Participante, error) {
	return utils.ListarParticipantes(db, dataInicio, dataFim, limit, offset)
}

func ContarParticipantes(db *sql.DB, dataInicio, dataFim string) (int, error) {
	return utils.ContarParticipantes(db, dataInicio, dataFim)
}

func CriarResiduo(db *sql.DB, lojaID, dataColetaStr, pesoKGStr, aproveitadoStr string) error {
	if lojaID == "" {
		return errors.New("loja obrigatória")
	}
	if dataColetaStr == "" {
		return errors.New("data de coleta obrigatória")
	}
	if pesoKGStr == "" {
		return errors.New("peso obrigatório")
	}
	if aproveitadoStr == "" {
		return errors.New("informe se foi aproveitado para adubo")
	}

	dataColeta, err := time.Parse("2006-01-02", dataColetaStr)
	if err != nil {
		return errors.New("data de coleta inválida")
	}
	if dataColeta.After(time.Now()) {
		return errors.New("data de coleta não pode ser futura")
	}

	pesoKG, err := strconv.ParseFloat(pesoKGStr, 64)
	if err != nil {
		return errors.New("peso inválido")
	}
	if pesoKG <= 0 {
		return errors.New("peso deve ser maior que zero")
	}
	if pesoKG > 9999.99 {
		return errors.New("O peso não pode ultrapassar 9.999,99 kg por registro")
	}

	aproveitado := aproveitadoStr == "Sim"
	if err := utils.InserirResiduo(db, lojaID, dataColeta, pesoKG, aproveitado); err != nil {
		return err
	}
	return utils.InserirAuditoria(db, lojaID, "residuo", "cadastro")
}

func ObterResiduos(db *sql.DB, dataInicio, dataFim, lojaID string, limit, offset int) ([]models.Residuo, error) {
	return utils.ListarResiduos(db, dataInicio, dataFim, lojaID, limit, offset)
}

func ContarResiduos(db *sql.DB, dataInicio, dataFim, lojaID string) (int, error) {
	return utils.ContarResiduos(db, dataInicio, dataFim, lojaID)
}

func CriarKit(db *sql.DB, lojaID, dataEntregaKitStr, qntKitStr string) error {
	if lojaID == "" {
		return errors.New("loja obrigatória")
	}
	if dataEntregaKitStr == "" {
		return errors.New("data de entrega obrigatória")
	}
	if qntKitStr == "" {
		return errors.New("quantidade obrigatória")
	}

	dataEntregaKit, err := time.Parse("2006-01-02", dataEntregaKitStr)
	if err != nil {
		return errors.New("data de entrega inválida")
	}
	if dataEntregaKit.After(time.Now()) {
		return errors.New("data de entrega não pode ser futura")
	}

	qntKit, err := strconv.Atoi(qntKitStr)
	if err != nil {
		return errors.New("quantidade inválida")
	}
	if qntKit <= 0 {
		return errors.New("quantidade deve ser maior que zero")
	}

	if err := utils.InserirKit(db, lojaID, dataEntregaKit, qntKit); err != nil {
		return err
	}
	return utils.InserirAuditoria(db, lojaID, "kit", "cadastro")
}

func ObterKits(db *sql.DB, dataInicio, dataFim, lojaID string, limit, offset int) ([]models.Kit, error) {
	return utils.ListarKits(db, dataInicio, dataFim, lojaID, limit, offset)
}

func ContarKits(db *sql.DB, dataInicio, dataFim, lojaID string) (int, error) {
	return utils.ContarKits(db, dataInicio, dataFim, lojaID)
}

func ObterFluxoKits(db *sql.DB) ([]models.PontoKits, error) {
	return utils.FluxoKitsPorPeriodo(db)
}

func ObterTotalKits(db *sql.DB) (int, error) {
	return utils.SomarTotalKits(db)
}

func ObterDadosLojas(db *sql.DB) (int, []models.PontoLojas, error) {
	total, err := utils.ContarLojasAtivas(db)
	if err != nil {
		return 0, nil, err
	}
	crescimento, err := utils.CrescimentoLojasPorMes(db)
	if err != nil {
		return 0, nil, err
	}
	return total, crescimento, nil
}

func InativarLoja(db *sql.DB, lojaID string) error {
	if err := utils.InativarLoja(db, lojaID); err != nil {
		return err
	}
	return utils.InserirAuditoria(db, lojaID, "eco_participante", "inativacao")
}

func AtivarLoja(db *sql.DB, lojaID string) error {
	if err := utils.AtivarLoja(db, lojaID); err != nil {
		return err
	}
	return utils.InserirAuditoria(db, lojaID, "eco_participante", "reativacao")
}

func ListarAuditorias(db *sql.DB, dataInicio, dataFim, lojaID string, limit, offset int) ([]models.RegistroAuditoria, error) {
	return utils.ListarAuditoriasEventos(db, dataInicio, dataFim, lojaID, limit, offset)
}

func ContarAuditorias(db *sql.DB, dataInicio, dataFim, lojaID string) (int, error) {
	return utils.ContarAuditoriasEventos(db, dataInicio, dataFim, lojaID)
}

func ObterResumoResiduos(db *sql.DB) (totalGeral, totalAdubo, totalDescarte, taxa float64, fluxo []models.PontoResiduos, err error) {
	totalGeral, totalAdubo, totalDescarte, err = utils.ResumoResiduos(db)
	if err != nil {
		return
	}
	fluxo, err = utils.FluxoResiduosPorPeriodo(db)
	if err != nil {
		return
	}
	if totalGeral > 0 {
		taxa = (totalAdubo / totalGeral) * 100
	}
	return
}

func BuscarLojasDisponiveis(db *sql.DB, q string) ([]models.Loja, error) {
	return utils.BuscarLojasDisponiveis(db, q)
}
