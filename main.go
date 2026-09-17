package main

import (
	"time"

	"github.com/dta0one/pokedexcli/internal/pokecache"
)

func main() {
	cfg := &config{
		cache:   pokecache.NewCache(5 * time.Minute),
		pokedex: make(map[string]Pokemon),
	}
	startRepl(cfg)
}
