package main

import (
	"context"
	"flag"
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
	level := levelValue{Level: tool.NoticeLevel}
	flag.Var(&level, "log", "The log level: DEBUG|INFO|WARN|ERROR (case insensitive, suport int offset)")
	out := flag.String("out", "public", "The output directory")
	cache := flag.String("cache", "cache", "The cache directory")
	host := flag.String("host", "https://sniffle.eu/", "The host absolute URL")
	dev := flag.Bool("dev", false, "Dev mode")
	flag.Parse()

	config := tool.Config{
		Logger:    slog.New(myhandler.New(os.Stdout, level.Level)),
		HostURL:   *host,
		Languages: []language.Language{language.English, language.French},
		Writefile: writefile.Os(*out),
		Fetcher: []fetch.Fetcher{
			fetch.CacheOnly(*cache),
			fetch.Net(nil, *cache, 1, time.Millisecond*100),
		},
		Dev: *dev,
	}

	tool.Run(context.Background(), &config, service.List)
}

type levelValue struct {
	slog.Level
}

func (l *levelValue) Set(s string) error { return l.Level.UnmarshalText([]byte(s)) }
