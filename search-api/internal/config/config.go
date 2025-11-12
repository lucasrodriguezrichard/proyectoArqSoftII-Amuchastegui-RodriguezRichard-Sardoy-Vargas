package config

import (
    "os"
)

type AppConfig struct {
    // Solr
    SolrURL  string
    SolrCore string

    // Caches
    LocalCacheTTLSeconds int
    DistCacheTTLSeconds  int
    MemcachedAddrs       string // comma-separated

    // RabbitMQ
    RabbitMQURI      string
    RabbitMQQueue    string
    RabbitMQExchange string

    // Reservations API
    ReservationsAPIURL string

    // Server
    Port string
}

func getenv(key, def string) string {
    v := os.Getenv(key)
    if v == "" {
        return def
    }
    return v
}

func FromEnv() AppConfig {
    return AppConfig{
        SolrURL:              getenv("SOLR_URL", "http://localhost:8983/solr"),
        SolrCore:             getenv("SOLR_CORE", "reservations"),
        LocalCacheTTLSeconds: atoiDefault(getenv("CACHE_TTL_SECONDS", "300"), 300),
        DistCacheTTLSeconds:  atoiDefault(getenv("DIST_CACHE_TTL_SECONDS", "900"), 900),
        MemcachedAddrs:       getenv("MEMCACHED_ADDRS", "localhost:11211"),
        RabbitMQURI:          getenv("RABBITMQ_URI", "amqp://admin:admin@localhost:5672/"),
        RabbitMQQueue:        getenv("RABBITMQ_QUEUE", "reservations_updates"),
        RabbitMQExchange:     getenv("RABBITMQ_EXCHANGE", "restaurant_events"),
        ReservationsAPIURL:   getenv("RESERVATIONS_API_URL", "http://localhost:8081"),
        Port:                 getenv("PORT", "8082"),
    }
}

func atoiDefault(s string, def int) int {
    var n int
    for _, ch := range s {
        if ch < '0' || ch > '9' {
            return def
        }
        n = n*10 + int(ch-'0')
    }
    return n
}

