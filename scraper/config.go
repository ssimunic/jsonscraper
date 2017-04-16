package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"regexp"
	"strings"
)

// Config holds information from config file.
type Config struct {
	URLs     []string  `json:"urls"`
	Targets  []*target `json:"targets"`
	Output   *output   `json:"output"`
	filePath string
	fileName string
}

type target struct {
	Selector   string `json:"selector"`
	Submatch   string `json:"submatch"`
	Tag        string `json:"tag"`
	Type       string `json:"type"`
	attrv      string
	submatchRe *regexp.Regexp
}

type output struct {
	Path string `json:"path"`
}

func newConfig(filePath string) (*Config, error) {
	c := new(Config)

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}

	c.filePath = filePath
	c.fileName = path.Base(filePath)

	for _, target := range c.Targets {
		if target.Submatch != "" {
			target.submatchRe, err = regexp.Compile(target.Submatch)
			if err != nil {
				return nil, err
			}
		}
		if !c.isValidTarget(target) {
			return nil, fmt.Errorf("missing target property")
		}
		if strings.HasPrefix(target.Type, "attr:") {
			target.attrv = strings.Split(target.Type, ":")[1]
		}
	}

	return c, nil
}

func (c *Config) outputPath() string {
	outputPath := c.Output.Path
	return strings.Replace(outputPath, "$FILENAME", c.fileName, -1)
}

func (c *Config) isValidTarget(t *target) bool {
	if t.Tag != "" && t.Type != "" && (t.Selector != "" || t.Submatch != "") {
		return true
	}
	return false
}
