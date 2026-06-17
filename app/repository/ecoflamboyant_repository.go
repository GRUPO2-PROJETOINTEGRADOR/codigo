package repo

import (
	"codigo/app/models"
	"database/sql"
	"sync"
	"time"
)

var (
	ensureAuditTable     sync.Once
)

func garantirTabelaAuditoria(db *sql.DB) {
	ensureAuditTable.Do(func() {
		db.Exec(`CREATE TABLE IF NOT EXISTS auditoria_eventos (
			id SERIAL PRIMARY KEY,
			loja_id VARCHAR(100) REFERENCES lojas(id),
			entidade VARCHAR(50),
			acao VARCHAR(20),
			data_evento TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`)
	})
}

func ListarLojas(db *sql.DB) ([]models.Loja, error) {
	rows, err := db.Query("SELECT id, nome FROM lojas ORDER BY nome ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []models.Loja
	for rows.Next() {
		var l models.Loja
		if err := rows.Scan(&l.ID, &l.Nome); err != nil {
			return nil, err
		}
		lista = append(lista, l)
	}
	return lista, nil
}

func ListarLojasParticipantes(db *sql.DB) ([]models.Loja, error) {
	rows, err := db.Query(`SELECT l.id, l.nome FROM lojas l INNER JOIN eco_participantes ep ON ep.loja_id = l.id ORDER BY l.nome`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []models.Loja
	for rows.Next() {
		var l models.Loja
		if err := rows.Scan(&l.ID, &l.Nome); err != nil {
			return nil, err
		}
		lista = append(lista, l)
	}
	return lista, nil
}

func CriarParticipante(db *sql.DB, lojaID string, dataEntrada time.Time, dataSaida *time.Time, nomeAnexo string, dadosAnexo []byte) error {
	query := `INSERT INTO eco_participantes (loja_id, status_participacao, data_entrada, data_saida, anexo_eco_nome, anexo_eco_dados)
		VALUES ($1, TRUE, $2, $3, $4, $5)`
	_, err := db.Exec(query, lojaID, dataEntrada, dataSaida, nomeAnexo, dadosAnexo)
	return err
}

func InserirResiduo(db *sql.DB, lojaID string, dataColeta time.Time, pesoKG float64, aproveitado bool) error {
	_, err := db.Exec(`INSERT INTO residuos_eco (loja_id, data_coleta, peso_kg, aproveitado) VALUES ($1, $2, $3, $4)`,
		lojaID, dataColeta, pesoKG, aproveitado)
	return err
}

func ListarResiduos(db *sql.DB, dataInicio, dataFim, lojaID string, limit, offset int) ([]models.Residuo, error) {
	query := `SELECT r.id, l.nome, r.data_coleta, r.peso_kg, r.aproveitado
		FROM residuos_eco r
		JOIN lojas l ON l.id = r.loja_id
		WHERE ($1 = '' OR r.data_coleta >= $1::date)
		  AND ($2 = '' OR r.data_coleta <= $2::date)
		  AND ($3 = '' OR r.loja_id = $3)
		ORDER BY r.data_coleta DESC
		LIMIT CASE WHEN $4 > 0 THEN $4 ELSE NULL END OFFSET $5`
	rows, err := db.Query(query, dataInicio, dataFim, lojaID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []models.Residuo
	for rows.Next() {
		var r models.Residuo
		if err := rows.Scan(&r.ID, &r.LojaNome, &r.DataColeta, &r.PesoKG, &r.Aproveitado); err != nil {
			return nil, err
		}
		lista = append(lista, r)
	}
	return lista, nil
}

func InserirKit(db *sql.DB, lojaID string, dataEntregaKit time.Time, qntKit int) error {
	_, err := db.Exec(`INSERT INTO kit (loja_id, data_entrega_kit, qnt_kit) VALUES ($1, $2, $3)`,
		lojaID, dataEntregaKit, qntKit)
	return err
}

func ContarResiduos(db *sql.DB, dataInicio, dataFim, lojaID string) (int, error) {
	var total int
	err := db.QueryRow(`SELECT COUNT(*)
		FROM residuos_eco r
		JOIN lojas l ON l.id = r.loja_id
		WHERE ($1 = '' OR r.data_coleta >= $1::date)
		  AND ($2 = '' OR r.data_coleta <= $2::date)
		  AND ($3 = '' OR r.loja_id = $3)`, dataInicio, dataFim, lojaID).Scan(&total)
	return total, err
}

func ListarKits(db *sql.DB, limit, offset int) ([]models.Kit, error) {
	rows, err := db.Query(`SELECT k.id, l.nome, k.data_entrega_kit, k.qnt_kit
		FROM kit k
		JOIN lojas l ON l.id = k.loja_id
		ORDER BY k.data_entrega_kit DESC
		LIMIT CASE WHEN $1 > 0 THEN $1 ELSE NULL END OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []models.Kit
	for rows.Next() {
		var k models.Kit
		if err := rows.Scan(&k.ID, &k.LojaNome, &k.DataEntregaKit, &k.QntKit); err != nil {
			return nil, err
		}
		lista = append(lista, k)
	}
	return lista, nil
}

func ContarParticipantes(db *sql.DB) (int, error) {
	var total int
	err := db.QueryRow(`SELECT COUNT(*)
		FROM eco_participantes ep
		JOIN lojas l ON l.id = ep.loja_id`).Scan(&total)
	return total, err
}

func ContarKits(db *sql.DB) (int, error) {
	var total int
	err := db.QueryRow(`SELECT COUNT(*) FROM kit`).Scan(&total)
	return total, err
}

func ContarAuditoriasEventos(db *sql.DB) (int, error) {
	var total int
	err := db.QueryRow(`SELECT COUNT(*)
		FROM auditoria_eventos a
		JOIN lojas l ON l.id = a.loja_id`).Scan(&total)
	return total, err
}

func ContarLojasAtivas(db *sql.DB) (int, error) {
	var total int
	err := db.QueryRow(`SELECT COUNT(*) FROM eco_participantes WHERE status_participacao = TRUE`).Scan(&total)
	return total, err
}

func CrescimentoLojasPorMes(db *sql.DB) ([]models.PontoLojas, error) {
	rows, err := db.Query(`SELECT TO_CHAR(data_entrada, 'Mon/YY') AS mes, COUNT(*) AS entradas
		FROM eco_participantes WHERE status_participacao = TRUE
		GROUP BY DATE_TRUNC('month', data_entrada), mes
		ORDER BY DATE_TRUNC('month', data_entrada)`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pontos []models.PontoLojas
	acumulado := 0
	for rows.Next() {
		var mes string
		var entradas int
		if err := rows.Scan(&mes, &entradas); err != nil {
			return nil, err
		}
		acumulado += entradas
		pontos = append(pontos, models.PontoLojas{Mes: mes, Total: acumulado, Entradas: entradas})
	}
	return pontos, nil
}

func SomarTotalKits(db *sql.DB) (int, error) {
	var total int
	err := db.QueryRow(`SELECT COALESCE(SUM(qnt_kit), 0) FROM kit`).Scan(&total)
	return total, err
}

func FluxoKitsPorPeriodo(db *sql.DB) ([]models.PontoKits, error) {
	rows, err := db.Query(`SELECT TO_CHAR(data_entrega_kit, 'DD/MM') AS periodo,
		COUNT(*) AS entregas,
		COALESCE(SUM(qnt_kit), 0) AS total_unidades
		FROM kit
		WHERE data_entrega_kit IS NOT NULL
		GROUP BY DATE_TRUNC('day', data_entrega_kit), periodo
		ORDER BY DATE_TRUNC('day', data_entrega_kit)`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pontos []models.PontoKits
	for rows.Next() {
		var p models.PontoKits
		if err := rows.Scan(&p.Periodo, &p.Entregas, &p.Total); err != nil {
			return nil, err
		}
		pontos = append(pontos, p)
	}
	return pontos, nil
}

func ListarParticipantes(db *sql.DB, limit, offset int) ([]models.Participante, error) {
	query := `SELECT ep.loja_id, l.nome, ep.status_participacao, ep.data_entrada, ep.data_saida, ep.anexo_eco_nome
		FROM eco_participantes ep
		JOIN lojas l ON l.id = ep.loja_id
		ORDER BY ep.status_participacao DESC, ep.data_entrada DESC
		LIMIT CASE WHEN $1 > 0 THEN $1 ELSE NULL END OFFSET $2`

	rows, err := db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []models.Participante
	for rows.Next() {
		var p models.Participante
		var ns sql.NullString
		var nt sql.NullTime

		if err := rows.Scan(&p.LojaID, &p.LojaName, &p.Status, &p.DataEntrada, &nt, &ns); err != nil {
			return nil, err
		}

		if nt.Valid {
			p.DataSaida = &nt.Time
		}
		if ns.Valid {
			p.AnexoEcoNome = ns.String
		}

		lista = append(lista, p)
	}
	return lista, nil
}

func BuscarTermoPorLoja(db *sql.DB, lojaID string) (string, []byte, error) {
	var nome string
	var dados []byte
	err := db.QueryRow(`SELECT anexo_eco_nome, anexo_eco_dados FROM eco_participantes WHERE loja_id = $1`, lojaID).Scan(&nome, &dados)
	if err != nil {
		return "", nil, err
	}
	return nome, dados, nil
}

func ResumoResiduos(db *sql.DB) (totalGeral, totalAdubo, totalDescarte float64, err error) {
	err = db.QueryRow(`SELECT
		COALESCE(SUM(peso_kg), 0),
		COALESCE(SUM(CASE WHEN aproveitado = true THEN peso_kg ELSE 0 END), 0),
		COALESCE(SUM(CASE WHEN aproveitado = false THEN peso_kg ELSE 0 END), 0)
		FROM residuos_eco`).Scan(&totalGeral, &totalAdubo, &totalDescarte)
	return
}

func InativarLoja(db *sql.DB, lojaID string) error {
	_, err := db.Exec(`UPDATE eco_participantes SET status_participacao = FALSE, data_saida = CURRENT_DATE WHERE loja_id = $1`, lojaID)
	return err
}

func AtivarLoja(db *sql.DB, lojaID string) error {
	_, err := db.Exec(`UPDATE eco_participantes SET status_participacao = TRUE, data_saida = NULL WHERE loja_id = $1`, lojaID)
	return err
}

func InserirAuditoria(db *sql.DB, lojaID, entidade, acao string) error {
	garantirTabelaAuditoria(db)
	_, err := db.Exec(`INSERT INTO auditoria_eventos (loja_id, entidade, acao) VALUES ($1, $2, $3)`, lojaID, entidade, acao)
	return err
}

func ListarAuditoriasEventos(db *sql.DB, limit, offset int) ([]models.RegistroAuditoria, error) {
	rows, err := db.Query(`SELECT a.id, l.nome, a.entidade, a.acao, a.data_evento
		FROM auditoria_eventos a
		JOIN lojas l ON l.id = a.loja_id
		ORDER BY a.data_evento DESC
		LIMIT CASE WHEN $1 > 0 THEN $1 ELSE NULL END OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []models.RegistroAuditoria
	for rows.Next() {
		var r models.RegistroAuditoria
		if err := rows.Scan(&r.ID, &r.LojaNome, &r.Entidade, &r.Acao, &r.DataEvento); err != nil {
			return nil, err
		}
		lista = append(lista, r)
	}
	return lista, nil
}

func FluxoResiduosPorPeriodo(db *sql.DB) ([]models.PontoResiduos, error) {
	rows, err := db.Query(`SELECT TO_CHAR(data_coleta, 'DD/MM') AS periodo,
		COALESCE(SUM(CASE WHEN aproveitado = true THEN peso_kg ELSE 0 END), 0) AS peso_adubo,
		COALESCE(SUM(CASE WHEN aproveitado = false THEN peso_kg ELSE 0 END), 0) AS peso_descarte
		FROM residuos_eco
		WHERE data_coleta IS NOT NULL
		GROUP BY DATE_TRUNC('day', data_coleta), periodo
		ORDER BY DATE_TRUNC('day', data_coleta)`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []models.PontoResiduos
	for rows.Next() {
		var p models.PontoResiduos
		if err := rows.Scan(&p.Periodo, &p.PesoAdubo, &p.PesoDescarte); err != nil {
			return nil, err
		}
		lista = append(lista, p)
	}
	return lista, nil
}
