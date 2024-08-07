package main

import (
	"context"
	"log/slog"
	"os"
	"sniffle/api/eu_ec_ice"
	"sniffle/myhandler"
	"sniffle/tool"
	"sniffle/tool/fetch"
	"sniffle/tool/writefile"
	"time"
)

func main() {
	logger := slog.New(myhandler.New(os.Stdout, slog.LevelDebug))
	// logger := slog.New(myhandler.New(os.Stdout, slog.LevelInfo))

	// fetcher := tool.FallBackFetcher(logger, http.DefaultTransport, "cache", 1)
	// fetcher := tool.NewFetcher(logger, http.DefaultTransport, "cache", 1)
	// fetcher := tool.ReadFetcher(logger, "cache")

	fetcher := tool.New(&tool.Config{
		Logger:    logger,
		Writefile: writefile.Os("public"),
		Fetcher: []fetch.Fetcher{
			fetch.CacheOnly("cache"),
			fetch.Net(nil, "cache", 1, time.Millisecond*100),
		},
	})

	ice, err := eu_ec_ice.Fetch(context.Background(), fetcher)
	if err != nil {
		logger.Error(err.Error())
	}

	logger.Info("iceSlice.len", "len", len(ice))

	fetcher.WriteFile("eu/ec/ice/file.txt", []byte("Hello World"))
}
