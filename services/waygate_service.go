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

func (s *WaygateService) DeleteWaygate(id int) error {
	return s.waygateRepository.Delete(id)
}

func (s *WaygateService) CanUserAccessWaygate(userId, waygateId int) (bool, error) {
	waygate, err := s.waygateRepository.GetByID(waygateId)
	if err != nil {
		return false, err
	}

	if waygate.UserId != userId {
		return false, nil
	}

	return true, nil
}
