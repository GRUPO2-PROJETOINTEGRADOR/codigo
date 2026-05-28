package services

import (
	"codigo/app/models"
	utils "codigo/app/repository"
	"errors"
	"time"
)

// OrientacaoService provides business logic for Orientacao Educativa.
// It depends on an OrientacaoRepository for data persistence.
type OrientacaoService struct {
	Repo utils.OrientacaoRepository
}

// CriarNovaOrientacao validates and saves a new OrientacaoEducativa.
func (s *OrientacaoService) CriarNovaOrientacao(o models.OrientacaoEducativa) error {
	if o.LojaID == "" || o.ResponsavelPresente == "" {
		return errors.New("Os campos 'Loja' e 'Responsável Presente' não podem estar vazios")
	}
	if o.DataOrientacao.After(time.Now()) {
		return errors.New("Data de orientação não pode ser futura")
	}
	return s.Repo.Salvar(o)
}

// ListarTodas returns all OrientacaoEducativa records.
func (s *OrientacaoService) ListarTodas() ([]models.OrientacaoEducativa, error) {
	return s.Repo.ListarTodas()
}
