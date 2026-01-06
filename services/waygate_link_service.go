package services

import (
	"github.com/Jaidenmagnan/waygates/models"
	"github.com/Jaidenmagnan/waygates/repositories"
)

type WaygateLinkService struct {
	waygateLinkRepository *repositories.WaygateLinkRepository
}

func NewWaygateLinkService(waygateLinkRepository *repositories.WaygateLinkRepository) *WaygateLinkService {
	return &WaygateLinkService{
		waygateLinkRepository: waygateLinkRepository,
	}
}

func (s *WaygateLinkService) CreateWaygateLink(name string, link string, waygateId int) (models.WaygateLink, error) {
	waygateLink, err := s.waygateLinkRepository.Create(models.CreateWaygateLink{
		Name:      name,
		Link:      link,
		WaygateId: waygateId,
	})

	return waygateLink, err
}

func (s *WaygateLinkService) ListWaygateLinks(waygateId int) ([]models.WaygateLink, error) {
	return s.waygateLinkRepository.GetByWaygateID(waygateId)
}

func (s *WaygateLinkService) DeleteWaygateLink(waygateLinkId int) error {
	return s.waygateLinkRepository.Delete(waygateLinkId)
}

func (s *WaygateLinkService) ResolveLink(waygateId int, linkId int) (string, string) {
	waygateLinks, err := s.ListWaygateLinks(waygateId)
	if err != nil {
		return "", ""
	}

	index := findIndexById(waygateLinks, linkId)

	nextLink := waygateLinks[(index+1)%len(waygateLinks)].Link
	prevLink := waygateLinks[(index-1)%len(waygateLinks)].Link

	return nextLink, prevLink
}

func findIndexById(waygateLinks []models.WaygateLink, waygateLinkId int) int {
	for i, wgl := range waygateLinks {
		if wgl.ID == waygateLinkId {
			return i
		}
	}
	return -1
}
