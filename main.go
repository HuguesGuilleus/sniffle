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
	logOut := flag.String("logout", "-", "The output file to append logs")
	logJson := flag.Bool("logjson", false, "Use json handler ")
	out := flag.String("out", "public", "The output directory")
	cache := flag.String("cache", "cache", "The cache directory")
	host := flag.String("host", "https://sniffle.eu/", "The host absolute URL")
	flag.Parse()

	logFile := os.Stderr
	if *logOut != "-" {
		err := error(nil)
		logFile, err = os.OpenFile(*logOut, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o664)
		if err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}
	}

	logHandler := myhandler.New(logFile, level.Level)
	if *logJson {
		logHandler = slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: level.Level})
	}

	config := tool.Config{
		Logger:    slog.New(logHandler),
		HostURL:   *host,
		Languages: []language.Language{language.English, language.French},
		Writefile: writefile.Os(*out),
		Fetcher: []fetch.Fetcher{
			fetch.CacheOnly(*cache),
			fetch.Net(nil, *cache, 1, time.Millisecond*100),
		},
	}

	tool.Run(context.Background(), &config, service.List)
}

type levelValue struct {
	slog.Level
}

func (l *levelValue) Set(s string) error { return l.Level.UnmarshalText([]byte(s)) }
