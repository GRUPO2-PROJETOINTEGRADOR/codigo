# 🗺️ Especificação Funcional e Arquitetural: Módulos de Conservação

Este documento estabelece as regras de negócio, o mapeamento de telas (Formulários/Botões) e a arquitetura de pacotes em GoLang para os módulos de **Segurança Alimentar** e **Orientações Educativas**.

---

## 1. Mapeamento de Telas vs. Banco de Dados (CRUD)

### 📋 Aba 1: Segurança Alimentar (`seguranca-alimentar.html`)

Este módulo gerencia o histórico de relatórios sanitários e inspeções das lojas do segmento de alimentação.

#### 1. Formulário de Cadastro / Edição (Create/Update)
* **Campos em Tela (Inputs do Usuário):**
  * Seleção da Loja (Dropdown alimentado via `SELECT * FROM lojas`) -> Envia `loja_id` (`VARCHAR`).
  * Data da Inspeção (Input tipo `date`) -> Envia `data_auditoria`.
  * Nome do Responsável (Input tipo `text`) -> Envia `responsavel_loja`.
  * Cargo do Responsável (Input tipo `text`) -> Envia `cargo_responsavel`.
  * Nota Final (Input tipo `number` de 0 a 10) -> Envia `nota`.
  * Relatório PDF (Input tipo `file` para upload) -> Envia o arquivo que o Go salvará na pasta `/uploads` e guardará o caminho em `anexo_tiller`.
* **Regra de Negócio no Go (Service Layer):**
  * O usuário **não digita** a classificação. O Go lê a nota e calcula automaticamente:
    * `Nota >= 9`: BOM
    * `Nota >= 7 e < 9`: REGULAR
    * `Nota >= 5 e < 7`: RUIM
    * `Nota < 5`: INACEITÁVEL
  * Após salvar com sucesso, dispara automaticamente um comando para a tabela `auditoria_eventos` com a ação `'CRIADO'` ou `'ALTERADO'`.

#### 2. Botões de Ação na Tabela de Listagem (Read/Delete)
* **Botão Visualizar Anexo:** Abre o PDF correspondente salvo em `anexo_tiller`.
* **Botão Editar:** Carrega os dados da linha de volta para o formulário.
* **Botão Excluir (Delete):** Dispara uma requisição DELETE enviando o ID da auditoria. Remove o registro e dispara um evento `'EXCLUIDO'` para o Log Consolidado.

---

### 🎓 Aba 2: Orientações Educativas (`orientacao-educativa.html`)

Este módulo registra as ações preventivas, treinamentos e notificações aplicadas aos lojistas.

#### 1. Formulário de Cadastro / Edição (Create/Update)
* **Campos em Tela (Inputs do Usuário):**
  * Seleção da Loja (Dropdown) -> Envia `loja_id` (`VARCHAR`).
  * Responsável Presente (Input tipo `text`) -> Envia `responsavel_presente`.
  * Cargo/Função (Input tipo `text`) -> Envia `funcao_responsavel`.
  * Data da Ação (Input tipo `date`) -> Envia `data_orientacao`.
  * Observações Pedagógicas (Textarea) -> Envia `observacoes`.
* **Regra de Negócio no Go (Service Layer):**
  * Gravação simples dos dados textuais.
  * Dispara automaticamente o Log Consolidado para a tabela `auditoria_eventos` marcando a entidade `'ORIENTACAO'`.

#### 2. Tabela de Listagem e Painel Lateral (Read/Delete)
* Exibe a linha do tempo das orientações realizadas.
* **Botão Excluir:** Remove a orientação caso tenha sido cadastrada de forma errônea.

---

## 2. Distribuição das Funcionalidades no Repositório Go

Seguindo estritamente o seu padrão arquitetural, veja a responsabilidade de cada arquivo e como as funções serão distribuídas para atender os dois módulos.

### 🗂️ Camada 1: Modelos (`/app/models`)
Contém apenas as Structs puras que espelham as tabelas do banco e os payloads de formulários.

* **`seguranca.go`**
  ```go
  type AuditoriaSeguranca struct {
      ID               int       `json:"id"`
      LojaID           string    `json:"loja_id"`
      DataAuditoria    time.Time `json:"data_auditoria"`
      ResponsavelLoja  string    `json:"responsavel_loja"`
      CargoResponsavel string    `json:"cargo_responsavel"`
      Nota             int       `json:"nota"`
      AnexoTiller      string    `json:"anexo_tiller"`
      Classificacao    string    `json:"classificacao"`
  }

  
### **🗂️ Camada 2: Repositórios (`/app/repositories`)**

Contém a execução crua das queries SQL (`SELECT`, `INSERT`, `UPDATE`, `DELETE`) usando o driver `database/sql`. **Nenhuma regra de negócio entra aqui.**

* **`seguranca_repository.go`**  
  * `SalvarAuditoria(ctx, models.AuditoriaSeguranca) error` \-\> Executa o `INSERT INTO auditorias_seguranca`.  
  * `ListarAuditorias() ([]models.AuditoriaSeguranca, error)` \-\> Executa o `SELECT`.  
  * `ExcluirAuditoria(id int) error` \-\> Executa o `DELETE`.  
* **`orientacao_repository.go`**  
  * `SalvarOrientacao(ctx, models.OrientacaoEducativa) error`  
  * `ListarOrientacoes() ([]models.OrientacaoEducativa, error)`  
  * `ExcluirOrientacao(id int) error`  
* **`log_repository.go`** *(Para o Painel Lateral)*  
  * `RegistrarEvento(lojaId string, entidade string, acao string) error` \-\> Roda o `INSERT INTO auditoria_eventos`.

### **🗂️ Camada 3: Serviços (`/app/services`)**

É o coração da aplicação. Processa as regras de negócio, calcula dados derivados e orquestra chamadas entre múltiplos repositórios (ex: salvar auditoria \+ gerar log de auditoria na mesma transação).

* **`seguranca_service.go`**  
  * `CriarNovaInspecao(dados models.AuditoriaSeguranca) error`  
    * *Lógica:* Roda as condicionais de `if/else` para definir se a `Classificacao` é BOM/REGULAR/RUIM baseado na nota.  
    * *Fluxo:* Chama `repoSeguranca.SalvarAuditoria()`. Se der certo, chama `repoLog.RegistrarEvento(dados.LojaID, "AUDITORIA", "CRIADO")`.  
* **`orientacao_service.go`**  
  * `RegistrarAcaoEducativa(dados models.OrientacaoEducativa) error`  
    * *Fluxo:* Chama `repoOrientacao.SalvarOrientacao()`. Se der certo, chama `repoLog.RegistrarEvento(dados.LojaID, "ORIENTACAO", "CRIADO")`.

### **🗂️ Camada 4: Controladores (`/app/controllers`)**

Captura as requisições HTTP da tela, lida com multipart/form-data (Upload de arquivos), faz o parse dos dados para as Structs, chama a camada de Serviço e decide qual HTML renderizar.

* **`seguranca_controller.go`**  
  * `HandleCadastrar(w http.ResponseWriter, r *http.Request)` \-\> Processa o formulário de envio, faz o upload do arquivo PDF para o servidor, monta a Struct e joga para o Service.  
  * `HandleListar(w http.ResponseWriter, r *http.Request)` \-\> Busca os dados do Service e injeta no template `seguranca-alimentar.html`.  
* **`orientacao_controller.go`**  
  * `HandleSalvar(w http.ResponseWriter, r *http.Request)`  
  * `HandleListar(w http.ResponseWriter, r *http.Request)`

### **🗂️ Camada 5: Rotas (`/app/routes`)**

Apenas conecta os endpoints HTTP aos métodos correspondentes dos controladores.

* **`routes.go`**  
* Go

func ConfigurarRotas() {

    // Segurança Alimentar

    http.HandleFunc("/conservacao/seguranca-alimentar", controllers.ListarAuditoriasHandler)

    http.HandleFunc("/conservacao/seguranca-alimentar/salvar", controllers.SalvarAuditoriaHandler)

    // Orientações Educativas

    http.HandleFunc("/conservacao/orientacao-educativa", controllers.ListarOrientacoesHandler)

    http.HandleFunc("/conservacao/orientacao-educativa/salvar", controllers.SalvarOrientacaoHandler)

}



## **3\. Fluxo de Execução Unificado (Como as camadas conversam)**

Para fixar a lógica de como essas pastas funcionam de forma independente, visualize o fluxo de um cadastro na aba de Segurança Alimentar:
```
\[Tela HTML: Clique em Salvar\]

       │

       ▼

\[controllers.SalvarAuditoriaHandler\] ──\> Processa arquivo PDF e faz parse do texto

       │

       ▼

\[services.SegurancaService\]          ──\> Calcula Classificação (BOM/RUIM) via lógica Go

       │

       ├───\> \[repositories.SegurancaRepository\] ──\> Salva registro no Banco de Dados

       │

       └───\> \[repositories.LogRepository\]       ──\> Grava o Log do Painel Lateral

       │

       ▼
\[Tela HTML\]                          ──\> Recarrega exibindo a tabela atualizada com o Log
```

Esta separação garante que, se no próximo período você decidir trocar o HTML do Frontend por um framework como React ou Vue, nenhuma linha de código SQL ou de regra de negócio precisará ser reescrita: você mudará apenas a camada de `controllers` para responder JSON em vez de renderizar arquivos `.html`.