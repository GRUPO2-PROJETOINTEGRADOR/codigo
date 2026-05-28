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

func (s *OrientacaoService) BuscaPorID(id int) (models.OrientacaoEducativa, error) {
	return s.Repo.BuscaPorID(id)
}

func (s *OrientacaoService) Atualizar(o models.OrientacaoEducativa) error {
	// Exemplo de validação simples antes de salvar
	if o.ResponsavelPresente == "" {
		return errors.New("o nome do responsável presente é obrigatório")
	}
	if o.ID <= 0 {
		return errors.New("id de orientação inválido para atualização")
	}
	if o.DataOrientacao.After(time.Now()) {
		return errors.New("Data de orientação não pode ser futura")
	}

	// Se passou nas regras, manda o repositório fazer o trabalho sujo
	return s.Repo.Atualizar(o)
}
