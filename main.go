package main

import (
	"context"
	"log/slog"
	"os"
	"sniffle/myhandler"
	"sniffle/service"
	"sniffle/tool"
	"sniffle/tool/fetch"
	"sniffle/tool/language"
	"sniffle/tool/writefile"
	"time"
)

func main() {
	logger := slog.New(myhandler.New(os.Stdout, tool.NoticeLevel))
	// logger = slog.New(myhandler.New(os.Stdout, slog.LevelInfo))
	// logger = slog.New(myhandler.New(os.Stdout, slog.LevelDebug))

	config := tool.Config{
		Logger:    logger,
		HostURL:   "https://sniffle.eu/",
		Languages: []language.Language{language.English, language.French},
		Writefile: writefile.Os("public"),
		Fetcher: []fetch.Fetcher{
			fetch.CacheOnly("cache"),
			fetch.Net(nil, "cache", 1, time.Millisecond*100),
		},
	}

	tool.Run(context.Background(), &config, service.List)
}
