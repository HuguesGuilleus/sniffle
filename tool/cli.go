package tool

import (
	"flag"
	"log/slog"
	"os"
	"path/filepath"
	"sniffle/myhandler"
	"sniffle/tool/fetch"
	"sniffle/tool/writefile"
	"time"
)

// Add some flags, call [flag.Parse], and use result in the config.
func CLI(delay map[string]time.Duration) *Config {
	config := new(Config)

	level := levelValue{Level: NoticeLevel}
	flag.Var(&level, "log", "The log level: DEBUG|INFO|WARN|ERROR (case insensitive, suport int offset)")
	logOut := flag.String("logout", "-", "The output file to append logs")
	logJson := flag.Bool("logjson", false, "Use json handler or text")
	out := flag.String("out", "public", "The output directory")
	cache := flag.String("cache", "cache", "The cache directory")
	cacheRemove := flag.Bool("cache-rm", false, "Remove HTTP(S) in the cache directory")
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
	if *logJson {
		config.Logger = slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: level.Level}))
	} else {
		config.Logger = slog.New(myhandler.New(logFile, level.Level))
	}

	config.Writefile = writefile.Os(*out)

	if *cacheRemove {
		if err := os.RemoveAll(filepath.Join(*cache, "http")); err != nil {
			config.Logger.Error("rm-cache", "err", err)
			os.Exit(1)
		}
		if err := os.RemoveAll(filepath.Join(*cache, "https")); err != nil {
			config.Logger.Error("rm-cache", "err", err)
			os.Exit(1)
		}
	}

	if delay == nil {
		delay = make(map[string]time.Duration, 1)
	}
	if _, ok := delay[""]; !ok {
		delay[""] = time.Millisecond * 100
	}

	config.LongTasksCache = writefile.Os(filepath.Join(*cache, "longtask"))
	config.Fetcher = []fetch.Fetcher{
		fetch.Cache(*cache),
		fetch.Net(nil, *cache, delay),
	}

	return config
}

type levelValue struct {
	slog.Level
}

func (l *levelValue) Set(s string) error { return l.Level.UnmarshalText([]byte(s)) }
