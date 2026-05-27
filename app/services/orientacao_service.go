package services

import (
	"codigo/app/models"
	utils "codigo/app/repository"
	"errors"
	"time"
)

type OrientacaoService struct {
	repo utils.OrientacaoRepository
}

func (s *OrientacaoService) CriarNovaOrientacao(o models.OrientacaoEducativa) error {

	if o.LojaID == "" || o.ResponsavelPresente == "" || o.DataOrientacao.After(time.Now()) { //É chamada para salvar, vai analisar se não há nenhum campo obrigatório vazio
		return errors.New("Os campos 'Loja' e 'Responsável Presente' não podem estar vazios!")
	}

	return s.repo.Salvar(o) //Se não houver nenhum campo vazio, ele executa a função repository de inserir no banco de dados
}

func (s *OrientacaoService) ListarTodas() ([]models.OrientacaoEducativa, error) {
	return s.repo.ListarTodas()
}
