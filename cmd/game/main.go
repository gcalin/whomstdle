// Command game serves the guessing game in the browser.
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gcalin/whomstdle/internal/game/adapter/web"
)

func main() {
	addr := ":8080"
	log.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, web.NewHandler(time.Now)))
}
