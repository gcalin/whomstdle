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

// newHandler builds the whole game: storage, app and web. Tests call it with a fixed clock.
func newHandler(now func() time.Time) http.Handler {
	return web.NewHandler(app.New(memory.NewGames(), memory.NewCatalog(), now))
}

func main() {
	addr := ":8080"
	log.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, newHandler(time.Now)))
}
