package validator

import (
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v3"
)

type CollectorType uint

const (
	Directory CollectorType = 1 << iota
)

type Config struct {
	baseDir      string
	Paths        []string `yaml:"paths"`
	ExcludeFiles []string `yaml:"exclude_files"`
	Layers       []Layer  `yaml:"layers"`
	Ruleset      Ruleset  `yaml:"ruleset"`
}

type Layer struct {
	Name       LayerName   `yaml:"name"`
	Collectors []Collector `yaml:"collectors"`
}

type Collector struct {
	Type  string `yaml:"type"`
	type_ CollectorType
	Regex string `yaml:"regex"`
	regex *regexp.Regexp
}

type Ruleset map[LayerName][]LayerName

func ParseConfig(configPaths []string) Config {
	config := Config{}

	for _, configPath := range configPaths {
		absPath, err := filepath.Abs(configPath)
		if err != nil {
			panic(err)
		}

		yamlFile, err := ioutil.ReadFile(absPath)
		if err != nil {
			panic(err)
		}

		tempConfig := Config{}
		err = yaml.Unmarshal(yamlFile, &tempConfig)
		if err != nil {
			panic(err)
		}

		err = mergo.Merge(&config, tempConfig)
		if err != nil {
			panic(err)
		}
	}

	absPath, _ := filepath.Abs(configPaths[0])
	absDir := filepath.Dir(absPath)
	config.baseDir = filepath.Dir(absDir)
	for i, p := range config.Paths {
		dir := filepath.Dir(p)
		config.Paths[i] = filepath.Join(absDir, dir)
	}

	for i, l := range config.Layers {
		for j, c := range l.Collectors {
			// TODO: Map type as well
			config.Layers[i].Collectors[j].regex = regexp.MustCompile(c.Regex)
		}
	}

	// return Config{
	// 	Paths:        []string{"./"},
	// 	ExcludeFiles: []string{},
	// 	Layers: []Layer{
	// 		{
	// 			Name: "Domain",
	// 			Collectors: []Collector{
	// 				{
	// 					Type:  Directory,
	// 					Regex: regexp.MustCompile("domain/.*"),
	// 				},
	// 			},
	// 		},
	// 		{
	// 			Name: "Application",
	// 			Collectors: []Collector{
	// 				{
	// 					Type:  Directory,
	// 					Regex: regexp.MustCompile("app/.*"),
	// 				},
	// 			},
	// 		},
	// 		{
	// 			Name: "Infrastructure",
	// 			Collectors: []Collector{
	// 				{
	// 					Type:  Directory,
	// 					Regex: regexp.MustCompile("infrastructure/.*"),
	// 				},
	// 			},
	// 		},
	// 	},
	// 	Ruleset: Ruleset{
	// 		"Domain":         []LayerName{},
	// 		"Application":    []LayerName{"Domain"},
	// 		"Infrastructure": []LayerName{"Application"},
	// 	},
	// }
	return config
}
