package services

import (
	"github.com/Jaidenmagnan/waygates/models"
	"github.com/Jaidenmagnan/waygates/repositories"
)

type WaygateService struct {
	waygateRepository *repositories.WaygateRepository
}

func NewWaygateService(waygateRepository *repositories.WaygateRepository) *WaygateService {
	return &WaygateService{
		waygateRepository: waygateRepository,
	}
}

func (s *WaygateService) CreateWaygate(name string, userId int) (models.Waygate, error) {
	waygate, err := s.waygateRepository.Create(models.CreateWaygate{
		Name:   name,
		UserId: userId,
	})
	return waygate, err
}

func (s *WaygateService) GetWaygateByID(id int) (models.Waygate, error) {
	waygate, err := s.waygateRepository.GetByID(id)
	return waygate, err
}

func (s *WaygateService) ListUserWaygates(userId int) ([]models.Waygate, error) {
	return s.waygateRepository.GetByUserID(userId)
}
