package main

import (
	"bufio"
	utils "codigo/app/repositories"
	"codigo/app/routes"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq" // Driver do Postgres
)

func main() {

	//1º passo, conectar ao servidor
	err := utils.Connect()
	if err != nil {
		log.Fatalln("Erro na conexão com o servidor!")
	}
	defer utils.DB.Close()

	//2º Passo, criar tabelas
	err = utils.Criar_banco()
	if err != nil {
		log.Fatalln("Erro na crição de tabelas SQL")
	}

	// Cria um servidor de arquivos que serve os arquivos da pasta "./static".
	fileserver := http.FileServer(http.Dir("./static"))

	// Associa o servidor de arquivos à rota "/static/".
	http.Handle("/static/", http.StripPrefix("/static/", fileserver))

	// Associa as rotas do sistema
	routes.Rotas()

	// Imprime no console a mensagem indicando que o servidor está rodando na porta 8081.
	fmt.Printf("port running on http://localhost:8081/\n")

	// Inicia o servidor HTTP na porta 8081. Se ocorrer um erro, ele será registrado e o programa será encerrado.
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err) // Registra o erro e encerra o programa.
	}

	var (
		opcao int
		LUC   string
	)

	scanner := bufio.NewScanner(os.Stdin)

	for {
		for {
			fmt.Print("JP MALL\n")
			fmt.Print("MENU\n")
			fmt.Print("1 - CRIAR LOJA\n")
			fmt.Print("2 - EXIBIR TODAS AS LOJAS\n")
			fmt.Print("3 - EDITAR DADOS DAS LOJAS\n")
			fmt.Print("4 - DELETAR LOJA\n")
			fmt.Print("5 - SAIR DO SISTEMA!\n")

			scanner.Scan()
			fmt.Sscan(scanner.Text(), &opcao)
			switch opcao {
			case 1:
				var (
					id, nome, categoria string
				)
				fmt.Print("LUC(ID) DA LOJA: ")
				scanner.Scan()
				id = scanner.Text()
				fmt.Print("NOME DA LOJA: ")
				scanner.Scan()
				nome = scanner.Text()
				fmt.Print("CATEGORIA: ")
				scanner.Scan()
				categoria = scanner.Text()

				result, err := utils.Insert_loja(id, nome, categoria) //A func Exec salva os valores dentro do banco
				if err != nil {
					fmt.Print("ERRO DE INSERT!")
					panic(err)
				}
				fmt.Println(result)

			case 2:
				fmt.Print("\tLUC\t|\tNOME\t|\tCATEGORIA\t\n")
				Lertabela, err := utils.Read_lojas()
				if err != nil {
					panic(err)
				}
				for _, v := range Lertabela {
					fmt.Printf("\t%v\t|\t%v\t|\t%v\t|\n", v.Id, v.Nome, v.Categoria)
				}

			case 3:
				var (
					LUC, dado string
				)

				fmt.Print("Qual a LUC da loja: ")
				scanner.Scan()
				LUC = scanner.Text()
				fmt.Println("Qual dado alterar?")
				fmt.Println("1 - LUC")
				fmt.Println("2 - NOME")
				fmt.Println("3 - CATEGORIA")
				scanner.Scan()
				fmt.Sscan(scanner.Text(), &opcao)

				switch opcao {
				case 1:
					fmt.Print("Qual a nova LUC?")
					scanner.Scan()
					dado = scanner.Text()
					result, _ := (utils.Update_lojas("id", dado, LUC))
					fmt.Println(result)

				case 2:
					fmt.Print("Qual o novo NOME?")
					scanner.Scan()
					dado = scanner.Text()
					result, _ := utils.Update_lojas("nome", dado, LUC)
					fmt.Println(result)
				case 3:
					fmt.Print("Qual a nova CATEGORIA?")
					scanner.Scan()
					dado = scanner.Text()
					result, _ := utils.Update_lojas("categoria", dado, LUC)
					fmt.Println(result)
				}
			case 4:
				fmt.Print("Qual a LUC a ser deletada?")
				scanner.Scan()
				LUC = scanner.Text()
				result, _ := utils.Delete_loja(LUC)
				fmt.Println(result)

			}
			if opcao == 5 {
				break
			}
		}
		if opcao == 5 {
			break
		}
	}
}
