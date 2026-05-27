package controllers

import (
    "encoding/json"
    "codigo/app/models"
    utils "codigo/app/repository"
    s "codigo/app/services"
    "html/template"
    "net/http"
    "time"
)

// existing code ... (keep unchanged)

func (c *OrientacaoController) ListarJSONHandler(w http.ResponseWriter, r *http.Request) {
    orientacoes, err := c.Service.ListarTodas()
    if err != nil {
        http.Error(w, "Erro ao buscar orientações", http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(orientacoes)
}



type OrientacaoController struct {
    Service s.OrientacaoService
}

func (c *OrientacaoController) ListarPaginaHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Busca as orientações do banco
	orientacoes, err := c.Service.ListarTodas()
	if err != nil {
		http.Error(w, "Erro ao buscar orientações", http.StatusInternalServerError)
		return
	}

	// 2. Busca as lojas do banco
	lojas, err := utils.Read_lojas()
	if err != nil {
		http.Error(w, "Erro ao buscar lojas", http.StatusInternalServerError)
		return
	}

	// 3. Monta o contexto com os dois dados
	contexto := models.ContextoOrientacao{
		Orientacoes: orientacoes,
		Lojas:       lojas,
	}

	// 4. RENDERIZA APENAS A TELA DE ORIENTAÇÃO
	// Como a tela já está completa, passamos apenas o caminho dela no ParseFiles
	tmpl := template.Must(template.ParseFiles(
		"templates/conservacao/orientacao-educativa.html",
	))

	// Mudamos para tmpl.Execute (sem a palavra Template), pois não precisamos
	// chamar um bloco específico como "layout", ele vai rodar o arquivo inteiro direto.
	tmpl.Execute(w, contexto)
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

	err := c.Service.CriarNovaOrientacao(novaOrientacao)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Atualiza a página limpando o formulário
	http.Redirect(w, r, "/conservacao/orientacao-educativa", http.StatusSeeOther)
}
