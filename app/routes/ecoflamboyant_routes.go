package routes

import (
	"codigo/app/controllers"
	"net/http"
)

func RotasEcoflamboyant() {
	var ctrl controllers.EcoflamboyantController

	http.HandleFunc("/conservacao/eco-flamboyant", ctrl.ListarEcoFlamboyantHandler)
	http.HandleFunc("/conservacao/eco-flamboyant/criar", ctrl.CriarParticipanteHandler)
	http.HandleFunc("/conservacao/eco-flamboyant/status", ctrl.AlterarStatusLoja)
	http.HandleFunc("/conservacao/eco-flamboyant/termo/", ctrl.DownloadTermo)
	http.HandleFunc("/api/conservacao/eco-flamboyant/buscar-lojas", ctrl.BuscarLojasDisponiveis)
	http.HandleFunc("/conservacao/eco-flamboyant/relatorio/pdf", ctrl.EmitirRelatorioPDF)
	http.HandleFunc("/conservacao/eco-flamboyant/loja/editar", ctrl.EditarParticipanteHandler)
	http.HandleFunc("/conservacao/eco-flamboyant/loja/remover", ctrl.RemoverParticipanteHandler)
	http.HandleFunc("/conservacao/eco-flamboyant/residuo/editar", ctrl.EditarResiduoHandler)
	http.HandleFunc("/conservacao/eco-flamboyant/residuo/excluir", ctrl.ExcluirResiduoHandler)
}
