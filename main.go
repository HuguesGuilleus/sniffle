package main

import (
	"context"
	"log/slog"
	"os"
	"sniffle/api/eu_ec_ice"
	"sniffle/front"
	"sniffle/myhandler"
	"sniffle/tool"
	"sniffle/tool/fetch"
	"sniffle/tool/language"
	"sniffle/tool/writefile"
	"time"
)

func main() {
	logger := slog.New(myhandler.New(os.Stdout, slog.LevelDebug))
	// logger := slog.New(myhandler.New(os.Stdout, slog.LevelInfo))

	// fetcher := tool.FallBackFetcher(logger, http.DefaultTransport, "cache", 1)
	// fetcher := tool.NewFetcher(logger, http.DefaultTransport, "cache", 1)
	// fetcher := tool.ReadFetcher(logger, "cache")

	t := tool.New(&tool.Config{
		Logger:    logger,
		HostURL:   "https://sniffle.eu/",
		Languages: []language.Language{language.English, language.French},
		Writefile: writefile.Os("public"),
		Fetcher: []fetch.Fetcher{
			fetch.CacheOnly("cache"),
			fetch.Net(nil, "cache", 1, time.Millisecond*100),
		},
	})

	front.WriteAssets(t)

	eu_ec_ice.Do(context.Background(), t)
	// ice, err := eu_ec_ice.Fetch(context.Background(), t)
	// if err != nil {
	// 	logger.Error(err.Error())
	// }

	// logger.Info("iceSlice.len", "len", len(ice))
}
