package tool

import (
	"flag"
	"log/slog"
	"os"
	"path/filepath"
	"sniffle/myhandler"
	"sniffle/tool/fetch"
	"sniffle/tool/writefs"
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
		config.LogHandler = slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: level.Level})
	} else {
		config.LogHandler = myhandler.New(logFile, level.Level)
	}

	config.Writefile = writefs.Os(*out)

	if *cacheRemove {
		if err := os.RemoveAll(filepath.Join(*cache, "http")); err != nil {
			slog.New(config.LogHandler).Error("rm-cache", "err", err)
			os.Exit(1)
		}
		if err := os.RemoveAll(filepath.Join(*cache, "https")); err != nil {
			slog.New(config.LogHandler).Error("rm-cache", "err", err)
			os.Exit(1)
		}
	}

	if delay == nil {
		delay = make(map[string]time.Duration, 1)
	}
	if _, ok := delay[""]; !ok {
		delay[""] = time.Millisecond * 100
	}

	config.Fetcher = []fetch.Fetcher{
		fetch.Cache(writefs.Os(*cache)),
		fetch.Net(nil, writefs.Os(*cache), delay),
	}

	config.LongTasksCache = writefs.Os(filepath.Join(*cache, "longtask"))
	config.LongTasksMap = make(map[string]func(*Tool, []byte) ([]byte, error))

	return config
}

// Return the cache file system from CLI value.
// It need to call CLI before.
func CLICache() writefs.CompleteFS {
	return writefs.Os(flag.CommandLine.Lookup("cache").Value.String())
}

type levelValue struct {
	slog.Level
}

func (l *levelValue) Set(s string) error { return l.Level.UnmarshalText([]byte(s)) }
