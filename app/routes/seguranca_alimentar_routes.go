package routes

import (
	"codigo/app/controllers"
	"net/http"
)

func RotasSegurancaAlimentar() {
	var segurancaController controllers.SegurancaAlimentarController

	http.HandleFunc("/conservacao/seguranca-alimentar", segurancaController.ListarPaginaHandler)

	http.HandleFunc("/api/lojas", segurancaController.ListarLojasHandler)
	http.HandleFunc("/api/inspecoes", segurancaController.ListarHandler)
	http.HandleFunc("/api/inspecoes/criar", segurancaController.SalvarHandler)

	http.HandleFunc("/api/inspecoes/editar", segurancaController.EditarHandler)
	http.HandleFunc("/api/inspecoes/deletar", segurancaController.ExcluirHandler)
	http.HandleFunc("/api/inspecoes/pdf", segurancaController.AbrirPDFHandler)
}
