package services

import (
	"codigo/app/models"
	utils "codigo/app/repository"
	"database/sql"
	"errors"
	"time"
)

func ListarLojas(db *sql.DB) ([]models.Loja, error) {
	return utils.ListarLojas(db)
}

func CriarParticipante(db *sql.DB, lojaID string, dataEntrada time.Time, dataSaida *time.Time, caminhoAnexo string) error {
	if lojaID == "" {
		return errors.New("loja obrigatória")
	}
	if dataEntrada.IsZero() {
		return errors.New("data de entrada obrigatória")
	}
	if caminhoAnexo == "" {
		return errors.New("termo de aceite obrigatório")
	}
	if dataSaida != nil && !dataSaida.After(dataEntrada) {
		return errors.New("data de saída deve ser após a data de entrada")
	}

	p := models.Participante{
		LojaID:      lojaID,
		DataEntrada: dataEntrada,
		DataSaida:   dataSaida,
		AnexoEco:    caminhoAnexo,
	}

	return utils.CriarParticipante(db, p)
}

func ListarParticipantes(db *sql.DB) ([]models.Participante, error) {
	return utils.ListarParticipantes(db)
}
