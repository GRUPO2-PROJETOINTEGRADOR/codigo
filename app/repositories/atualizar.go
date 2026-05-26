package utils

import "fmt"

func Update_lojas(coluna, novo, LUC string) (string, error) {

	query := fmt.Sprintf("UPDATE lojas SET %s = $1 WHERE id = $2", coluna)
	update, err := DB.Prepare(query)
	update.Exec(novo, LUC)
	defer update.Close()
	return "Sistema atualizado com sucesso!", err
}
