package sasrd

import (
	"log"
	"os"
)

func logWarning(msg string, args ...interface{}) {
	log.Printf("[WARN] "+msg, args)
}

func logInfo(msg string, args ...interface{}) {
	log.Printf("[INFO] "+msg, args)
}

func logError(msg string, args ...interface{}) {
	log.Printf("[ERROR] "+msg, args)
}

func logFatal(msg string, args ...interface{}) {
	log.Printf("[ERROR] "+msg, args)
	os.Exit(1)
}
