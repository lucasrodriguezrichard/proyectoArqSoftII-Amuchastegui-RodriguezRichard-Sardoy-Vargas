package service

import (
	"context"
	"log"

	"github.com/blassardoy/restaurant-reservas/search-api/internal/cache"
	"github.com/blassardoy/restaurant-reservas/search-api/internal/repository"
)

type SyncService struct {
	repo      repository.SearchRepository
	resClient *ReservationClient
	cache     *cache.DualCache
}

func NewSyncService(repo repository.SearchRepository, resClient *ReservationClient, cacheLayer *cache.DualCache) *SyncService {
	return &SyncService{repo: repo, resClient: resClient, cache: cacheLayer}
}

func (s *SyncService) HandleEvent(ctx context.Context, op string, id string) error {
	var err error
	switch op {
	case "create", "update", "confirm":
		doc, derr := s.resClient.GetReservationByID(id)
		if derr != nil {
			return derr
		}
		if op == "create" {
			err = s.repo.Index(ctx, *doc)
		} else {
			err = s.repo.Update(ctx, *doc)
		}
	case "delete":
		err = s.repo.Delete(ctx, id)
	default:
		log.Printf("unknown op: %s", op)
		return nil
	}
	if err == nil && s.cache != nil {
		s.cache.Clear()
	}
	return err
}
