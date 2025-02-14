package main

import (
  "log"
  "os"
  "github.com/yungbote/slotter/backend/services/authorization/internal/api"
)

func main() {
  port := os.Getenv("PORT")
  if port == "" {
    port = "8081"
  }
  router := api.NewRouter()
  log.Printf("INFO: Authorization service listensing on port: %s", port)
  if err := router.Run(":" + port); err != nil {
    log.Fatalf("ERROR: Failed to run auth service: %v", err)
  }
}
