package services

import (
	"bytes"
	utils "codigo/app/repository"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-pdf/fpdf"
)

func GerarRelatorioPDF(db *sql.DB, dataInicio, dataFim, lojaID string) ([]byte, error) {
	residuos, err := utils.ListarResiduos(db, dataInicio, dataFim, lojaID, 0, 0)
	if err != nil {
		log.Printf("Erro ao buscar residuos para PDF: %v", err)
	}

	var (
		volumeTotal, volumeAdubo, volumeDescarte float64
		totalKits                                int
		lojasAtivas                              int
	)
	volumeTotal, volumeAdubo, volumeDescarte, err = utils.ResumoResiduos(db)
	if err != nil {
		log.Printf("Erro ao buscar resumo residuos PDF: %v", err)
	}

	totalKits, err = utils.SomarTotalKits(db)
	if err != nil {
		log.Printf("Erro ao buscar total kits PDF: %v", err)
	}

	lojasAtivas, err = utils.ContarLojasAtivas(db)
	if err != nil {
		log.Printf("Erro ao contar lojas ativas PDF: %v", err)
	}

	var taxa float64
	if volumeTotal > 0 {
		taxa = (volumeAdubo / volumeTotal) * 100
	}

	periodo := "Todo o período"
	if dataInicio != "" && dataFim != "" {
		periodo = dataInicio + " até " + dataFim
	} else if dataInicio != "" {
		periodo = "A partir de " + dataInicio
	} else if dataFim != "" {
		periodo = "Até " + dataFim
	}

	restaurante := "Todas as lojas"
	if lojaID != "" {
		var nome string
		db.QueryRow(`SELECT nome FROM lojas WHERE id = $1`, lojaID).Scan(&nome)
		if nome != "" {
			restaurante = nome
		}
	}

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(true, 20)
	pdf.AddPage()

	pdf.SetFont("Helvetica", "B", 20)
	pdf.CellFormat(190, 14, "Relatorio Eco Flamboyant", "", 1, "C", false, 0, "")

	pdf.SetFont("Helvetica", "", 10)
	pdf.SetTextColor(100, 100, 100)
	pdf.CellFormat(190, 6, "Periodo: "+periodo, "", 1, "C", false, 0, "")
	pdf.CellFormat(190, 6, "Gerado em: "+time.Now().Format("02/01/2006"), "", 1, "C", false, 0, "")
	pdf.CellFormat(190, 6, "Restaurante: "+restaurante, "", 1, "C", false, 0, "")

	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(8)

	pdf.SetFont("Helvetica", "B", 14)
	pdf.CellFormat(190, 10, "Metricas do Periodo", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	drawTableHeader(pdf, []string{"Indicador", "Valor"}, []float64{130, 60})

	pdf.SetFont("Helvetica", "", 10)
	drawMetricRow(pdf, "Volume Total Gerado", fmt.Sprintf("%.1f kg", volumeTotal))
	drawMetricRow(pdf, "Total Aproveitado", fmt.Sprintf("%.1f kg", volumeAdubo))
	drawMetricRow(pdf, "Total Descartado", fmt.Sprintf("%.1f kg", volumeDescarte))
	drawMetricRow(pdf, "Taxa de Aproveitamento", fmt.Sprintf("%.1f%%", taxa))
	drawMetricRow(pdf, "Total de Kits (Cestas)", fmt.Sprintf("%d unid", totalKits))
	drawMetricRow(pdf, "Lojas Parceiras Ativas", fmt.Sprintf("%d", lojasAtivas))

	pdf.Ln(10)

	pdf.SetFont("Helvetica", "B", 14)
	pdf.CellFormat(190, 10, "Dados Historicos", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	drawTableHeader(pdf, []string{"Loja", "Data", "Peso", "Destino"}, []float64{60, 40, 40, 50})

	pdf.SetFont("Helvetica", "", 9)
	for i, r := range residuos {
		if i%2 == 0 {
			pdf.SetFillColor(245, 245, 245)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}
		destino := "Descarte Comum"
		if r.Aproveitado {
			destino = "Adubo (Horta)"
		}
		data := r.DataColeta.Format("02/01/2006")
		peso := fmt.Sprintf("%.1f kg", r.PesoKG)
		drawTableRow(pdf, []string{r.LojaNome, data, peso, destino}, []float64{60, 40, 40, 50}, true)
	}

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar PDF: %w", err)
	}
	return buf.Bytes(), nil
}

func drawTableHeader(pdf *fpdf.Fpdf, headers []string, widths []float64) {
	pdf.SetFont("Helvetica", "B", 10)
	pdf.SetFillColor(139, 26, 26)
	pdf.SetTextColor(255, 255, 255)
	for i, h := range headers {
		pdf.CellFormat(widths[i], 8, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetTextColor(0, 0, 0)
}

func drawMetricRow(pdf *fpdf.Fpdf, label, value string) {
	pdf.SetFont("Helvetica", "", 10)
	pdf.SetFillColor(245, 245, 245)
	pdf.CellFormat(130, 7, "  "+label, "1", 0, "L", true, 0, "")
	pdf.CellFormat(60, 7, value, "1", 1, "C", true, 0, "")
}

func drawTableRow(pdf *fpdf.Fpdf, cells []string, widths []float64, fill bool) {
	pdf.SetFont("Helvetica", "", 9)
	for i, c := range cells {
		pdf.CellFormat(widths[i], 7, " "+c, "1", 0, "L", fill, 0, "")
	}
	pdf.Ln(-1)
}


