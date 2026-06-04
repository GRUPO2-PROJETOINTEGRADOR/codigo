package repo

import (
	"codigo/app/models"
	"database/sql"
	"time"
)

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

func CriarParticipante(db *sql.DB, p models.Participante) error {
	query := `INSERT INTO eco_participantes (loja_id, status_participacao, data_entrada, data_saida, anexo_eco)
		VALUES ($1, TRUE, $2, $3, $4)`
	_, err := db.Exec(query, p.LojaID, p.DataEntrada, p.DataSaida, p.AnexoEco)
	return err
}

func InserirResiduo(db *sql.DB, lojaID string, dataColeta time.Time, pesoKG float64, aproveitado bool) error {
	_, err := db.Exec(`INSERT INTO residuos_eco (loja_id, data_coleta, peso_kg, aproveitado) VALUES ($1, $2, $3, $4)`,
		lojaID, dataColeta, pesoKG, aproveitado)
	return err
}

func ListarResiduos(db *sql.DB) ([]models.Residuo, error) {
	rows, err := db.Query(`SELECT r.id, l.nome, r.data_coleta, r.peso_kg, r.aproveitado
		FROM residuos_eco r
		JOIN lojas l ON l.id = r.loja_id
		ORDER BY r.data_coleta DESC`)
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

func ListarKits(db *sql.DB) ([]models.Kit, error) {
	rows, err := db.Query(`SELECT k.id, l.nome, k.data_entrega_kit, k.qnt_kit
		FROM kit k
		JOIN lojas l ON l.id = k.loja_id
		ORDER BY k.data_entrega_kit DESC`)
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

func ListarParticipantes(db *sql.DB) ([]models.Participante, error) {
	query := `SELECT ep.loja_id, l.nome, ep.status_participacao, ep.data_entrada, ep.data_saida, ep.anexo_eco
		FROM eco_participantes ep
		JOIN lojas l ON l.id = ep.loja_id
		ORDER BY ep.data_entrada DESC`

	rows, err := db.Query(query)
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
			p.AnexoEco = ns.String
		}

		lista = append(lista, p)
	}
	return lista, nil
}
