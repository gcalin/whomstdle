// Command game serves the guessing game in the browser.
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gcalin/whomstdle/internal/game/adapter/memory"
	"github.com/gcalin/whomstdle/internal/game/adapter/web"
	"github.com/gcalin/whomstdle/internal/game/app"
)

func main() {
	games := memory.NewGames()
	catalog := memory.NewCatalog()
	svc := app.New(games, catalog, time.Now)
	addr := ":8080"
	log.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, web.NewHandler(svc)))
}
