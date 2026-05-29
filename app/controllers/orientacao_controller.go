package controllers

import (
	"codigo/app/models"
	utils "codigo/app/repository"
	s "codigo/app/services"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
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

	//log.Printf("===> QUANTIDADE LIDA DO BANCO: %d", len(orientacoes)) LOG PARA TESTE
	//log.Printf("===> DADOS: %+v", orientacoes)

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
	w.WriteHeader(http.StatusOK)
}

func (c *OrientacaoController) Editar(w http.ResponseWriter, r *http.Request) {
	// 1. Segurança: Só aceita se o HTML enviar como POST (envio de formulário)
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido. Use POST.", http.StatusMethodNotAllowed)
		return
	}

	// 2. Manda o Go ler o pacote de dados vindo do HTML
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Erro ao processar os dados do formulário", http.StatusBadRequest)
		return
	}

	idStr := r.FormValue("id")                         // Pega do <input type="hidden" name="id">
	responsavel := r.FormValue("responsavel_presente") // Pega do <input name="responsavel_presente">
	funcao := r.FormValue("funcao_responsavel")        // Pega do <input name="funcao_responsavel">
	dataStr := r.FormValue("data_orientacao")          // Pega do <input type="date" name="data_orientacao">
	observacoes := r.FormValue("observacoes")          // Pega do <textarea name="observacoes">

	// 4. Conversão de Tipos: O HTML envia tudo como texto (string). O Go precisa converter.
	// Convertendo o ID de String para Inteiro (int)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "O ID enviado não é um número válido", http.StatusBadRequest)
		return
	}

	// Convertendo a Data do formato HTML (2026-05-28) para o formato de tempo do Go (time.Time)
	dataOrientacao, err := time.Parse("2006-01-02", dataStr)
	if err != nil {
		http.Error(w, "Formato de data inválido. Use o padrão do calendário.", http.StatusBadRequest)
		return
	}

	// 5. Montando a nossa Struct com os dados limpos e convertidos
	orientacaoAtualizada := models.OrientacaoEducativa{
		ID:                  id,
		ResponsavelPresente: responsavel,
		FuncaoResponsavel:   funcao,
		DataOrientacao:      dataOrientacao,
		Observacoes:         observacoes,
		// LojaID não entra aqui porque decidimos travar a edição da loja!
	}

	// 6. Envia para o Service aplicar as regras e salvar no banco
	err = c.Service.Atualizar(orientacaoAtualizada)
	if err != nil {
		// Se der erro em alguma validação ou no SQL, avisa a tela
		http.Error(w, "Erro ao salvar alterações: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 7. SUCESSO COMPLETO!
	http.Redirect(w, r, "/conservacao/orientacao-educativa", http.StatusSeeOther)
}

func (c *OrientacaoController) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método não permitido. Use DELETE.", http.StatusMethodNotAllowed)
		return
	}

	IDStr := r.FormValue("id")
	id, err := strconv.Atoi(IDStr)
	if err != nil {
		log.Printf("Erro ao converter ID '%s' para inteiro: %v", IDStr, err)
		http.Error(w, "ID de registro inválido", http.StatusBadRequest)
		return
	}
	IdOrientacao := models.OrientacaoEducativa{
		ID: id,
	}

	err = c.Service.Delete(IdOrientacao)
	if err != nil {
		log.Printf("Erro ao deletar registro: %v", err)
		http.Error(w, "Erro ao deletar registro: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
