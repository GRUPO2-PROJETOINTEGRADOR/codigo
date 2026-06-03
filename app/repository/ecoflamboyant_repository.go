package repo

import (
	"codigo/app/models"
	"database/sql"
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
