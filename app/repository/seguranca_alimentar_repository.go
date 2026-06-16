package repo

import (
	"codigo/app/models"
	"errors"
)

func CriarAuditoria(a models.SegurancaAlimentar) error {
	query := `
		INSERT INTO auditorias_seguranca 
		(loja_id, data_auditoria, responsavel_loja, cargo_responsavel, nota, anexo_tiller, classificacao, tipo_inspecao, nc_grave, pdf_nome, pdf_tipo, pdf_tamanho, pdf_arquivo)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
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
		a.TipoInspecao,
		a.NCGrave,
		a.PDFNome,
		a.PDFTipo,
		a.PDFTamanho,
		a.PDFArquivo,
	)

	return err
}

func ListarAuditorias() ([]models.SegurancaAlimentar, error) {
	query := `
		SELECT 
			id, 
			loja_id, 
			data_auditoria, 
			responsavel_loja, 
			cargo_responsavel, 
			nota, 
			anexo_tiller, 
			classificacao,
			tipo_inspecao,
			nc_grave,
			COALESCE(pdf_nome, ''),
			COALESCE(pdf_tipo, ''),
			COALESCE(pdf_tamanho, 0)
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
			&a.TipoInspecao,
			&a.NCGrave,
			&a.PDFNome,
			&a.PDFTipo,
			&a.PDFTamanho,
		)

		if err != nil {
			return nil, err
		}

		auditorias = append(auditorias, a)
	}

	return auditorias, nil
}

func DeletarAuditoria(id int) error {
	query := `
		DELETE FROM auditorias_seguranca
		WHERE id = $1
	`

	result, err := DB.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("nenhuma auditoria foi encontrada")
	}

	return nil
}

func AtualizarAuditoria(a models.SegurancaAlimentar) error {
	var query string
	var err error
	var result interface {
		RowsAffected() (int64, error)
	}

	if len(a.PDFArquivo) > 0 {
		query = `
			UPDATE auditorias_seguranca
			SET
				loja_id = $1,
				data_auditoria = $2,
				responsavel_loja = $3,
				cargo_responsavel = $4,
				nota = $5,
				anexo_tiller = $6,
				classificacao = $7,
				tipo_inspecao = $8,
				nc_grave = $9,
				pdf_nome = $10,
				pdf_tipo = $11,
				pdf_tamanho = $12,
				pdf_arquivo = $13
			WHERE id = $14
		`
		result, err = DB.Exec(
			query,
			a.LojaID,
			a.DataAuditoria,
			a.ResponsavelLoja,
			a.CargoResponsavel,
			a.Nota,
			a.AnexoTiller,
			a.Classificacao,
			a.TipoInspecao,
			a.NCGrave,
			a.PDFNome,
			a.PDFTipo,
			a.PDFTamanho,
			a.PDFArquivo,
			a.ID,
		)
	} else {
		query = `
			UPDATE auditorias_seguranca
			SET
				loja_id = $1,
				data_auditoria = $2,
				responsavel_loja = $3,
				cargo_responsavel = $4,
				nota = $5,
				classificacao = $6,
				tipo_inspecao = $7,
				nc_grave = $8
			WHERE id = $9
		`
		result, err = DB.Exec(
			query,
			a.LojaID,
			a.DataAuditoria,
			a.ResponsavelLoja,
			a.CargoResponsavel,
			a.Nota,
			a.Classificacao,
			a.TipoInspecao,
			a.NCGrave,
			a.ID,
		)
	}
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("nenhuma auditoria foi encontrada")
	}

	return nil
}

func BuscarAnexoAuditoria(id int) (string, error) {
	var anexo string
	query := `
		SELECT anexo_tiller
		FROM auditorias_seguranca
		WHERE id = $1
	`
	err := DB.QueryRow(query, id).Scan(&anexo)
	if err != nil {
		return "", err
	}
	return anexo, nil
}

type PDFAuditoria struct {
	Nome        string
	Tipo        string
	Tamanho     int64
	Arquivo     []byte
	AnexoTiller string
}

func BuscarPDFAuditoria(id int) (PDFAuditoria, error) {
	var pdf PDFAuditoria
	query := `
		SELECT COALESCE(pdf_nome, ''), COALESCE(pdf_tipo, ''), COALESCE(pdf_tamanho, 0), pdf_arquivo, COALESCE(anexo_tiller, '')
		FROM auditorias_seguranca
		WHERE id = $1
	`
	err := DB.QueryRow(query, id).Scan(&pdf.Nome, &pdf.Tipo, &pdf.Tamanho, &pdf.Arquivo, &pdf.AnexoTiller)
	return pdf, err
}