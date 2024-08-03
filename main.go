package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"sniffle/api/eu_ec_ice"
	"sniffle/myhandler"
	"sniffle/tool"
)

func main() {
	logger := slog.New(myhandler.New(os.Stdout, slog.LevelDebug))

	fetcher := tool.FallBackFetcher(logger, http.DefaultTransport, "cache", 1)
	// fetcher := tool.NewFetcher(logger, http.DefaultTransport, "cache", 1)
	// fetcher := tool.ReadFetcher(logger, "cache")

	ice, err := eu_ec_ice.Fetch(context.Background(), fetcher)
	if err != nil {
		logger.Error(err.Error())
	}

	logger.Info("iceSlice.len", "len", len(ice))
}
