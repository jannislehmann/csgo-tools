package demoparser

import "github.com/Cludch/csgo-tools/pkg/demo"

type UseCase interface {
	Parse(dir string, demoFile *demo.Demo) error
}
