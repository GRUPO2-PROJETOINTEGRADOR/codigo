package routes

import (
	"codigo/app/controllers"
	"net/http"
)

func RotasEcoflamboyant() {
	var ctrl controllers.EcoflamboyantController

	http.HandleFunc("/conservacao/eco-flamboyant", ctrl.ListarEcoFlamboyantHandler)
	http.HandleFunc("/conservacao/eco-flamboyant/criar", ctrl.CriarParticipanteHandler)
}
