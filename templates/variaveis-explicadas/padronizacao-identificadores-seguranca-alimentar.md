# Padronização de Identificadores — Segurança Alimentar

## Objetivo

Alinhar os nomes de campos HTML, variáveis JavaScript e payloads de fetch no arquivo `seguranca-alimentar.html` com as colunas da tabela `auditorias_seguranca` definidas em `codigo/database/schema.sql`, eliminando a necessidade de conversões no backend Go.

## Schema de referência

```sql
CREATE TABLE IF NOT EXISTS auditorias_seguranca (
    id                  SERIAL PRIMARY KEY,
    loja_id             INT REFERENCES lojas(id),
    data_auditoria      DATE NOT NULL,
    responsavel_loja    VARCHAR(100),
    cargo_responsavel   VARCHAR(100),
    nota                INTEGER,
    anexo_tiller        VARCHAR(255),
    classificacao       VARCHAR(20)
);
```

## Alterações realizadas

### 1. Atributos `name` dos campos do formulário

| Linha | Antes | Depois |
|-------|-------|--------|
| 480 | `name="responsavel"` | `name="responsavel_loja"` |
| 485 | `name="funcao"` | `name="cargo_responsavel"` |
| 516 | `name="anexo_url"` | `name="anexo_tiller"` |

### 2. `data-sort` no cabeçalho da tabela

| Linha | Antes | Depois |
|-------|-------|--------|
| 371 | `data-sort="responsavel"` | `data-sort="responsavel_loja"` |

### 3. Variáveis no submit handler (`submit` event)

| Antes | Depois |
|-------|--------|
| `const responsavel = ...` | `const responsavel_loja = ...` |
| `const funcao = ...` | `const cargo_responsavel = ...` |

### 4. Validação dos campos obrigatórios

```
!responsavel_loja  (antes: !responsavel)
!cargo_responsavel (antes: !funcao)
```

### 5. Payload do fetch

Antes:

```js
const body = { loja_id, tipo, data_auditoria, responsavel, funcao, nota, anexo_url: pdfFileName.value };
```

Depois:

```js
const body = { loja_id, tipo, data_auditoria, responsavel_loja, cargo_responsavel, nota, anexo_tiller: pdfFileName.value };
```

### 6. Template strings na renderização da tabela e modal de detalhes

| Linha | Antes | Depois |
|-------|-------|--------|
| 781 | `i.responsavel` | `i.responsavel_loja` |
| 911, 982 | `i.responsavel` | `i.responsavel_loja` |
| 912, 981 | `i.funcao` | `i.cargo_responsavel` |
| 914, 983 | `i.anexo_url` | `i.anexo_tiller` |

### 7. Filtro de busca (search)

| Linha | Antes | Depois |
|-------|-------|--------|
| 781 | `(i.responsavel \|\| '').toLowerCase()` | `(i.responsavel_loja \|\| '').toLowerCase()` |

## Observações

- O campo `tipo` (tipo de inspeção) **não existe** na tabela `auditorias_seguranca` do `schema.sql`. O formulário coleta esse campo e o envia no payload, mas a coluna correspondente deverá ser adicionada ao schema ou o backend deverá tratá-lo separadamente.
- O campo `classificacao` não é enviado pelo formulário — é calculado no frontend a partir da `nota` e enviado pelo backend após o cálculo.
- O `id` dos inputs (`input-responsavel`, `input-funcao`) **não foi alterado** para evitar quebrar os `getElementById` correspondentes. Apenas o atributo `name` (usado no payload HTTP) foi ajustado.

## Arquivo alterado

- `NOVO/seguranca-alimentar.html`
