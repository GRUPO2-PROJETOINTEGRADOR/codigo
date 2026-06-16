-- 1. Tabela Central
CREATE TABLE IF NOT EXISTS lojas (
    id VARCHAR(100) PRIMARY KEY, -- Código da LUC da loja
    nome VARCHAR(100) NOT NULL, -- Nome Fantasia da loja
    categoria VARCHAR(50) NOT NULL -- Ex: 'ALIMENTACAO', 'VESTUARIO'
);

-- 2. Módulo Segurança Alimentar
CREATE TABLE IF NOT EXISTS auditorias_seguranca (
    id SERIAL PRIMARY KEY,
    loja_id VARCHAR(100) REFERENCES lojas(id),
    data_auditoria DATE NOT NULL,
    responsavel_loja VARCHAR(100),
    cargo_responsavel VARCHAR(100),
    nota INTEGER,
    anexo_tiller VARCHAR(255),
    classificacao VARCHAR(20),
    tipo_inspecao VARCHAR(20) DEFAULT 'Rotina',
    nc_grave BOOLEAN DEFAULT FALSE
);

-- 3. Eco Participantes
CREATE TABLE IF NOT EXISTS eco_participantes (
    loja_id VARCHAR(100) PRIMARY KEY REFERENCES lojas(id) ON DELETE CASCADE,
    status_participacao BOOLEAN DEFAULT TRUE,
    data_entrada DATE NOT NULL,
    data_saida DATE,
    anexo_eco_nome VARCHAR(255),
    anexo_eco_dados BYTEA
);

CREATE TABLE IF NOT EXISTS kit (
    id SERIAL PRIMARY KEY,
    loja_id VARCHAR(100) NOT NULL REFERENCES eco_participantes(loja_id) ON DELETE CASCADE,
    data_entrega_kit DATE,
    qnt_kit INTEGER
);

-- 3.1. Resíduos Eco
CREATE TABLE IF NOT EXISTS residuos_eco (
    id SERIAL PRIMARY KEY,
    loja_id VARCHAR(100) NOT NULL REFERENCES eco_participantes(loja_id) ON DELETE CASCADE,
    peso_kg DECIMAL(10,2) NOT NULL,
    data_coleta DATE NOT NULL,
    aproveitado BOOLEAN DEFAULT TRUE
);

-- 4. Módulo Orientação Educativa
CREATE TABLE IF NOT EXISTS orientacoes_educativas (
    id SERIAL PRIMARY KEY,
    loja_id VARCHAR(100) NOT NULL REFERENCES lojas(id),
    responsavel_presente VARCHAR(255) NOT NULL,
    funcao_responsavel VARCHAR(255),
    data_orientacao DATE NOT NULL,
    observacoes TEXT
);

-- 5. Log Consolidado do Painel Lateral
CREATE TABLE IF NOT EXISTS auditoria_eventos (
    id SERIAL PRIMARY KEY,
    loja_id VARCHAR(100) REFERENCES lojas(id),
    entidade VARCHAR(50),
    acao VARCHAR(20),
    data_evento TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
