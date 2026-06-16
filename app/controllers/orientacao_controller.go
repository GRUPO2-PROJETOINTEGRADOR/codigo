package controllers

import (
	"codigo/app/models"
	repo "codigo/app/repository"
	s "codigo/app/services"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

type OrientacaoController struct {
	Service s.OrientacaoService
}

func (c *OrientacaoController) ListarPaginaHandler(w http.ResponseWriter, r *http.Request) { //Lê os dados da tabela e exibe
	// 1. Busca as orientações do banco
	orientacoes, err := c.Service.ListarTodas()
	if err != nil {
		http.Error(w, "Erro ao buscar orientações", http.StatusInternalServerError)
		return
	}

	//log.Printf("===> QUANTIDADE LIDA DO BANCO: %d", len(orientacoes)) LOG PARA TESTE
	//log.Printf("===> DADOS: %+v", orientacoes)

	// 2. Busca as lojas do banco
	lojas, err := repo.Read_lojas()
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
	tmpl := template.Must(template.New("orientacao-educativa.html").Funcs(template.FuncMap{
		"json": func(v interface{}) template.JS {
			b, _ := json.Marshal(v)
			return template.JS(b)
		},
	}).ParseFiles("templates/conservacao/orientacao-educativa.html"))
	tmpl.Execute(w, contexto)

}

func (c *OrientacaoController) ExibirStats(w http.ResponseWriter, r *http.Request) { //
	Total, err := c.Service.TotalTreinos()
	if err != nil {
		http.Error(w, "Erro ao Exibir Total Treinos", http.StatusInternalServerError)
		log.Printf("Erro na exibição TotalTreinos: %s", err)
		return
	}

	TotalLojas, err := c.Service.LojasTreinos()
	if err != nil {
		http.Error(w, "Erro ao Exibir Total Lojas", http.StatusInternalServerError)
		log.Printf("Erro na exibição TotalLojas: %s", err)
		return
	}

	ultimaData, err := c.Service.BuscarUltimaData()
	if err != nil {
		http.Error(w, "Erro ao ler ultima data", http.StatusInternalServerError)
		log.Printf("Erro leitura última data: %v", err)
		return
	}
	var Ultimo *string
	if ultimaData != nil {
		formatted := ultimaData.Format("02/01/2006")
		Ultimo = &formatted
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"total":           Total,
		"lojas_treinadas": TotalLojas,
		"ultimo_registro": Ultimo,
	})

}

func (c *OrientacaoController) SalvarHandler(w http.ResponseWriter, r *http.Request) { //Chama a função para salvar no banco
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	dataStr := r.FormValue("data_orientacao")
	dataformat, err := time.Parse("2006-01-02", dataStr)
	if err != nil {
		http.Error(w, "Data de orientação inválida ou ausente", http.StatusBadRequest)
		return
	}

	now := time.Now()

	novaOrientacao := models.OrientacaoEducativa{
		LojaID:              r.FormValue("loja_id"), //lê do html a var name="loja_id"
		ResponsavelPresente: r.FormValue("responsavel_presente"),
		FuncaoResponsavel:   r.FormValue("funcao_responsavel"),
		DataOrientacao:      dataformat,
		Observacoes:         r.FormValue("observacoes"),
		Signatario:          r.FormValue("signatario"),
		DataAssinatura:      &now,
	}

	err = c.Service.CriarNovaOrientacao(novaOrientacao)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Atualiza a página limpando o formulário
	w.WriteHeader(http.StatusOK)
}

func (c *OrientacaoController) EditarHandler(w http.ResponseWriter, r *http.Request) { //Envia solicitação para alterar banco
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
	signatario := r.FormValue("signatario")            // Pega do <input name="signatario">

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

	now := time.Now()

	// 5. Montando a nossa Struct com os dados limpos e convertidos
	orientacaoAtualizada := models.OrientacaoEducativa{
		ID:                  id,
		ResponsavelPresente: responsavel,
		FuncaoResponsavel:   funcao,
		DataOrientacao:      dataOrientacao,
		Observacoes:         observacoes,
		Signatario:          signatario,
		DataAssinatura:      &now,
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

func (c *OrientacaoController) DeleteHandler(w http.ResponseWriter, r *http.Request) { //Solicita o Delete de dados do Banco

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

/*func (c *OrientacaoController) ListarJSONHandler(w http.ResponseWriter, r *http.Request) { //Função para realizar a pesquisa na tabela
	orientacoes, err := c.Service.ListarTodas()
	if err != nil {
		http.Error(w, "Erro ao buscar orientações", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "aplication/json")
	json.NewEncoder(w).Encode(orientacoes)
}*/
