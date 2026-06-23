package routes

import (
	"codigo/app/controllers"
	"net/http"
)

func Rotas() {
	RotasOrientacao()
	RotasSegurancaAlimentar()
	RotasEcoflamboyant()

	var ecoCtrl controllers.EcoflamboyantController
	http.HandleFunc("/conservacao/eco-flamboyant/residuos", ecoCtrl.CriarResiduoHandler)
	http.HandleFunc("/conservacao/eco-flamboyant/kits", ecoCtrl.CriarKitHandler)

	var dashCtrl controllers.DashboardController
	http.HandleFunc("/conservacao/dashboard", dashCtrl.ListarPaginaHandler)

	http.HandleFunc("/conservacao/seguranca-alimentar.html", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/conservacao/seguranca-alimentar", http.StatusMovedPermanently)
	})

	http.HandleFunc("/conservacao/dashboard.html", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/conservacao/dashboard", http.StatusMovedPermanently)
	})

	http.HandleFunc("/conservacao/eco-flamboyant.html", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/conservacao/eco-flamboyant", http.StatusMovedPermanently)
	})

	http.HandleFunc("/conservacao/orientacao-educativa.html", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/conservacao/orientacao-educativa", http.StatusMovedPermanently)
	})

}
