package main

import "flag"

type flags struct {
	cfgPath     string
	storageType string
}

const defaultConfigPath = "local.yaml"
const defaultStorageType = "inMemory"

func parseFlags() *flags {
	var cfgPath = flag.String("path", defaultConfigPath, "path to config")
	var storageType = flag.String("storage", defaultStorageType, "path to config")

	flag.Parse()

	return &flags{
		cfgPath:     *cfgPath,
		storageType: *storageType,
	}
}
