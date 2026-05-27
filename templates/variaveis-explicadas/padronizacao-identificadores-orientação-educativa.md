# Padronização de Identificadores — Orientação Educativa

## Objetivo

Alinhar os nomes de campos HTML, variáveis JavaScript e payloads de fetch no arquivo `orientacao-educativa.html` com as colunas da tabela `orientacoes_educativas` definidas em `codigo/database/schema.sql`, eliminando a necessidade de conversões no backend Go.

## Schema de referência

```sql
CREATE TABLE orientacoes_educativas (
    id                  SERIAL PRIMARY KEY,
    loja_id             INTEGER NOT NULL REFERENCES lojas(id),
    responsavel_presente VARCHAR(255) NOT NULL,
    funcao_responsavel  VARCHAR(255),
    data_orientacao     DATE NOT NULL,
    observacoes         TEXT
);
```

## Alterações realizadas

### 1. ID do `<select>` de loja

| Arquivo | Antes | Depois |
|---------|-------|--------|
| `orientacao-educativa.html:213` | `id="loja-select"` | `id="loja_id"` |
| `orientacao-educativa.html:356` | `getElementById('loja-select')` | `getElementById('loja_id')` |

### 2. Variáveis no submit handler

| Antes | Depois | Motivo |
|-------|--------|--------|
| `const responsavel = ...` | `const responsavel_presente = ...` | coluna `responsavel_presente` |
| `const funcao = ...` | `const funcao_responsavel = ...` | coluna `funcao_responsavel` |
| `const data = ...` | `const data_orientacao = ...` | coluna `data_orientacao` |
| `const obs = ...` | `const observacoes = ...` | coluna `observacoes` |

### 3. Validação dos campos obrigatórios

A condição `if` que verifica campos vazios foi atualizada para usar os novos nomes:

```
!loja_id || !responsavel_presente || !funcao_responsavel || !data_orientacao || !observacoes
```

### 4. Payload do fetch

Antes:

```js
body: JSON.stringify({
    loja_id,
    responsavel_presente: responsavel,
    funcao_responsavel: funcao,
    data_orientacao: data,
    observacoes: obs
})
```

Depois (shorthand ES6, pois chave ≡ valor):

```js
body: JSON.stringify({
    loja_id,
    responsavel_presente,
    funcao_responsavel,
    data_orientacao,
    observacoes
})
```

### 5. Template strings na renderização da tabela

Os campos já eram lidos com os nomes do schema (ex.: `o.responsavel_presente`, `o.funcao_responsavel`, `o.data_orientacao`, `o.observacoes`), portanto **nenhuma alteração foi necessária** na renderização.

## Arquivo alterado

- `NOVO/orientacao-educativa.html`

## O que NÃO foi alterado

- Layout visual e funcionamento da interface permanecem idênticos.
- Apenas identificadores foram renomeados; nenhuma lógica de negócio foi modificada.

## Compatibilidade

- Backend Go (`main.go`) recebe agora payloads com chaves exatamente iguais às colunas do banco, permitindo `json.Unmarshal` direto sem mapeamento adicional.
- As respostas da API já utilizam os nomes canônicos, então a renderização no frontend continua funcionando sem adaptações.
