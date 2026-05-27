## ==============================================================================
## CHECKLIST DE IMPLEMENTAÇÃO BACKEND: MÓDULO ORIENTAÇÃO EDUCATIVA
## ==============================================================================
## Este documento serve como guia de tasks passo a passo para o desenvolvimento
## em camadas (MVC/Clean Architecture) integrado ao banco de dados PostgreSQL.
## ==============================================================================

---

## 🎯 SUGESTÕES DE ENGENHARIA ANTES DE COMEÇAR (Aditamentos ao Guia)

1. **Contexto de Renderização Duplo:** O Handler principal da página (`GET /conservacao/orientacao-educativa`) precisará carregar duas coisas em paralelo do banco de dados:
   - A lista de todas as orientações já gravadas (para a tabela).
   - A lista de todas as lojas cadastradas (para preencher o dropdown do formulário).
2. **Input Oculto no Formulário de Edição:** Para reaproveitar a mesma estrutura visual, utilize um campo oculto `<input type="hidden" name="id" value="{{.ID}}">` para que o backend diferencie uma "Criação" (ID zero/vazio) de uma "Edição" (ID preenchido).

---

## 📋 MATRIZ DE TASKS (Para Gerenciamento da Sprint 2)

### [ ] TASK 1: Camada de Modelos (`/app/models/orientacao.go`)
- [ ] Criar a Struct `OrientacaoEducativa` mapeando os tipos primitivos compatíveis com o Go e tags JSON.
- [ ] Garantir que o campo `LojaID` seja do tipo `string` (para suportar o formato VARCHAR/LUC das lojas).
- [ ] Criar uma Struct auxiliar `ContextoOrientacao` para envelopar os dados enviados à View (Lista de Orientações + Lista de Lojas).

### [ ] TASK 2: Camada de Repositório (`/app/repositories/orientacao_repository.go`)
- [ ] Implementar `SalvarOrientacao(o models.OrientacaoEducativa) error`.
- [ ] Implementar `ListarTodasOrientacoes() ([]models.OrientacaoEducativa, error)`.
- [ ] Implementar `BuscarOrientacaoPorID(id int) (models.OrientacaoEducativa, error)`.
- [ ] Implementar `AtualizarOrientacao(o models.OrientacaoEducativa) error`.
- [ ] Implementar `ExcluirOrientacao(id int) error`.
- *Nota: Usar queries parametrizadas (`$1`, `$2`) para prevenir SQL Injection.*

### [ ] TASK 3: Camada de Serviço (`/app/services/orientacao_service.go`)
- [ ] Validar regras de negócio básicas (ex: impedir texto de observação vazio ou data no futuro).
- [ ] Orquestrar o Log de Auditoria: Chamar `LogRepository.RegistrarEvento` após operações bem-sucedidas de INSERT, UPDATE e DELETE.

### [ ] TASK 4: Camada de Controladores (`/app/controllers/orientacao_controller.go`)
- [ ] **Handler Exibir Página (GET):** Buscar orientações e lojas do banco; mesclar no `ContextoOrientacao` e injetar no template HTML.
- [ ] **Handler Salvar/Editar (POST):** Ler dados via `r.FormValue()`, efetuar conversões de tipo (`strconv.Atoi` para o ID e `time.Parse` para data) e encaminhar para o Service.
- [ ] **Handler Excluir (POST):** Capturar o ID via formulário embutido e disparar a deleção.

### [ ] TASK 5: Camada de Roteamento (`/app/routes/routes.go`)
- [ ] Registrar os endpoints HTTP associados a cada Handler do Controlador.

### [ ] TASK 6: Integração Frontend (`/templates/conservacao/orientacao-educativa.html`)
- [ ] Implementar o laço dinâmico `{{range .Lojas}}` na tag `<select name="loja_id">`.
- [ ] Implementar o laço dinâmico `{{range .Orientacoes}}` nas linhas da tabela (`<tr>`).
- [ ] Adicionar os botões/fomulários de Editar e Excluir na tabela de listagem.

---

## 🛠️ CODIFICAÇÃO DAS CAMADAS: TEMPLATES ESQUELETO

### 📦 1. MODELO (`/app/models/orientacao.go`)
``` package models

import "time"

type OrientacaoEducativa struct {
	ID                  int       `json:"id"`
	LojaID              string    `json:"loja_id"` // LUC da Loja
	ResponsavelPresente string    `json:"responsavel_presente"`
	FuncaoResponsavel   string    `json:"funcao_responsavel"`
	DataOrientacao      time.Time `json:"data_orientacao"`
	Observacoes         string    `json:"observacoes"`
}

// Contexto para renderizar a página completa de forma dinâmica
type ContextoOrientacao struct {
	Orientacoes []OrientacaoEducativa
	Lojas       []Loja // Assumindo que você possui um Model Loja
}
```

### 🗄️ 2. REPOSITÓRIO (`/app/repositories/orientacao_repository.go`)

```go
package repositories

import (
	"database/sql"
	"PROTOTIPO/app/models"
	"PROTOTIPO/app/utils" // Local onde se encontra o seu pool database.DB
)

type OrientacaoRepository struct{}

func (repo *OrientacaoRepository) Salvar(o models.OrientacaoEducativa) error {
	query := `INSERT INTO orientacoes_educativas (loja_id, responsavel_presente, funcao_responsavel, data_orientacao, observacoes) 
	          VALUES ($1, $2, $3, $4, $5)`
	_, err := utils.DB.Exec(query, o.LojaID, o.ResponsavelPresente, o.FuncaoResponsavel, o.DataOrientacao, o.Observacoes)
	return err
}

func (repo *OrientacaoRepository) ListarTodas() ([]models.OrientacaoEducativa, error) {
	query := `SELECT id, loja_id, responsavel_presente, funcao_responsavel, data_orientacao, observacoes FROM orientacoes_educativas ORDER BY data_orientacao DESC`
	rows, err := utils.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []models.OrientacaoEducativa
	for rows.Next() {
		var o models.OrientacaoEducativa
		if err := rows.Scan(&o.ID, &o.LojaID, &o.ResponsavelPresente, &o.FuncaoResponsavel, &o.DataOrientacao, &o.Observacoes); err != nil {
			return nil, err
		}
		lista = append(lista, o)
	}
	return lista, nil
}

func (repo *OrientacaoRepository) Deletar(id int) error {
	query := `DELETE FROM orientacoes_educativas WHERE id = $1`
	_, err := utils.DB.Exec(query, id)
	return err
}

```

### 🧠 3. SERVIÇO (`/app/services/orientacao_service.go`)

```go
package services

import (
	"errors"
	"PROTOTIPO/app/models"
	"PROTOTIPO/app/repositories"
)

type OrientacaoService struct {
	repo    repositories.OrientacaoRepository
	repoLog repositories.LogRepository // Para alimentar o Painel Lateral automaticamente
}

func (s *OrientacaoService) CriarNovaOrientacao(o models.OrientacaoEducativa) error {
	if o.ResponsavelPresente == "" || o.LojaID == "" {
		return errors.New("os campos 'Loja Atendida' e 'Responsável' são obrigatórios")
	}

	err := s.repo.Salvar(o)
	if err != nil {
		return err
	}

	// Registra o evento de log de forma transparente para auditoria
	_ = s.repoLog.RegistrarEvento(o.LojaID, "ORIENTACAO", "CRIADO")
	return nil
}

```

### 🎮 4. CONTROLADOR (`/app/controllers/orientacao_controller.go`)

```go
package controllers

import (
	"html/template"
	"net/http"
	"strconv"
	"time"
	"PROTOTIPO/app/models"
	"PROTOTIPO/app/services"
)

type OrientacaoController struct {
	service services.OrientacaoService
}

// Renderiza a página misturando a tabela e o dropdown populado do banco
func (c *OrientacaoController) ListarPaginaHandler(w http.ResponseWriter, r *http.Request) {
	orientacoes, err := c.service.ListarTodas()
	if err != nil {
		http.Error(w, "Erro ao buscar orientações", http.StatusInternalServerError)
		return
	}

	// OBS: Aqui você também chamará o repositório de lojas para buscar []models.Loja
	var lojasFake []models.Loja // Substituir pela busca real das lojas cadastradas

	contexto := models.ContextoOrientacao{
		Orientacoes: orientacoes,
		Lojas:       lojasFake,
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/sidebar.html", "templates/conservacao/orientacao-educativa.html"))
	tmpl.Execute(w, contexto)
}

// Captura a requisição POST do formulário
func (c *OrientacaoController) SalvarHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	dataStr := r.FormValue("data_orientacao")
	dataConvertida, err := time.Parse("2006-01-02", dataStr)
	if err != nil {
		http.Error(w, "Formato de data inválido. Use AAAA-MM-DD", http.StatusBadRequest)
		return
	}

	novaOrientacao := models.OrientacaoEducativa{
		LojaID:              r.FormValue("loja_id"),
		ResponsavelPresente: r.FormValue("responsavel_presente"),
		FuncaoResponsavel:   r.FormValue("funcao_responsavel"),
		DataOrientacao:      dataConvertida,
		Observacoes:         r.FormValue("observacoes"),
	}

	err = c.service.CriarNovaOrientacao(novaOrientacao)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Redireciona para recarregar a tela limpando o fluxo do formulário
	http.Redirect(w, r, "/conservacao/orientacao-educativa", http.StatusSeeOther)
}

```

### 🛣️ 5. MAPEAMENTO DE ROTAS (`/app/routes/routes.go`)

```go
package routes

import (
	"net/http"
	"PROTOTIPO/app/controllers"
)

func CarregarRotasOrientacao(c *controllers.OrientacaoController) {
	// Exibe a tela com a tabela e o dropdown preenchidos
	http.HandleFunc("/conservacao/orientacao-educativa", c.ListarPaginaHandler)
	
	// Ações de processamento de dados (submissões de formulário)
	http.HandleFunc("/conservacao/orientacao-educativa/salvar", c.SalvarHandler)
	// http.HandleFunc("/conservacao/orientacao-educativa/excluir", c.ExcluirHandler)
}

```

---

## 🎨 6. ADAPTAÇÕES NO SEU HTML (`/templates/conservacao/orientacao-educativa.html`)

### A. População do Dropdown Dinâmico (Dentro do seu Form/Modal de Cadastro)

```html
<label for="loja_id">Loja Atendida *</label>
<select name="loja_id" id="loja_id" required>
    <option value="">Selecione a loja...</option>
    {{range .Lojas}}
        <option value="{{.ID}}">{{.Nome}} ({{.ID}})</option>
    {{end}}
</select>

```

### B. Injeção dos Registros e Botões de Ações na Tabela

```html
<tbody>
    {{range .Orientacoes}}
    <tr>
        <td>{{.LojaID}}</td>
        <td>{{.DataOrientacao.Format "02/01/2006"}}</td>
        <td>{{.ResponsavelPresente}}</td>
        <td>{{.FuncaoResponsavel}}</td>
        <td>{{.Observacoes}}</td>
        
        <td class="flex items-center gap-2">
            <a href="/conservacao/orientacao-educativa/editar?id={{.ID}}" class="p-1.5 text-blue-500 hover:bg-blue-50 rounded transition-colors" title="Editar Registro">
                <svg xmlns="[http://www.w3.org/2000/svg](http://www.w3.org/2000/svg)" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-pencil"><path d="M12 20h9"/><path d="M16.5 3.5a2.12 2.12 0 0 1 3 3L7 19l-4 1 1-4Z"/></svg>
            </a>

            <form action="/conservacao/orientacao-educativa/excluir" method="POST" onsubmit="return confirm('Deseja realmente remover esta orientação educativa?');" style="display:inline;">
                <input type="hidden" name="id" value="{{.ID}}">
                <button type="submit" class="p-1.5 text-gray-500 hover:text-[#D93030] rounded transition-colors" title="Excluir Registro">
                    <svg xmlns="[http://www.w3.org/2000/svg](http://www.w3.org/2000/svg)" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-trash-2"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/><line x1="10" x2="10" y1="11" y2="17"/><line x1="14" x2="14" y1="11" y2="17"/></svg>
                </button>
            </form>
        </td>
    </tr>
    {{else}}
    <tr>
        <td colspan="6" class="text-center text-gray-500 py-4">Nenhuma orientação educativa registrada até o momento.</td>
    </tr>
    {{end}}
</tbody>

```

# ==============================================================================

```

```