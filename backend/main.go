package main

import (
	"context"
	"os"
	"os/signal"
	database "syncjob/Database"
	handler "syncjob/Handler"
	jobs "syncjob/Handler/Jobs"
	"syncjob/Logger"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	ADDR string = ":5461"
)

func main() {
	if err := database.LoadSMPM(); err != nil {
		logger.LogErr("%v", err)
	}
	defer database.Close()

	r := gin.Default()
	r.SetTrustedProxies([]string{"0.0.0.0"})

	logger.LogDebug("Starting server...")
	r.GET("api/pairs", handler.GetSlaveMasterPairs)
	r.GET("api/pairs/:key", handler.GetSlaveMasterPair)
	r.PUT("api/pairs/:key", handler.PutSlaveMasterPair)
	r.POST("api/pairs", handler.PostSlaveMasterPair)
	r.DELETE("api/pairs/:key", handler.DeleteSlaveMasterPair)

	r.POST("api/credentials/ping", handler.PostCredentialsPing)

	r.POST("api/pairs/:key/sync", handler.PostSlaveMasterPairSync)

	jobs.Init()

	go func() {
		logger.LogInfo("Server running at: http://localhost%s\n", ADDR)
		if err := r.Run(ADDR); err != nil {
			logger.LogFatal("HTTP server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	_, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()
}
