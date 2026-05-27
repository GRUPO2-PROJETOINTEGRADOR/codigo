package routes

import (
	"codigo/app/controllers"
	"net/http"
)

func RotasSegurancaAlimentar() {
	var segurancaController controllers.SegurancaAlimentarController

	http.HandleFunc("/conservacao/seguranca-alimentar", segurancaController.ListarPaginaHandler)
	http.HandleFunc("/conservacao/seguranca-alimentar/salvar", segurancaController.SalvarHandler)
	http.HandleFunc("/conservacao/seguranca-alimentar/editar", segurancaController.EditarHandler)
	http.HandleFunc("/conservacao/seguranca-alimentar/excluir", segurancaController.ExcluirHandler)
}
