package services

import (
	"codigo/app/models"
	"errors"
	"log"
	"time"
)

// OrientacaoRepository defines the interface for data persistence.
type OrientacaoRepository interface {
	Salvar(o models.OrientacaoEducativa) error
	ListarTodas() ([]models.OrientacaoEducativa, error)
	BuscaPorID(id int) (models.OrientacaoEducativa, error)
	Atualizar(o models.OrientacaoEducativa) error
	Delete(o models.OrientacaoEducativa) error
	TotalTreinos() (int, error)
	LojasTreinos() (int, error)
	BuscarUltimaData() (*time.Time, error)
}

// OrientacaoService provides business logic for Orientacao Educativa.
// It depends on an OrientacaoRepository for data persistence.
type OrientacaoService struct {
	Repo OrientacaoRepository
}

// CriarNovaOrientacao validates and saves a new OrientacaoEducativa.
func (s *OrientacaoService) CriarNovaOrientacao(o models.OrientacaoEducativa) error {
	if o.LojaID == "" || o.ResponsavelPresente == "" || o.Signatario == "" {
		return errors.New("Os campos 'Loja', 'Responsável Presente' e 'Signatário' são obrigatórios")
	}
	//if o.DataOrientacao.After(time.Now()) {
	//	return errors.New("Data de orientação não pode ser futura")
	//}
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
	if o.Signatario == "" {
		return errors.New("o signatário é obrigatório")
	}
	if o.ID <= 0 {
		return errors.New("id de orientação inválido para atualização")
	}

	// Se passou nas regras, manda o repositório fazer o trabalho sujo
	return s.Repo.Atualizar(o)
}

func (s *OrientacaoService) Delete(o models.OrientacaoEducativa) error {
	if o.ID == 0 {
		log.Printf("Campo ID vazio!")
		return errors.New("Campo ID vazio!")
	}
	return s.Repo.Delete(o)
}

func (s *OrientacaoService) TotalTreinos() (int, error) {
	return s.Repo.TotalTreinos()
}

func (s *OrientacaoService) LojasTreinos() (int, error) {
	return s.Repo.LojasTreinos()
}

func (s *OrientacaoService) BuscarUltimaData() (*time.Time, error) {
	return s.Repo.BuscarUltimaData()
}
