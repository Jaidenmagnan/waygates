package services

import (
	"errors"
	"strings"

	"github.com/Jaidenmagnan/waygates/models"
	"github.com/Jaidenmagnan/waygates/repositories"
)

var ErrWaygateLinkNotFound = errors.New("waygate link not found")

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
		Link:      normalizeLink(link),
		WaygateId: waygateId,
	})

	return waygateLink, err
}

func normalizeLink(link string) string {
	link = strings.TrimSpace(link)
	if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") {
		return link
	}

	return "https://" + link
}

func (s *WaygateLinkService) ListWaygateLinks(waygateId int) ([]models.WaygateLink, error) {
	return s.waygateLinkRepository.GetByWaygateID(waygateId)
}

func (s *WaygateLinkService) DeleteWaygateLink(waygateLinkId int) error {
	return s.waygateLinkRepository.Delete(waygateLinkId)
}

func (s *WaygateLinkService) ResolveLink(waygateId int, linkId int) (string, string) {
	nextLink, prevLink, _, err := s.ResolveNeighbors(waygateId, linkId)
	if err != nil {
		return "", ""
	}

	return nextLink.Link, prevLink.Link
}

func (s *WaygateLinkService) ResolveNeighbors(waygateId int, linkId int) (models.WaygateLink, models.WaygateLink, models.WaygateLink, error) {
	waygateLinks, err := s.ListWaygateLinks(waygateId)
	if err != nil {
		return models.WaygateLink{}, models.WaygateLink{}, models.WaygateLink{}, err
	}

	if len(waygateLinks) == 0 {
		return models.WaygateLink{}, models.WaygateLink{}, models.WaygateLink{}, ErrWaygateLinkNotFound
	}

	index := findIndexById(waygateLinks, linkId)
	if index == -1 {
		return models.WaygateLink{}, models.WaygateLink{}, models.WaygateLink{}, ErrWaygateLinkNotFound
	}

	nextLink := waygateLinks[(index+1)%len(waygateLinks)]
	prevLink := waygateLinks[(index-1+len(waygateLinks))%len(waygateLinks)]
	currentLink := waygateLinks[index]

	return nextLink, prevLink, currentLink, nil
}

func findIndexById(waygateLinks []models.WaygateLink, waygateLinkId int) int {
	for i, wgl := range waygateLinks {
		if wgl.ID == waygateLinkId {
			return i
		}
	}
	return -1
}
