package repo

func Insert_loja(id, nome, categoria string) (string, error) {

	input, err := DB.Prepare("INSERT INTO lojas (id, nome, categoria) VALUES ($1, $2, $3)")
	_, err = input.Exec(id, nome, categoria)

	if err != nil {
		panic(err)
	}
	defer input.Close()

	return "LOJA CRIADA COM SUCESSO!", err

}
