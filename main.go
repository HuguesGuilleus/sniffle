package main

import (
	"flag"
	"sniffle/service"
	"sniffle/tool"
	"sniffle/tool/language"
	"time"
)

func main() {
	host := flag.String("host", "https://sniffle.eu/", "The host absolute URL")

	config := tool.CLI(map[string]time.Duration{"": time.Millisecond * 100})

	config.HostURL = *host
	config.Languages = []language.Language{language.English, language.French}
	config.LongTasksMap = service.LongTask

	tool.Run(config, service.List...)
}
