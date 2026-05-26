-- 1. Tabela Central
CREATE TABLE IF NOT EXISTS lojas (
    id VARCHAR(100) PRIMARY KEY, --Código da LUC da loja
    nome VARCHAR(100) NOT NULL, -- Nome Fantasia da loja
    categoria VARCHAR(50) NOT NULL -- Ex: 'ALIMENTACAO', 'VESTUARIO'
);

-- 2. Módulo Segurança Alimentar
CREATE TABLE IF NOT EXISTS auditorias_seguranca (
    id SERIAL PRIMARY KEY, --Id da auditoria proporção (loja 1:N auditorias)
    loja_id INT REFERENCES lojas(id), --É atrelado ao Id da loja cadastrada na tabela lojas
    data_auditoria DATE NOT NULL, --O usuário vai cadastrar o dia da auditoria realizada
    responsavel_loja VARCHAR(100), --Cadastro do nome do responsável pela loja que recebeu a inspeção
	cargo_responsavel VARCHAR(100),
    nota INTEGER, -- Nota da loja recebida pelo relatório Tiller
	anexo_tiller VARCHAR(255), -- Anexo do caminho do arquivo
    classificacao VARCHAR(20) -- INACEITAVEL, RUIM, REGULAR, BOM -- Classificação a ser exibida a depender do intervalo de nota
);

-- 3. Eco Participantes agora usa o loja_id como CHAVE PRIMÁRIA e CHAVE ESTRANGEIRA ao mesmo tempo
CREATE TABLE IF NOT EXISTS eco_participantes (
    loja_id INTEGER PRIMARY KEY REFERENCES lojas(id) ON DELETE CASCADE, --Id da loja participante atrelado à lojas
    status_participacao BOOLEAN DEFAULT TRUE, -- TRUE = Ativo, FALSE = Encerrado -- Status de participação da loja ao projeto eco flamboyant
	data_entrada DATE NOT NULL, --Inserir a data em que a loja entrou no projeto eco flamboyant
	data_saida DATE, --Caso a loja seja desligada do projeto a data será salva
	anexo_eco VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS kit (
    id SERIAL PRIMARY KEY, --Gera um id único de registro para o resíduo
    loja_id INTEGER NOT NULL REFERENCES eco_participantes(loja_id) ON DELETE CASCADE, --é atrlado único e exclusivamente às lojas participantes do eco_flamboyant
	data_entrega_kit DATE,
	qnt_kit INTEGER
);

-- 3.1. Resíduos Eco aponta DIRETO para o loja_id
CREATE TABLE IF NOT EXISTS residuos_eco (
    id SERIAL PRIMARY KEY, --Gera um id único de registro para o resíduo
    loja_id INTEGER NOT NULL REFERENCES eco_participantes(loja_id) ON DELETE CASCADE, --é atrlado único e exclusivamente às lojas participantes do eco_flamboyant
    peso_kg DECIMAL(10,2) NOT NULL, --Guarda a quantidade de Kg, descartados por cada loja
    data_coleta DATE NOT NULL, --Armazena a data da geração dos resíduos
	aproveitado BOOLEAN DEFAULT TRUE
);

-- 4. Módulo Orientação Educativa
CREATE TABLE IF NOT EXISTS orientacoes_educativas (
    id SERIAL PRIMARY KEY, --Gera um id único para cada registro
    loja_id INTEGER NOT NULL REFERENCES lojas(id), --Lojas de orientação educativa devem obrigatoriamente estar cadastrada
    responsavel_presente VARCHAR(255) NOT NULL, --Armazena o responsável da loja que recebeu a orientação
    funcao_responsavel VARCHAR(255), -- Armazena o "cargo" do responsável na loja
    data_orientacao DATE NOT NULL, -- Armazena a data de
    observacoes TEXT --registro das anotações sobre a orientação
    --data_criacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP --Registra automaticamente a data da criação da orientação educativa
);

-- 5. Log Consolidado do Painel Lateral
CREATE TABLE IF NOT EXISTS auditoria_eventos (
    id SERIAL PRIMARY KEY, --Id único para cada alteração no sistema
    loja_id INT REFERENCES lojas(id), -- Cada alteração é sempre atrelada a uma loja para rastreabilidade
    entidade VARCHAR(50), -- 'AUDITORIA', 'RESIDUO', 'ORIENTACAO'
    acao VARCHAR(20),     -- 'CRIADO', 'ALTERADO', 'EXCLUIDO'
    data_evento TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Registra automaticamente a criação de um evento
);