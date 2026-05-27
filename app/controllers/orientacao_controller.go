package controllers

import (
	"codigo/app/models"
	s "codigo/app/services"
	"html/template"
	"net/http"
	"time"
)

type OrientacaoController struct {
	service s.OrientacaoService
}

func (c *OrientacaoController) ListarPaginaHandler(w http.ResponseWriter, r *http.Request) {
	orientacoes, err := c.service.ListarTodas()
	if err != nil {
		http.Error(w, "Erro ao bucar orientações", http.StatusInternalServerError)
		return
	}

	contexto := models.ContextoOrientacao{
		Orientacoes: orientacoes,
	}

	tmpl := template.Must(template.ParseFiles(
		"templates/layout.html",
		"templates/conservacao/orientacao-educativa.html",
	))
	tmpl.ExecuteTemplate(w, "layout", contexto)
}

func (c *OrientacaoController) SalvarHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	dataStr := r.FormValue("data_orientacao")
	dataformat, _ := time.Parse("2006-01-02", dataStr)

	novaOrientacao := models.OrientacaoEducativa{
		LojaID:              r.FormValue("loja_id"), //lê do html a var name="loja_id"
		ResponsavelPresente: r.FormValue("responsavel_presente"),
		FuncaoResponsavel:   r.FormValue("funcao_responsavel"),
		DataOrientacao:      dataformat,
		Observacoes:         r.FormValue("observacoes"),
	}
	
	err := c.service.CriarNovaOrientacao(novaOrientacao)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Atualiza a página limpando o formulário
	http.Redirect(w, r, "/conservacao/orientacao-educativa", http.StatusSeeOther)
}
