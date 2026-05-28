package repo

import "fmt"

func Delete_loja(id string) (string, error) {
	q_delete, err := DB.Prepare("DELETE FROM lojas WHERE id = $1")
	if err != nil {
		return "", err
	}
	defer q_delete.Close()
	q_delete.Exec(id)

	return fmt.Sprintf("UNIDADE %s DELETADA!", id), err
}
