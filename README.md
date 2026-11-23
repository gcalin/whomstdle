# whomstdle

A daily celebrity guessing game, built as a learning project. See
`Daily Celebrity Guessing Game — Learning Build Plan.md` for the full 13-phase
plan. This is **Phase 0 — Foundations**: a terminal version of the game and the
repository skeleton the later phases grow into.

## Layout

```
.
├── cmd/
│   └── game/           # Phase 0 terminal entrypoint + the hard-coded cast/questions
│       └── main.go
├── internal/
│   └── game/           # pure game rules — no terminal, web, or database code
│       ├── game.go
│       └── game_test.go
├── go_practice/        # Tour-of-Go scratch files (its own module; not part of the game)
├── go.mod              # module github.com/gcalin/whomstdle
├── NOTES.md            # the learning log — three sentences per phase
└── README.md
```

The split between `internal/game` (rules) and `cmd/game` (I/O + data) is the
point: the rules have no dependencies, so they are trivial to unit-test, and
Phase 1 can put a web handler in front of exactly the same package.

## Run it

```sh
go run ./cmd/game
```

Then type commands at the prompt:

```
ask alive        # answers a question about the secret person
guess Einstein   # guess who it is (case-insensitive, aliases work)
help             # list question ids and commands
quit
```

## Test it

```sh
go test ./...
```

`go_practice/` has its own `go.mod`, so it is a separate module and is skipped
by `go test ./...` (those scratch files each define their own `main` and can't
compile together — run them individually with `go run ./go_practice/<file>.go`).

## What's scaffolding vs. your Phase 0 work

Provided as a worked example so the repo compiles and tests pass from the start:
the `game` package, a handful of tests mirroring the plan's two required cases,
and a 5-person cast in `cmd/game`.

Still yours to do for Phase 0 (don't let the AI do it — review and extend
instead): grow the cast and tests, make sure you can explain slices vs. arrays,
interfaces, and Go's error handling, and keep `NOTES.md` up to date. Finish the
phase with a commit and the `v0.1-minimal` tag once Phase 1 is reached (Phase 0
itself has no tag in the plan).
