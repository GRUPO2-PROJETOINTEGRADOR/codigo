package routes

import "net/http"

func Rotas() {
	http.HandleFunc("/teste", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("rotas funcionando"))
	})

	RotasOrientacao()
	RotasSegurancaAlimentar()
}
