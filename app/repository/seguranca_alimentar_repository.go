package utils

import (
	"codigo/app/models"
)

func CriarAuditoria(a models.SegurancaAlimentar) error {
	query := `
		INSERT INTO auditorias_seguranca 
		(loja_id, data_auditoria, responsavel_loja, cargo_responsavel, nota, anexo_tiller, classificacao)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := DB.Exec(
		query,
		a.LojaID,
		a.DataAuditoria,
		a.ResponsavelLoja,
		a.CargoResponsavel,
		a.Nota,
		a.AnexoTiller,
		a.Classificacao,
	)

	return err
}

func ListarAuditorias() ([]models.SegurancaAlimentar, error) {
	query := `
		SELECT id, loja_id, data_auditoria, responsavel_loja, cargo_responsavel, nota, anexo_tiller, classificacao
		FROM auditorias_seguranca
		ORDER BY data_auditoria DESC
	`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var auditorias []models.SegurancaAlimentar

	for rows.Next() {
		var a models.SegurancaAlimentar

		err := rows.Scan(
			&a.ID,
			&a.LojaID,
			&a.DataAuditoria,
			&a.ResponsavelLoja,
			&a.CargoResponsavel,
			&a.Nota,
			&a.AnexoTiller,
			&a.Classificacao,
		)

		if err != nil {
			return nil, err
		}

		auditorias = append(auditorias, a)
	}

	return auditorias, nil
}
