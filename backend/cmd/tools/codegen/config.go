package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type FilePath string

func (f FilePath) String() string {
	return string(f)
}

type Config struct {
	Typescript struct {
		DataContracts []string `json:"data-contracts"`
		SwaggerFile   string   `json:"swaggerfile"`
		OutputRoutes  string   `json:"output-routes"`
	}
}

func (c *Config) Load(path string) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0o644)
	if err != nil {
		return err
	}

	err = json.NewDecoder(file).Decode(c)
	if err != nil {
		return err
	}

	// override the FilePaths to be relative to the config file
	dir, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		return err
	}

	for i := range c.Typescript.DataContracts {
		c.Typescript.DataContracts[i] = filepath.Join(dir, c.Typescript.DataContracts[i])
	}

	c.Typescript.SwaggerFile = filepath.Join(dir, c.Typescript.SwaggerFile)
	c.Typescript.OutputRoutes = filepath.Join(dir, c.Typescript.OutputRoutes)

	return nil
}

func (c *Config) Dump() string {
	b, _ := json.MarshalIndent(c, "", "  ")
	return string(b)
}

type OutputOnly struct {
	Output FilePath `json:"output"`
}
