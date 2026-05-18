package main

import (
	"backend/funcoes"
	"bufio"
	"fmt"
	"os"

	_ "github.com/lib/pq" // Driver do Postgres
)

func main() {

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

				result, err := funcoes.Insert_loja(id, nome, categoria) //A func Exec salva os valores dentro do banco
				if err != nil {
					fmt.Print("ERRO DE INSERT!")
					panic(err)
				}
				fmt.Println(result)

			case 2:
				fmt.Print("\tLUC\t|\tNOME\t|\tCATEGORIA\t\n")
				Lertabela, err := funcoes.Read_lojas()
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
					result, _ := (funcoes.Update_lojas("id", dado, LUC))
					fmt.Println(result)

				case 2:
					fmt.Print("Qual o novo NOME?")
					scanner.Scan()
					dado = scanner.Text()
					result, _ := funcoes.Update_lojas("nome", dado, LUC)
					fmt.Println(result)
				case 3:
					fmt.Print("Qual a nova CATEGORIA?")
					scanner.Scan()
					dado = scanner.Text()
					result, _ := funcoes.Update_lojas("categoria", dado, LUC)
					fmt.Println(result)
				}
			case 4:
				fmt.Print("Qual a LUC a ser deletada?")
				scanner.Scan()
				LUC = scanner.Text()
				result, _ := funcoes.Delete_loja(LUC)
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
