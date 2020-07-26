package validator

import (
	"regexp"
)

type CollectorType uint

const (
	Directory CollectorType = 1 << iota
)

type Config struct {
	Paths        []string
	ExcludeFiles []string
	Layers       []Layer
	Ruleset      Ruleset
}

type Layer struct {
	Name       LayerName
	Collectors []Collector
}

type Collector struct {
	Type  CollectorType
	Regex *regexp.Regexp
}

type Ruleset map[LayerName][]LayerName

func ParseConfig() Config {
	return Config{
		Paths:        []string{"./"},
		ExcludeFiles: []string{},
		Layers: []Layer{
			{
				Name: "Domain",
				Collectors: []Collector{
					{
						Type:  Directory,
						Regex: regexp.MustCompile("domain/.*"),
					},
				},
			},
			{
				Name: "Application",
				Collectors: []Collector{
					{
						Type:  Directory,
						Regex: regexp.MustCompile("app/.*"),
					},
				},
			},
			{
				Name: "Infrastructure",
				Collectors: []Collector{
					{
						Type:  Directory,
						Regex: regexp.MustCompile("infrastructure/.*"),
					},
				},
			},
		},
		Ruleset: Ruleset{
			"Domain":         []LayerName{},
			"Application":    []LayerName{"Domain"},
			"Infrastructure": []LayerName{"Application"},
		},
	}
}
