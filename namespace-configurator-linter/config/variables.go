package config

import (
	"log"
	"os"
)

var (
	ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLog  = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
)

type Folder struct {
	Products []struct {
		ProductID        string   `yaml:"product-id,omitempty"`
		NamespacePattern []string `yaml:"namespace_pattern,omitempty"`
		Exclusions       []string `yaml:"exclusions,omitempty"`
	} `yaml:"products,omitempty"`
}
