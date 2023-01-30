package model

import (
	"github.com/flamego/i18n"
)

type Config struct {
	HttpAddr           string     `json:"httpAddr,omitempty" yaml:"httpAddr,omitempty" toml:"httpAddr,omitempty"`
	HttpPort           string     `json:"httpPort,omitempty" yaml:"httpPort,omitempty" toml:"httpPort,omitempty"`
	Description        string     `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	DocsDirectory      string     `json:"docsDirectory,omitempty" yaml:"docsDirectory,omitempty" toml:"docsDirectory,omitempty"`
	CustomDirectory    string     `json:"customDirectory,omitempty" yaml:"customDirectory,omitempty" toml:"customDirectory,omitempty"`
	HasLandingPage     bool       `json:"hasLandingPage,omitempty" yaml:"hasLandingPage,omitempty" toml:"hasLandingPage,omitempty"`
	HasNavBar          bool       `json:"hasNavBar,omitempty" yaml:"hasNavBar,omitempty" toml:"hasNavBar,omitempty"`
	DocsBasePath       string     `json:"docsBasePath,omitempty" yaml:"docsBasePath,omitempty" toml:"docsBasePath,omitempty"`
	EditPageLinkFormat string     `json:"editPageLinkFormat,omitempty" yaml:"editPageLinkFormat,omitempty" toml:"editPageLinkFormat,omitempty"`
	Languages          []Language `json:"languages,omitempty" yaml:"languages,omitempty" toml:"languages,omitempty"`
}

type Language struct {
	Name        string `json:"name" yaml:"name" toml:"name"`
	Description string `json:"description" yaml:"description" toml:"description"`
}

func (c *Config) I18nLanguages() []i18n.Language {
	if len(c.Languages) == 0 {
		return []i18n.Language{
			{Name: "en-US", Description: "English"},
		}
	}
	languages := make([]i18n.Language, len(c.Languages))
	for i, lan := range c.Languages {
		languages[i] = i18n.Language{
			Name:        lan.Name,
			Description: lan.Description,
		}
	}
	return languages
}
