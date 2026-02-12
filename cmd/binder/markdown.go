package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v3"
)

func writeMarkdownScenes(fd *os.File, sceneFiles []string) error {
	lastSceneIndex := len(sceneFiles) - 1
	for i, sceneFile := range sceneFiles {
		sceneText, err := os.ReadFile(sceneFile)
		if err != nil {
			return err
		}
		if _, err := fd.Write(sceneText); err != nil {
			return err
		}
		if i < lastSceneIndex {
			if _, err := fd.WriteString("\n\n***\n\n"); err != nil {
				return err
			}
		}
	}
	return nil
}

func sceneWordCount(path string) (string, error) {
	out, err := exec.Command("wc", "-w", path).Output()
	if err != nil {
		return "", err
	}
	// wc -w outputs "  123 path/to/file", grab just the count
	fields := strings.Fields(strings.TrimSpace(string(out)))
	if len(fields) == 0 {
		return "0", nil
	}
	return fields[0], nil
}

func markdown(ctx context.Context, cmd *cli.Command) error {
	outdir := cmd.String("outdir")
	wordcount := cmd.Bool("wordcount")
	_ = os.RemoveAll(outdir)
	if err := os.MkdirAll(outdir, 0755); err != nil {
		return err
	}
	frontMatter, book, err := LoadBook(cmd.String("input"))
	if err != nil {
		return err
	}
	cnum := 1
	for chapter := range book.GetChapters() {
		if err := chapter.Validate(); err != nil {
			return err
		}
		chapterOutPath := filepath.Join(outdir, fmt.Sprintf("%03d-%s.md", cnum, chapter.HeadingToFilename()))
		cnum += 1
		fd, err := os.OpenFile(chapterOutPath, os.O_CREATE|os.O_WRONLY|os.O_SYNC, 0644)
		if err != nil {
			return err
		}
		if chapter.Heading != "" {
			if _, err := fmt.Fprintf(fd, "# %s\n\n", chapter.Heading); err != nil {
				return err
			}
		}
		if wordcount {
			for _, scene := range chapter.Scenes {
				wc, err := sceneWordCount(scene)
				if err != nil {
					return err
				}
				fmt.Printf("%s: %s words\n", filepath.Base(scene), wc)
			}
		}
		if err := writeMarkdownScenes(fd, chapter.Scenes); err != nil {
			return err
		}
		if err := fd.Close(); err != nil {
			return err
		}
	}
	return writeMetadata(frontMatter, outdir)
}
