package logstats

import (
	"regexp"
	"time"
)

type Configuration struct {
	Range time.Duration `yaml:"range"`
	Read  int           `yaml:"read"`

	Kinds map[string]KindConfig `yaml:"kinds"`

	Exclude ExcludeConfig `yaml:"exclude"`
}

func (c *Configuration) Compile() {
	for kind, cfg := range c.Kinds {
		cfg.Compile()
		c.Kinds[kind] = cfg
	}

	c.Exclude.Compile()
}

type ExcludeConfig struct {
	UriPatterns  []string `yaml:"uris"`
	FilePatterns []string `yaml:"files"`
	IPs          []string `yaml:"ips"`

	uriRegexes  []*regexp.Regexp
	fileRegexes []*regexp.Regexp
}

func (c *ExcludeConfig) Compile() {
	c.uriRegexes = make([]*regexp.Regexp, len(c.UriPatterns))

	for idx, pattern := range c.UriPatterns {
		c.uriRegexes[idx] = regexp.MustCompile(pattern)
	}

	c.fileRegexes = make([]*regexp.Regexp, len(c.FilePatterns))

	for idx, pattern := range c.FilePatterns {
		c.fileRegexes[idx] = regexp.MustCompile(pattern)
	}
}

func (c *ExcludeConfig) Matches(line *LogLine) bool {
	for _, ip := range c.IPs {
		if ip == line.IP {
			return true
		}
	}

	for _, regex := range c.uriRegexes {
		if regex.MatchString(line.Uri) {
			return true
		}
	}

	path := line.FilePath()

	for _, regex := range c.fileRegexes {
		if regex.MatchString(path) {
			return true
		}
	}

	return false
}

type KindConfig struct {
	UriPattern  string `yaml:"uri"`
	FilePattern string `yaml:"file"`

	uriRegex  *regexp.Regexp
	fileRegex *regexp.Regexp
}

func (c *KindConfig) Compile() {
	c.uriRegex = regexp.MustCompile(c.UriPattern)
	c.fileRegex = regexp.MustCompile(c.FilePattern)
}

func (c *KindConfig) Matches(line *LogLine) bool {
	if c.UriPattern != "" {
		if c.uriRegex.MatchString(line.Uri) {
			return true
		}
	}

	if c.FilePattern != "" {
		if c.fileRegex.MatchString(line.FilePath()) {
			return true
		}
	}

	return false
}
