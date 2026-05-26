package utils

import "log"

type Lojas struct {
	Id        string
	Nome      string
	Categoria string
}

func Read_lojas() ([]Lojas, error) {

	dados, err := DB.Query("SELECT id, nome, categoria FROM lojas")
	if err != nil {
		log.Printf("Erro em READ_LOJAS, err: %e", err)
	}
	defer dados.Close()

	var lojas []Lojas

	for dados.Next() {
		var tb Lojas
		dados.Scan(
			&tb.Id,
			&tb.Nome,
			&tb.Categoria,
		)

		lojas = append(lojas, tb)

	}

	return lojas, nil
}
