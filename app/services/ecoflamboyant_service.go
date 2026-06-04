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

	pesoKG, err := strconv.ParseFloat(pesoKGStr, 64)
	if err != nil {
		return errors.New("peso inválido")
	}
	if pesoKG <= 0 {
		return errors.New("peso deve ser maior que zero")
	}

	aproveitado := aproveitadoStr == "Sim"
	return utils.InserirResiduo(db, lojaID, dataColeta, pesoKG, aproveitado)
}

func ObterResiduos(db *sql.DB) ([]models.Residuo, error) {
	return utils.ListarResiduos(db)
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

	qntKit, err := strconv.Atoi(qntKitStr)
	if err != nil {
		return errors.New("quantidade inválida")
	}
	if qntKit <= 0 {
		return errors.New("quantidade deve ser maior que zero")
	}

	return utils.InserirKit(db, lojaID, dataEntregaKit, qntKit)
}

func ObterKits(db *sql.DB) ([]models.Kit, error) {
	return utils.ListarKits(db)
}
