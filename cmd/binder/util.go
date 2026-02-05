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
	metadataPath := filepath.Join(outdir, "metadata.yaml")
	err = os.WriteFile(metadataPath, contents, 0644)
	return err
}
