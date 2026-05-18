package funcoes

import (
	"database/sql"
	"fmt"
)

const server = "host=localhost port=5432 user=postgres password=1234 dbname=conservacao sslmode=disable" //Dados do servidor
func Insert_loja(id, nome, categoria string) (string, error) {

	database, err := sql.Open("postgres", server)
	input, err := database.Prepare("INSERT INTO lojas (id, nome, categoria) VALUES ($1, $2, $3)")
	_, err = input.Exec(id, nome, categoria)

	if err != nil {
		panic(err)
	}
	database.Close()
	return "LOJA CRIADA COM SUCESSO!", err

}

type Lojas struct {
	Id        string
	Nome      string
	Categoria string
}

func Read_lojas() ([]Lojas, error) {

	database, err := sql.Open("postgres", server)

	if err != nil {
		panic(err)
	}

	dados, err := database.Query("SELECT id, nome, categoria FROM lojas")

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

	database.Close()
	return lojas, nil
}

func Update_lojas(coluna, novo, LUC string) (string, error) {
	database, err := sql.Open("postgres", server)
	if err != nil {
		panic(err)
	}
	query := fmt.Sprintf("UPDATE lojas SET %s = $1 WHERE id = $2", coluna)
	update, err := database.Prepare(query)
	update.Exec(novo, LUC)
	defer database.Close()

	return "Sistema atualizado com sucesso!", err
}

func Delete_loja(id string) (string, error) {
	db, err := sql.Open("postgres", server)
	if err != nil {
		panic(err)
	}
	q_delete, err := db.Prepare("DELETE FROM lojas WHERE id = $1")
	q_delete.Exec(id)
	db.Close()
	return fmt.Sprintf("UNIDADE %s DELETADA!", id), err
}
