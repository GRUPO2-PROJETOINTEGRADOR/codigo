# Sistema de Gestao Integrada para Shopping Center - JP Mall / Eco Flamboyant

Este projeto e uma aplicacao desenvolvida em Go (GoLang) para a gestao de lojas, auditorias de seguranca alimentar, monitoramento de residuos do projeto de sustentabilidade (Eco Flamboyant) e orientacoes educativas para lojistas. O sistema integra um banco de dados PostgreSQL e possui tanto uma interface de linha de comando (CLI) em Go quanto arquivos estaticos de front-end para visualizacao e relatorios.

---

## Estrutura do Projeto

O repositorio esta estruturado da seguinte forma:

* **app/**: Contem os arquivos de codigo-fonte da aplicacao Go.
  * **main.go**: Ponto de entrada do sistema. Contem o menu CLI interativo para operacoes CRUD de lojas e a estrutura inicial para o servidor HTTP local.
  * **utils/**: Modulos auxiliares de persistencia de dados e conexao com o banco:
    * **connectDb.go**: Lida com a conexao e ping ao banco de dados PostgreSQL usando variaveis de ambiente carregadas via godotenv.
    * **criarbanco.go**: Le o arquivo SQL inicial para criar as tabelas necessarias no banco de dados.
    * **inserirloja.go**: Realiza a insercao de novas lojas no banco de dados.
    * **lerlojas.go**: Consulta e retorna a lista de lojas registradas.
    * **atualizar.go**: Atualiza os dados cadastrais de uma loja (LUC, nome ou categoria).
    * **delete.go**: Exclui uma loja cadastrada com base em seu identificador.
* **database/**: Contem os scripts de definicao do banco de dados.
  * **tabelas.sql**: Script SQL responsavel por estruturar todas as tabelas do banco de dados.
* **static/**: Arquivos front-end para exibicao dos modulos e dashboards.
  * **index.html**: Painel principal da interface web.
  * **css/** e **styles/**: Arquivos de estilizacao para a interface grafica.
  * **pages/**: Telas especificas do sistema:
    * **dashboard.html**: Dashboard de visualizacao geral.
    * **eco-flamboyant.html**: Tela do projeto de sustentabilidade e descarte de residuos.
    * **nova-inspecao.html**: Formulario para registro de auditorias e inspecoes.
    * **orientacao-educativa.html**: Tela de registro de orientacoes de boas praticas.
    * **relatorios.html** e **relatorio-seguranca-alimentar.html**: Exibicao de relatorios consolidados.

---

## Modelagem do Banco de Dados

O banco de dados PostgreSQL e estruturado no arquivo tabelas.sql com as seguintes entidades:

1. **Lojas (lojas)**:
   * id: Identificador unico da loja (LUC).
   * nome: Nome fantasia do estabelecimento.
   * categoria: Segmento da loja (por exemplo, 'ALIMENTACAO', 'VESTUARIO').

2. **Auditorias de Seguranca Alimentar (auditorias_seguranca)**:
   * id: Chave primaria da auditoria.
   * loja_id: Chave estrangeira referenciando a loja auditada.
   * data_auditoria: Data em que a inspecao foi realizada.
   * responsavel_loja: Nome da pessoa que acompanhou a auditoria na loja.
   * nota: Nota atribuida ao estabelecimento no relatorio.
   * classificacao: Classificacao gerada com base na nota (por exemplo: BOM, REGULAR, RUIM, INACEITAVEL).

3. **Eco Participantes (eco_participantes)**:
   * loja_id: Chave primaria e estrangeira referenciando lojas com exclusao em cascata.
   * status_participacao: Booleano indicando participacao ativa ou encerrada.
   * data_entrada: Data de ingresso da loja no projeto socioambiental.
   * data_saida: Data de desligamento, se aplicavel.
   * anexo_eco: Caminho para documentos anexos.

4. **Residuos Eco (residuos_eco)**:
   * id: Identificador do descarte.
   * loja_id: Referencia a loja em eco_participantes.
   * peso_kg: Peso dos residuos descartados para reciclagem/compostagem.
   * data_entrega: Data em que o descarte foi realizado.

5. **Orientacao Educativa (orientacoes_educativas)**:
   * id: Identificador da orientacao.
   * loja_id: Referencia a loja orientada.
   * responsavel_presente: Funcionario que recebeu as instrucoes.
   * funcao_responsavel: Cargo do colaborador na loja.
   * data_orientacao: Data de realizacao da atividade.
   * observacoes: Detalhes e anotacoes da orientacao.
   * data_criacao: Registro automatico de timestamp do sistema.

6. **Log de Eventos (auditoria_eventos)**:
   * id: Identificador do evento.
   * loja_id: Referencia a loja associada a mudanca.
   * entidade: Area alterada ('AUDITORIA', 'RESIDUO', 'ORIENTACAO').
   * acao: Tipo de operacao ('CRIADO', 'ALTERADO', 'EXCLUIDO').
   * data_evento: Timestamp automatico da acao.

---

## Configuracao do Ambiente

### Pre-requisitos
* GoLang instalado na maquina.
* Banco de dados PostgreSQL configurado e em execucao.

### Passos de Instalar e Configurar

1. **Configurar as Variaveis de Ambiente**:
   Crie um arquivo chamado `.env` na raiz do projeto contendo as definicoes do seu banco de dados PostgreSQL. Voce pode se basear no arquivo `.env.example`:
   ```env
   PORT=5432
   DB_HOST=localhost
   DB_USER=seu_usuario_postgres
   DB_PASS=sua_senha_postgres
   DATABASE_NAME=nome_do_seu_banco
   ```

2. **Instalar as Dependencias do Go**:
   No diretorio raiz do projeto, execute o comando abaixo para carregar os modulos necessarios:
   ```bash
   go mod tidy
   ```

---

## Como Executar a Aplicacao

### Modo CLI (Linha de Comando)
Atualmente, o ponto de entrada principal em `app/main.go` executa uma interface via terminal que permite realizar as operacoes CRUD diretamente na tabela de lojas.

Para iniciar, execute o seguinte comando no terminal:
```bash
go run app/main.go
```

Voce visualizara o menu abaixo no terminal:
```
JP MALL
MENU
1 - CRIAR LOJA
2 - EXIBIR TODAS AS LOJAS
3 - EDITAR DADOS DAS LOJAS
4 - DELETAR LOJA
5 - SAIR DO SISTEMA!
```

### Servidor HTTP Web (Opcional)
Se desejar disponibilizar a visualizacao das paginas web estaticas localizadas no diretorio `/static`, basta descomentar o bloco de codigo do servidor de arquivos em `app/main.go`:

```go
// Cria um servidor de arquivos que serve os arquivos da pasta "./static".
fileserver := http.FileServer(http.Dir("./static"))

// Associa o servidor de arquivos a rota raiz ("/").
http.Handle("/", fileserver)

// Inicia o servidor HTTP na porta 8081.
if err := http.ListenAndServe(":8081", nil); err != nil {
    log.Fatal(err)
}
```

Apos descomentar e rodar a aplicacao com `go run app/main.go`, voce podera abrir os dashboards e formularios atraves do navegador no endereco:
`http://localhost:8081/`
