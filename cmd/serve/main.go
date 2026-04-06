// Minimal GND server using Hooks() mode.
// This replaces the old mcpconfig-based server.
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

	dsr "github.com/dpopsuev/origami-gnd"
	fwmcp "github.com/dpopsuev/origami/mcp"
)

func main() {
	port := flag.Int("port", 9100, "HTTP port")
	healthz := flag.Bool("healthz", false, "probe /healthz and exit")
	flag.Parse()

	if *healthz {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/healthz", *port))
		if err != nil || resp.StatusCode != http.StatusOK {
			os.Exit(1)
		}
		os.Exit(0)
	}

	os.Exit(serve(*port))
}

func serve(port int) int {
	factory := dsr.Factory()
	cfg := fwmcp.SessionFactoryToConfig(factory)
	cfg.Name = "origami-gnd"
	cfg.Version = "1.0"

	server := fwmcp.NewCircuitServer(&cfg)
	defer server.Shutdown()

	mux := http.NewServeMux()
	mcpHandler := sdkmcp.NewStreamableHTTPHandler(
		func(_ *http.Request) *sdkmcp.Server { return server.MCPServer },
		&sdkmcp.StreamableHTTPOptions{Stateless: false},
	)
	mux.Handle("/mcp", mcpHandler)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(200) })

	addr := fmt.Sprintf(":%d", port)
	fmt.Fprintf(os.Stderr, "gnd listening on %s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Fprintf(os.Stderr, "gnd: %v\n", err)
		return 1
	}
	return 0
}
