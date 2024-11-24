package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jraams/bigtable-emulator-dumper/bigtable"
	"github.com/jraams/bigtable-emulator-dumper/config"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()

	bigtableService, err := bigtable.New(ctx, cfg)
	if err != nil {
		log.Fatalf("Could not init bigtable: %s", err.Error())
	}
	defer bigtableService.Close()

	r := gin.Default()

	err = r.SetTrustedProxies(nil)
	if err != nil {
		log.Fatal("Could not set trusted proxies to nil")
	}

	r.GET("/", func(c *gin.Context) {
		data, err := bigtableService.FetchAllTables(c)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, data)
	})

	r.GET("/:tableName", func(c *gin.Context) {
		tableName := c.Param("tableName")

		data, err := bigtableService.FetchSingleTable(c, tableName)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, data)
	})

	if err := r.Run(cfg.Address); err != nil {
		log.Fatal(err)
	}
}
