# Passo a Passo: Integração de Formulários e Tabelas (Front-end e Back-end)

Este guia explica detalhadamente como conectar o **front-end (HTML)** com os **handlers do Go** na arquitetura atual do sistema, usando a página de **Orientação Educativa** como exemplo.

A estrutura de integração utiliza os `templates/html` nativos da linguagem Go.

---

## 1. Como o Front-end recebe os Dados (Listagem na Tabela)

No seu `OrientacaoController.ListarPaginaHandler`, você envia para o template um `contexto`:

```go
contexto := models.ContextoOrientacao{
    Orientacoes: orientacoes, // Uma lista (slice) com todas as orientações do banco
}
tmpl.ExecuteTemplate(w, "layout", contexto)
```

Para exibir esses dados dinamicamente no HTML (`orientacao-educativa.html`), você deve utilizar a tag `{{range}}`, que funciona como um laço de repetição (`for`).

### Passo a Passo da Tabela:

1. Localize a tag `<tbody>` da sua tabela no HTML.
2. Remova as linhas estáticas (`<tr>`) usadas para mock de design.
3. Adicione a tag `{{range .Orientacoes}}` antes da `<tr>` que vai se repetir.
4. Substitua os dados estáticos pelos atributos da struct, usando `{{.NomeDoCampo}}`.
5. Feche o loop com `{{end}}`.

**Exemplo prático no HTML:**
```html
<tbody class="divide-y divide-border">
    <!-- Início do Loop Go -->
    {{range .Orientacoes}}
    <tr class="hover:bg-muted/50 transition">
        <td class="px-4 py-4 text-sm">{{.LojaID}}</td>
        <td class="px-4 py-4 text-sm">{{.ResponsavelPresente}}</td>
        <td class="px-4 py-4 text-sm">{{.FuncaoResponsavel}}</td>
        <!-- Para a data, pode ser necessário formatar no backend se vier suja -->
        <td class="px-4 py-4 text-sm">{{.DataOrientacao.Format "02/01/2006"}}</td>
        <td class="px-4 py-4 text-sm">{{.Observacoes}}</td>
        <td class="px-4 py-4">
            <!-- Botões de Ação -->
        </td>
    </tr>
    {{end}}
    <!-- Fim do Loop Go -->
</tbody>
```

---

## 2. Como o Front-end envia os Dados (Formulário de Cadastro)

Para que o formulário da página envie as informações corretamente para a rota `/conservacao/orientacao-educativa/salvar`, é preciso que as tags `<form>` e `<input>` (ou `<select>`) obedeçam às regras que o backend espera.

### Passo a Passo do Formulário:

1. **Tag Form:** Adicione as tags `action` e `method` no contêiner principal do formulário.
    * `action`: a URL exata do handler que salva os dados (`/conservacao/orientacao-educativa/salvar`).
    * `method`: deve ser `"POST"`.

```html
<form action="/conservacao/orientacao-educativa/salvar" method="POST" class="space-y-6">
```

2. **Atributo `name` nos Inputs:** O backend captura os dados usando `r.FormValue("nome_do_campo")`. Isso significa que cada `<input>` ou `<select>` no HTML deve possuir o atributo `name` exatamente igual à chave lida no Go.

Com base no seu `SalvarHandler`:
* `r.FormValue("data_orientacao")`
* `r.FormValue("loja_id")`
* `r.FormValue("responsavel_presente")`
* `r.FormValue("funcao_responsavel")`
* `r.FormValue("observacoes")`

**Exemplo prático no HTML:**
```html
<!-- Input Loja/LUC -->
<input type="text" name="loja_id" placeholder="Digite a LUC" required>

<!-- Input Responsável -->
<input type="text" name="responsavel_presente" required>

<!-- Input Função do Responsável -->
<input type="text" name="funcao_responsavel" required>

<!-- Input Data da Orientação -->
<input type="date" name="data_orientacao" required>

<!-- Textarea de Observações -->
<textarea name="observacoes"></textarea>

<!-- Botão de Envio (deve ser type="submit") -->
<button type="submit">Salvar Orientação</button>
```

---

## 3. Resumo de Integração

1. **Rotas e Handlers:** As rotas declaradas em `routes.go` conectam a URL digitada aos *Handlers*.
2. **Método GET (Listar):** O backend consulta o banco, monta uma *struct de contexto* e injeta no template. O front-end intercepta a struct com `{{.Variavel}}` ou iterando com `{{range}}`.
3. **Método POST (Salvar):** O front-end submete o formulário referenciando uma ação e um método. As tags `name=""` de cada elemento de formulário viajam no corpo da requisição e são lidas no Go usando `r.FormValue()`. 

Seguindo esses padrões em qualquer tela do sistema, a comunicação entre o banco de dados e a interface funcionará perfeitamente e de forma dinâmica.
