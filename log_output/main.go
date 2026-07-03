package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func main() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	id := uuid.New().String()

	for {
		now := time.Now().UTC().Format(time.RFC3339Nano)
		fmt.Printf("%s: %s\n", now, id)
		<-ticker.C
	}
}
