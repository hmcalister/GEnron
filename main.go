package main

import (
	"flag"
	"log/slog"

	"github.com/hmcalister/genron/config"
)

func main() {
	configFilePath := flag.String("configFilePath", "config.yaml", "Set the file path to the config file. Accepts JSON, YAML, TOML, and envfiles. See README for config specifications.")
	flag.Parse()

}
