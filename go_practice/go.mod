// go_practice is a scratchpad for Tour-of-Go / Learn-Go-with-Tests exercises.
// Each file is its own `package main` with its own `func main`, so the folder
// can't compile as a single package — run them one at a time, e.g.
// `go run ./slice.go`.
//
// It has its own go.mod purely so these scratch files are a *separate* module
// and are therefore excluded from the root module's `go build ./...` and
// `go test ./...` (a Phase 0 "done when" check). Nothing here is imported by
// the game.
module github.com/gcalin/whomstdle/go_practice

go 1.27
