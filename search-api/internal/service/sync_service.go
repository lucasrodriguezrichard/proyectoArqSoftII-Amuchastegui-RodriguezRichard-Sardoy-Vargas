package service

import (
    "context"
    "log"

    "github.com/blassardoy/restaurant-reservas/search-api/internal/repository"
)

type SyncService struct {
    repo    repository.SearchRepository
    resClient *ReservationClient
}

func NewSyncService(repo repository.SearchRepository, resClient *ReservationClient) *SyncService {
    return &SyncService{repo: repo, resClient: resClient}
}

func (s *SyncService) HandleEvent(ctx context.Context, op string, id string) error {
    switch op {
    case "create", "update", "confirm":
        doc, err := s.resClient.GetReservationByID(id)
        if err != nil { return err }
        if op == "create" {
            return s.repo.Index(ctx, *doc)
        }
        return s.repo.Update(ctx, *doc)
    case "delete":
        return s.repo.Delete(ctx, id)
    default:
        log.Printf("unknown op: %s", op)
        return nil
    }
}

