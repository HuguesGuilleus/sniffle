package testingtool

import (
	"context"
	"fmt"
)

type TestFetcher map[string][]byte

func (tf TestFetcher) FetchGET(_ context.Context, url string) ([]byte, error) {
	data, ok := tf[url]
	if !ok {
		return nil, fmt.Errorf("not found %q", url)
	}
	return data, nil
}

func (tf TestFetcher) Error(msg string, args ...any) {
	fmt.Printf("[%q] ", msg)
	fmt.Println(args...)
	panic(msg)
}
