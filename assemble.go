package binder

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

// AssemblyConfig holds the parameters for assembling a manuscript.
type AssemblyConfig struct {
	InputFile string
	OutputDir string
	WordCount bool
}

// WordCountResult holds the word count for a single scene file.
type WordCountResult struct {
	Scene string
	Count int
}

// AssembleMarkdown assembles a book's scenes into per-chapter markdown files
// and writes a metadata.yaml for pandoc. Returns the parsed FrontMatter and
// any word count results (if config.WordCount is true).
func AssembleMarkdown(config AssemblyConfig) (*FrontMatter, []WordCountResult, error) {
	_ = os.RemoveAll(config.OutputDir)
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		return nil, nil, err
	}
	frontMatter, book, err := LoadBook(config.InputFile)
	if err != nil {
		return nil, nil, err
	}
	var counts []WordCountResult
	cnum := 1
	for chapter := range book.GetChapters() {
		if err := chapter.Validate(); err != nil {
			return nil, nil, err
		}
		chapterOutPath := filepath.Join(config.OutputDir, fmt.Sprintf("%03d-%s.md", cnum, chapter.HeadingToFilename()))
		cnum += 1
		fd, err := os.OpenFile(chapterOutPath, os.O_CREATE|os.O_WRONLY|os.O_SYNC, 0644)
		if err != nil {
			return nil, nil, err
		}
		if chapter.Heading != "" {
			if _, err := fmt.Fprintf(fd, "# %s\n\n", chapter.Heading); err != nil {
				fd.Close()
				return nil, nil, err
			}
		}
		if config.WordCount {
			for _, scene := range chapter.Scenes {
				wc, err := SceneWordCount(scene)
				if err != nil {
					fd.Close()
					return nil, nil, err
				}
				counts = append(counts, WordCountResult{
					Scene: filepath.Base(scene),
					Count: wc,
				})
			}
		}
		if err := WriteMarkdownScenes(fd, chapter.Scenes); err != nil {
			fd.Close()
			return nil, nil, err
		}
		if err := fd.Close(); err != nil {
			return nil, nil, err
		}
	}
	if err := WriteMetadata(frontMatter, config.OutputDir); err != nil {
		return nil, nil, err
	}
	return frontMatter, counts, nil
}

// WriteMarkdownScenes writes the contents of scene files to fd, separated
// by scene break markers.
func WriteMarkdownScenes(fd *os.File, sceneFiles []string) error {
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

// SceneWordCount counts the words in a file using pure Go.
func SceneWordCount(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	count := 0
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	return count, nil
}

// FormatWordCount formats a WordCountResult as a human-readable string.
func FormatWordCount(result WordCountResult) string {
	return fmt.Sprintf("%s: %d words", result.Scene, result.Count)
}

// OutputFiles returns a sorted list of the assembled chapter markdown files
// in the output directory.
func OutputFiles(outputDir string) ([]string, error) {
	pattern := filepath.Join(outputDir, "0*.md")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	return matches, nil
}
