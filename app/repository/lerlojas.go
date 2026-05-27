package utils

import (
	"codigo/app/models"
	"log"
)

func Read_lojas() ([]models.Loja, error) {

	dados, err := DB.Query("SELECT id, nome, categoria FROM lojas")
	if err != nil {
		log.Printf("Erro em READ_LOJAS, err: %e", err)
	}
	defer dados.Close()

	var lojas []models.Loja

	for dados.Next() {
		var tb models.Loja
		dados.Scan(
			&tb.ID,
			&tb.Nome,
			&tb.Categoria,
		)

		lojas = append(lojas, tb)

	}

	return lojas, nil
}
