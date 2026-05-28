package repo

import (
	"codigo/app/models"
	"log"
)

type OrientacaoRepository struct{}

func (repo *OrientacaoRepository) Salvar(o models.OrientacaoEducativa) error {
	query := `INSERT INTO orientacoes_educativas (loja_id, responsavel_presente, funcao_responsavel, data_orientacao, observacoes) VALUES ($1, $2, $3, $4, $5)`
	_, err := DB.Exec(query, o.LojaID, o.ResponsavelPresente, o.FuncaoResponsavel, o.DataOrientacao, o.Observacoes)
	if err != nil {
		log.Printf("ERRO INSERT orientacoes_educativas, err: %e\n", err)
		return err
	}
	return err
}

func (repo *OrientacaoRepository) ListarTodas() ([]models.OrientacaoEducativa, error) {
	query := `SELECT o.id, o.loja_id, l.nome, o.responsavel_presente, o.funcao_responsavel, o.data_orientacao, o.observacoes
        FROM orientacoes_educativas o
        INNER JOIN lojas l ON o.loja_id = l.id
        ORDER BY o.data_orientacao DESC`

	rows, err := DB.Query(query) //Ler todos os dados da tabela, e armazena bagunçado em rows

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []models.OrientacaoEducativa //cria uma lista com todas as structs lidas dentro
	for rows.Next() {                      //Itera no banco de dados e armazena em variáveis correspondentes
		var o models.OrientacaoEducativa
		if err := rows.Scan(&o.ID, &o.LojaID, &o.NomeLoja, &o.ResponsavelPresente, &o.FuncaoResponsavel, &o.DataOrientacao, &o.Observacoes); err != nil {
			return nil, err
		}
		lista = append(lista, o) //ao final da leitura de cada linha, adiciona um "pacote" inteiro na lista
	}
	return lista, nil //Retorna a lista com os dados para leitura e renderização no front
}

func (repo *OrientacaoRepository) BuscaPorID(id int) (models.OrientacaoEducativa, error) {
	var o models.OrientacaoEducativa
	query := `SELECT id, loja_id, responsavel_presente, funcao_resposnavel, data_orientacao, observacoes, FROM
	orientacoes_educativas WHERE id = $1`
	err := DB.QueryRow(query, id).Scan(&o.ID, &o.LojaID, &o.ResponsavelPresente, &o.FuncaoResponsavel, &o.DataOrientacao, &o.Observacoes)
	if err != nil {
		return o, err
	}
	return o, nil
}

func (repo *OrientacaoRepository) Atualizar(o models.OrientacaoEducativa) error {

	query := `
		UPDATE orientacoes_educativas 
		SET responsavel_presente = $1, 
		    funcao_responsavel = $2, 
		    data_orientacao = $3, 
		    observacoes = $4 
		WHERE id = $5
	`

	// Executa a query passando os valores na ordem dos $1, $2, etc.
	_, err := DB.Exec(query,
		o.ResponsavelPresente,
		o.FuncaoResponsavel,
		o.DataOrientacao,
		o.Observacoes,
		o.ID,
	)

	if err != nil {
		return err // Se o banco der erro (ex: tipo de dado inválido), joga para cima
	}

	return nil
}

func (repo *OrientacaoRepository) Delete(o models.OrientacaoEducativa) error {
	query := `DELETE FROM orientacoes_educativas WHERE id =$1`
	_, err := DB.Exec(query, o.ID)
	if err != nil {
		return err // Se o banco der erro (ex: tipo de dado inválido), joga para cima
	}
	return nil
}
