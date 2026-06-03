package routes

import (
	"codigo/app/controllers"
	"net/http"
)

func Rotas() {
	RotasOrientacao()
	RotasSegurancaAlimentar()
	RotasEcoflamboyant()

	var dashCtrl controllers.DashboardController
	http.HandleFunc("/conservacao/dashboard", dashCtrl.ListarPaginaHandler)

	var relCtrl controllers.RelatoriosController
	http.HandleFunc("/conservacao/relatorios", relCtrl.ListarPaginaHandler)

	http.HandleFunc("/conservacao/relatorio-seguranca-alimentar", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/conservacao/seguranca-alimentar", http.StatusMovedPermanently)
	})

	http.HandleFunc("/conservacao/relatorio-seguranca-alimentar.html", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/conservacao/seguranca-alimentar", http.StatusMovedPermanently)
	})

	http.HandleFunc("/conservacao/seguranca-alimentar.html", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/conservacao/seguranca-alimentar", http.StatusMovedPermanently)
	})

	http.HandleFunc("/conservacao/relatorios.html", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/conservacao/relatorios", http.StatusMovedPermanently)
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "templates/index.html")
	})
}
