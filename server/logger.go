package server

import (
	"log"
)

func LogDebug(msg string) {
	log.Printf("[DEBUG] %s\n", msg)
}

func LogInfo(msg string) {
	log.Printf("[INFO] %s\n", msg)
}

func LogWarn(msg string) {
	log.Printf("[WARN] %s\n", msg)
}

func LogError(msg string) {
	log.Printf("[ERROR] %s\n", msg)
}
