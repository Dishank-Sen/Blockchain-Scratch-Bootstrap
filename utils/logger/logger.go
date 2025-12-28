package logger

import (
	"context"
	"fmt"
	"log"
	"os"
)

func Info(msg string){
	log.Printf("Info : %s", msg)
}

func Error(ctx context.Context, cancel context.CancelFunc, msg string) {
	fmt.Println("Error[CLI]:", msg)

	// trigger cancellation so all goroutines watching ctx.Done() exit safely
	if cancel != nil {
		cancel()
	}

	// let defers run in the caller, then exit
	os.Exit(1)
}