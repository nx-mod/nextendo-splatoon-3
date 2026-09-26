package main

// debugLogf — verbose diagnostics, off unless asked for.
//
// This existed only as call sites. 29 of them across frames_diag.go, stats.go,
// matchmaking.go, main.go and gamesync.go, and no definition in any revision of
// this repository: "Update to latest source" (74f2401) brought the callers in
// without the file that declares them, so the package has never compiled since.
//
// Gated on an environment variable, the way allowUnverified() and
// NPLN_CAPTURED_USER already are here, because the traffic it prints is
// per-HTTP/2-frame and would bury the ordinary log. frames_diag.go explains why
// that tracing has to be hand-rolled at all: grpc-go carries its own HTTP/2
// stack, so GODEBUG=http2debug produces nothing for it.
//
//   NPLN_DEBUG=1 ./server

import (
	"log"
	"os"
	"strings"
	"sync"
)

var (
	debugOnce    sync.Once
	debugEnabled bool
)

// debugLogf writes one formatted diagnostic line when NPLN_DEBUG is set to
// anything other than empty, "0" or "false". The check is done once: these call
// sites sit in per-frame paths where a Getenv each time would show up.
func debugLogf(format string, args ...any) {
	debugOnce.Do(func() {
		switch strings.ToLower(strings.TrimSpace(os.Getenv("NPLN_DEBUG"))) {
		case "", "0", "false", "no", "off":
			debugEnabled = false
		default:
			debugEnabled = true
		}
	})
	if !debugEnabled {
		return
	}
	log.Printf(format, args...)
}
