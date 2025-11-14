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
	// Carga opcional de variables desde .env para los servicios locales
	_ = godotenv.Load()

	// Estructura central con todas las variables de entorno necesarias
	cfg := config.FromEnv()

	// Inicializa la capa de caché distribuida (Memcached) y local para búsquedas
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

	// Ensambla cliente de Solr, repositorio y servicio de sincronización con RabbitMQ
	solrClient := solr.New(cfg.SolrURL, cfg.SolrCore)
	repo := repository.NewSolrRepository(solrClient)
	resClient := service.NewReservationClient(cfg.ReservationsAPIURL)
	syncSvc := service.NewSyncService(repo, resClient, dualCache)

	// Lanza en segundo plano el consumidor de eventos que sincroniza Solr cuando llegan mensajes
	ctx, _ := context.WithCancel(context.Background())
	go func() {
		if err := rabbitmq.NewConsumer(cfg.RabbitMQURI, cfg.RabbitMQExchange, cfg.RabbitMQQueue, syncSvc).Run(ctx); err != nil {
			log.Printf("rabbitmq consumer stopped: %v", err)
		}
	}()

	// Servicio de búsqueda HTTP + reindexación inicial para poblar Solr antes de atender tráfico
	searchSvc := service.NewSearchService(repo, dualCache, resClient)
	go func() {
		reindexCtx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		if err := searchSvc.Reindex(reindexCtx); err != nil {
			log.Printf("initial solr reindex failed: %v", err)
		} else {
			log.Printf("initial solr index populated")
		}
	}()
	r := httptransport.NewRouterWithService(searchSvc)

	// Inicia el servidor HTTP que atiende /search y /reservations/:id
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
