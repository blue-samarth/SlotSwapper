package middleware

import (
    "net/http"
    "sync"
    "time"

    "github.com/gin-gonic/gin"
)

type RateLimiter struct {
    visitors map[string]*Visitor
    mu       sync.RWMutex
    limit    int
    window   time.Duration
}


type Visitor struct {
    count      int
    lastReset  time.Time
    mu         sync.Mutex
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
    rl := &RateLimiter{
        visitors: make(map[string]*Visitor),
        limit:    limit,
        window:   window,
    }

    go rl.cleanupRoutine()
    
    return rl
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        ip := c.ClientIP()

        if !rl.allow(ip) {
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error": "Rate limit exceeded. Please try again later.",
                "retry_after": int(rl.window.Seconds()),
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}

func (rl *RateLimiter) allow(ip string) bool {
    rl.mu.Lock()
    visitor, exists := rl.visitors[ip]
    if !exists {
        visitor = &Visitor{
            count:     0,
            lastReset: time.Now(),
        }
        rl.visitors[ip] = visitor
    }
    rl.mu.Unlock()
    
    visitor.mu.Lock()
    defer visitor.mu.Unlock()
    
    now := time.Now()
    if now.Sub(visitor.lastReset) > rl.window {
        visitor.count = 0
        visitor.lastReset = now
    }
    
    if visitor.count >= rl.limit {
        return false
    }
    
    visitor.count++
    return true
}

func (rl *RateLimiter) cleanupRoutine() {
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()
    
    for range ticker.C {
        rl.mu.Lock()
        now := time.Now()
        for ip, visitor := range rl.visitors {
            visitor.mu.Lock()
            if now.Sub(visitor.lastReset) > rl.window*2 {
                delete(rl.visitors, ip)
            }
            visitor.mu.Unlock()
        }
        rl.mu.Unlock()
    }
}

func GlobalRateLimiter() gin.HandlerFunc {
    limiter := NewRateLimiter(100, 1*time.Minute)
    return limiter.Middleware()
}

func AuthRateLimiter() gin.HandlerFunc {
    limiter := NewRateLimiter(5, 1*time.Minute)
    return limiter.Middleware()
}

func SwapRateLimiter() gin.HandlerFunc {
    limiter := NewRateLimiter(20, 1*time.Minute)
    return limiter.Middleware()
}
