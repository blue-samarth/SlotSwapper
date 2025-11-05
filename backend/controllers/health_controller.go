package controllers

import (
    "net/http"
    "os"
    "runtime"
    "time"

    "github.com/gin-gonic/gin"
    "SlotSwapper/config"
)

var startTime = time.Now()

// HealthResponse represents the health check response
type HealthResponse struct {
    Status      string            `json:"status"`
    Service     string            `json:"service"`
    Version     string            `json:"version"`
    Uptime      string            `json:"uptime"`
    Timestamp   string            `json:"timestamp"`
    Environment string            `json:"environment,omitempty"`
    Checks      map[string]string `json:"checks"`
}

// DetailedHealthResponse includes system metrics
type DetailedHealthResponse struct {
    Status      string                 `json:"status"`
    Service     string                 `json:"service"`
    Version     string                 `json:"version"`
    Uptime      string                 `json:"uptime"`
    Timestamp   string                 `json:"timestamp"`
    Environment string                 `json:"environment,omitempty"`
    Checks      map[string]string      `json:"checks"`
    System      map[string]interface{} `json:"system"`
}

// HealthCheck returns basic health status
func HealthCheck(c *gin.Context) {
    checks := make(map[string]string)
    
    // Check database connectivity
    db := config.GetDB()
    sqlDB, err := db.DB()
    if err != nil || sqlDB.Ping() != nil {
        checks["database"] = "unhealthy"
    } else {
        checks["database"] = "healthy"
    }
    
    // Determine overall status
    status := "healthy"
    for _, check := range checks {
        if check == "unhealthy" {
            status = "unhealthy"
            break
        }
    }
    
    response := HealthResponse{
        Status:      status,
        Service:     "SlotSwapper API",
        Version:     getVersion(),
        Uptime:      getUptime(),
        Timestamp:   time.Now().Format(time.RFC3339),
        Environment: getEnvironment(),
        Checks:      checks,
    }
    
    statusCode := http.StatusOK
    if status == "unhealthy" {
        statusCode = http.StatusServiceUnavailable
    }
    
    c.JSON(statusCode, response)
}

// DetailedHealthCheck returns detailed health with system metrics
func DetailedHealthCheck(c *gin.Context) {
    checks := make(map[string]string)
    
    // Check database connectivity
    db := config.GetDB()
    sqlDB, err := db.DB()
    if err != nil || sqlDB.Ping() != nil {
        checks["database"] = "unhealthy"
    } else {
        stats := sqlDB.Stats()
        checks["database"] = "healthy"
        checks["db_open_connections"] = string(rune(stats.OpenConnections))
        checks["db_in_use"] = string(rune(stats.InUse))
        checks["db_idle"] = string(rune(stats.Idle))
    }
    
    // System metrics
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    system := map[string]interface{}{
        "goroutines":       runtime.NumGoroutine(),
        "memory_alloc_mb":  m.Alloc / 1024 / 1024,
        "memory_sys_mb":    m.Sys / 1024 / 1024,
        "num_gc":           m.NumGC,
        "go_version":       runtime.Version(),
        "cpu_count":        runtime.NumCPU(),
    }
    
    // Determine overall status
    status := "healthy"
    for _, check := range checks {
        if check == "unhealthy" {
            status = "unhealthy"
            break
        }
    }
    
    response := DetailedHealthResponse{
        Status:      status,
        Service:     "SlotSwapper API",
        Version:     getVersion(),
        Uptime:      getUptime(),
        Timestamp:   time.Now().Format(time.RFC3339),
        Environment: getEnvironment(),
        Checks:      checks,
        System:      system,
    }
    
    statusCode := http.StatusOK
    if status == "unhealthy" {
        statusCode = http.StatusServiceUnavailable
    }
    
    c.JSON(statusCode, response)
}

// ReadinessCheck returns readiness status (for k8s)
func ReadinessCheck(c *gin.Context) {
    db := config.GetDB()
    sqlDB, err := db.DB()
    
    if err != nil || sqlDB.Ping() != nil {
        c.JSON(http.StatusServiceUnavailable, gin.H{
            "status": "not_ready",
            "reason": "database_unavailable",
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "status": "ready",
    })
}

// LivenessCheck returns liveness status (for k8s)
func LivenessCheck(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status": "alive",
    })
}

// Helper functions
func getVersion() string {
    version := os.Getenv("APP_VERSION")
    if version == "" {
        return "1.0.0"
    }
    return version
}

func getUptime() string {
    uptime := time.Since(startTime)
    return uptime.Round(time.Second).String()
}

func getEnvironment() string {
    env := os.Getenv("APP_ENV")
    if env == "" {
        return "development"
    }
    return env
}