package main

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func writeMetadata(fm *FrontMatter, outdir string) error {
	contents, err := yaml.Marshal(fm)
	if err != nil {
		return err
	}
	// Wrap in YAML front matter delimiters so Pandoc parses it as metadata
	wrapped := append([]byte("---\n"), contents...)
	wrapped = append(wrapped, []byte("---\n")...)
	metadataPath := filepath.Join(outdir, "metadata.yaml")
	err = os.WriteFile(metadataPath, wrapped, 0644)
	return err
}
