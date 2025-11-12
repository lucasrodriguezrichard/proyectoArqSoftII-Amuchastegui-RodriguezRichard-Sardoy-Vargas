package main

import (
	"context"
	"log"
	"time"

	"github.com/blassardoy/restaurant-reservas/search-api/internal/config"
	"github.com/blassardoy/restaurant-reservas/search-api/internal/rabbitmq"
	"github.com/blassardoy/restaurant-reservas/search-api/internal/repository"
	"github.com/blassardoy/restaurant-reservas/search-api/internal/service"
	"github.com/blassardoy/restaurant-reservas/search-api/internal/solr"
	httptransport "github.com/blassardoy/restaurant-reservas/search-api/internal/transport/http"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := config.FromEnv()

	// Wire Solr repository
	solrClient := solr.New(cfg.SolrURL, cfg.SolrCore)
	repo := repository.NewSolrRepository(solrClient)
	resClient := service.NewReservationClient(cfg.ReservationsAPIURL)
	syncSvc := service.NewSyncService(repo, resClient)

	// Start consumer in background
	ctx, _ := context.WithCancel(context.Background())
	go func() {
		if err := rabbitmq.NewConsumer(cfg.RabbitMQURI, cfg.RabbitMQExchange, cfg.RabbitMQQueue, syncSvc).Run(ctx); err != nil {
			log.Printf("rabbitmq consumer stopped: %v", err)
		}
	}()

	// HTTP router
	searchSvc := service.NewSearchService(repo, time.Duration(cfg.LocalCacheTTLSeconds)*time.Second)
	r := httptransport.NewRouterWithService(searchSvc)

	addr := ":" + cfg.Port
	log.Printf("search-api listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
