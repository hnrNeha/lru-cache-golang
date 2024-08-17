package main

import (
	"backend/cache"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    lruCache *cache.LRUCache

    cacheHits = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "lru_cache_hits_total",
        Help: "Total number of cache hits",
    })
    cacheMisses = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "lru_cache_misses_total",
        Help: "Total number of cache misses",
    })
)

func init() {
    prometheus.MustRegister(cacheHits)
    prometheus.MustRegister(cacheMisses)
}

func main() {
    lruCache = cache.NewLRUCache(1024)
    r := gin.Default()

    // CORS middleware
    r.Use(cors.Default())

    // Expose Prometheus metrics endpoint
    r.GET("/metrics", gin.WrapH(promhttp.Handler()))

    r.GET("/cache/:key", func(c *gin.Context) {
        key := c.Param("key")
        value, found := lruCache.Get(key)
        if found {
            cacheHits.Inc()
            c.JSON(http.StatusOK, gin.H{"value": value})
        } else {
            cacheMisses.Inc()
            c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
        }
    })

    r.POST("/cache", func(c *gin.Context) {
        var json struct {
            Key   string `json:"key"`
            Value string `json:"value"`
        }
        if err := c.BindJSON(&json); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
            return
        }
        lruCache.Set(json.Key, json.Value, 5*time.Second)
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })

    print("starting server")
    r.Run(":8080")
}
