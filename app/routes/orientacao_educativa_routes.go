package routes

import (
	"codigo/app/controllers"
	repo "codigo/app/repository"
	"codigo/app/services"
	"net/http"
)

func RotasOrientacao() {
	orientacaoController := controllers.OrientacaoController{
		Service: services.OrientacaoService{
			Repo: &repo.OrientacaoRepository{},
		},
	}

	http.HandleFunc("/conservacao/orientacao-educativa", orientacaoController.ListarPaginaHandler)
	http.HandleFunc("/conservacao/orientacoes/stats", orientacaoController.ExibirStats)
	http.HandleFunc("/api/orientacoes", orientacaoController.ListarJSONHandler)
	http.HandleFunc("/conservacao/orientacao-educativa/salvar", orientacaoController.SalvarHandler)
	http.HandleFunc("/conservacao/orientacao-educativa/atualizar", orientacaoController.EditarHandler)
	http.HandleFunc("/conservacao/orientacao-educativa/delete", orientacaoController.DeleteHandler)

	http.HandleFunc("/teste", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("rotas funcionando"))
	})
}
