package utils

import (
	"codigo/app/models"
	"log"
)

type OrientacaoRepository struct{}

func (repo *OrientacaoRepository) Salvar(o models.OrientacaoEducativa) error {
	query := `INSERT INTO orientacoes_educativas (loja_id, responsavel_presente, funcao_responsavel, data_orientacao, observacoes)`
	_, err := DB.Exec(query, o.LojaID, o.ResponsavelPresente, o.FuncaoResponsavel, o.DataOrientacao, o.Observacoes)
	if err != nil {
		log.Printf("ERRO INSERT orientacoes_educativas, err: %e\n", err)
		return err
	}
	return err
}

func (repo *OrientacaoRepository) ListarTodas() ([]models.OrientacaoEducativa, error) {
	query := `SELECT id, loja_id, responsavel_presente, funcao_responsavel, data_orientacao, observacoes 
	FROM orientacoes_educativas ORDER BY data_orientacao DESC`

	rows, err := DB.Query(query) //Ler todos os dados da tabela, e armazena bagunçado em rows

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []models.OrientacaoEducativa //cria uma lista com todas as structs lidas dentro
	for rows.Next() {                      //Itera no banco de dados e armazena em variáveis correspondentes
		var o models.OrientacaoEducativa
		if err := rows.Scan(&o.ID, &o.LojaID, &o.ResponsavelPresente, &o.FuncaoResponsavel, &o.DataOrientacao, &o.Observacoes); err != nil {
			return nil, err
		}
		lista = append(lista, o) //ao final da leitura de cada linha, adiciona um "pacote" inteiro na lista
	}
	return lista, nil //Retorna a lista com os dados para leitura e renderização no front
}
