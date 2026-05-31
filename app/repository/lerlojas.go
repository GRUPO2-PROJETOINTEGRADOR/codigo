package repo

import (
	"codigo/app/models"
	"log"
)

func Read_lojas() ([]models.Loja, error) {

	dados, err := DB.Query("SELECT id, nome, categoria FROM lojas")
	if err != nil {
		log.Printf("Erro em READ_LOJAS, err: %e", err)
		return nil, err
	}
	defer dados.Close()

	var lojas []models.Loja

	for dados.Next() {
		var tb models.Loja
		err := dados.Scan(
			&tb.ID,
			&tb.Nome,
			&tb.Categoria,
		)
		if err != nil {
			log.Println("Erro repo Ler lojas: ", err)
			return nil, err
		}

		lojas = append(lojas, tb)

	}

	return lojas, err
}
