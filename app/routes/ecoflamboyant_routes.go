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
}
