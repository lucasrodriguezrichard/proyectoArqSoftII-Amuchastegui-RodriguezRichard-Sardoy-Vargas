package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/blassardoy/restaurant-reservas/search-api/internal/cache"
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

	// Cache layers
	var distributed *cache.DistributedCache
	if addrs := sanitizeAddrs(cfg.MemcachedAddrs); len(addrs) > 0 {
		dist, err := cache.NewDistributedCache(addrs, time.Duration(cfg.DistCacheTTLSeconds)*time.Second)
		if err != nil {
			log.Printf("memcached disabled: %v", err)
		} else {
			distributed = dist
		}
	}
	cacheCodec := cache.JSONCodec{New: func() any { return &repository.SearchResult{} }}
	dualCache := cache.NewDual(time.Duration(cfg.LocalCacheTTLSeconds)*time.Second, distributed, cacheCodec)

	// Wire Solr repository and sync service
	solrClient := solr.New(cfg.SolrURL, cfg.SolrCore)
	repo := repository.NewSolrRepository(solrClient)
	resClient := service.NewReservationClient(cfg.ReservationsAPIURL)
	syncSvc := service.NewSyncService(repo, resClient, dualCache)

	// Start consumer in background
	ctx, _ := context.WithCancel(context.Background())
	go func() {
		if err := rabbitmq.NewConsumer(cfg.RabbitMQURI, cfg.RabbitMQExchange, cfg.RabbitMQQueue, syncSvc).Run(ctx); err != nil {
			log.Printf("rabbitmq consumer stopped: %v", err)
		}
	}()

	// HTTP router
	searchSvc := service.NewSearchService(repo, dualCache, resClient)
	r := httptransport.NewRouterWithService(searchSvc)

	addr := ":" + cfg.Port
	log.Printf("search-api listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func sanitizeAddrs(raw string) []string {
	fields := strings.Split(raw, ",")
	out := make([]string, 0, len(fields))
	for _, f := range fields {
		f = strings.TrimSpace(f)
		if f != "" {
			out = append(out, f)
		}
	}
	return out
}
