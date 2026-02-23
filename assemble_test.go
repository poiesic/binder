package binder

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAssembleMarkdown_ValidBook(t *testing.T) {
	outdir := t.TempDir()

	config := AssemblyConfig{
		InputFile: "testdata/valid_book.yaml",
		OutputDir: outdir,
		WordCount: false,
	}
	fm, counts, err := AssembleMarkdown(config)
	require.NoError(t, err)
	require.NotNil(t, fm)
	assert.Empty(t, counts)

	assert.Equal(t, "Test Book", fm.Title)
	assert.Equal(t, "Test Author", fm.Author)

	// Should produce 4 chapter files + metadata.yaml
	entries, err := os.ReadDir(outdir)
	require.NoError(t, err)

	var mdFiles []string
	hasMetadata := false
	for _, e := range entries {
		if e.Name() == "metadata.yaml" {
			hasMetadata = true
		} else {
			mdFiles = append(mdFiles, e.Name())
		}
	}
	assert.True(t, hasMetadata, "metadata.yaml should be generated")
	assert.Len(t, mdFiles, 4)

	// Verify chapter file names
	assert.Equal(t, "001-interlude.md", mdFiles[0])
	assert.Equal(t, "002-chapter-one.md", mdFiles[1])
	assert.Equal(t, "003-interlude.md", mdFiles[2])
	assert.Equal(t, "004-chapter-two.md", mdFiles[3])

	// Verify chapter content has heading
	chapterOne, err := os.ReadFile(filepath.Join(outdir, "002-chapter-one.md"))
	require.NoError(t, err)
	assert.Contains(t, string(chapterOne), "# Chapter One")

	// Verify interlude has no heading (starts directly with content)
	interlude1, err := os.ReadFile(filepath.Join(outdir, "001-interlude.md"))
	require.NoError(t, err)
	assert.NotContains(t, string(interlude1), "#")
	assert.Contains(t, string(interlude1), "This is interlude 1.")
}

func TestAssembleMarkdown_WithWordCount(t *testing.T) {
	outdir := t.TempDir()

	config := AssemblyConfig{
		InputFile: "testdata/valid_book.yaml",
		OutputDir: outdir,
		WordCount: true,
	}
	fm, counts, err := AssembleMarkdown(config)
	require.NoError(t, err)
	require.NotNil(t, fm)
	require.NotEmpty(t, counts)

	// Should have a count for each scene (6 total: interlude1, foo, baz, interlude2, bar, quux)
	assert.Len(t, counts, 6)

	// Each test scene has a few words
	for _, wc := range counts {
		assert.Greater(t, wc.Count, 0, "scene %s should have words", wc.Scene)
	}
}

func TestAssembleMarkdown_MissingInput(t *testing.T) {
	outdir := t.TempDir()

	config := AssemblyConfig{
		InputFile: "testdata/nonexistent.yaml",
		OutputDir: outdir,
		WordCount: false,
	}
	fm, counts, err := AssembleMarkdown(config)
	require.Error(t, err)
	assert.Nil(t, fm)
	assert.Nil(t, counts)
}

func TestAssembleMarkdown_MetadataContent(t *testing.T) {
	outdir := t.TempDir()

	config := AssemblyConfig{
		InputFile: "testdata/valid_book.yaml",
		OutputDir: outdir,
		WordCount: false,
	}
	_, _, err := AssembleMarkdown(config)
	require.NoError(t, err)

	metadata, err := os.ReadFile(filepath.Join(outdir, "metadata.yaml"))
	require.NoError(t, err)

	content := string(metadata)
	assert.Contains(t, content, "---")
	assert.Contains(t, content, "title: Test Book")
	assert.Contains(t, content, "author: Test Author")
}

func TestAssembleMarkdown_SceneBreaks(t *testing.T) {
	outdir := t.TempDir()

	config := AssemblyConfig{
		InputFile: "testdata/valid_book.yaml",
		OutputDir: outdir,
		WordCount: false,
	}
	_, _, err := AssembleMarkdown(config)
	require.NoError(t, err)

	// Chapter with 2 scenes (foo, baz) should have a scene break
	chapterOne, err := os.ReadFile(filepath.Join(outdir, "002-chapter-one.md"))
	require.NoError(t, err)
	assert.Contains(t, string(chapterOne), "***")
}

func TestWriteMarkdownScenes(t *testing.T) {
	outFile := filepath.Join(t.TempDir(), "test.md")
	fd, err := os.OpenFile(outFile, os.O_CREATE|os.O_WRONLY, 0644)
	require.NoError(t, err)

	scenes := []string{
		"testdata/manuscript/foo.md",
		"testdata/manuscript/bar.md",
	}
	err = WriteMarkdownScenes(fd, scenes, false)
	require.NoError(t, err)
	fd.Close()

	content, err := os.ReadFile(outFile)
	require.NoError(t, err)

	text := string(content)
	assert.Contains(t, text, "This is foo.")
	assert.Contains(t, text, "***")
	assert.Contains(t, text, "This is bar.")
	assert.NotContains(t, text, "## foo")
}

func TestWriteMarkdownScenes_SingleScene(t *testing.T) {
	outFile := filepath.Join(t.TempDir(), "test.md")
	fd, err := os.OpenFile(outFile, os.O_CREATE|os.O_WRONLY, 0644)
	require.NoError(t, err)

	scenes := []string{"testdata/manuscript/foo.md"}
	err = WriteMarkdownScenes(fd, scenes, false)
	require.NoError(t, err)
	fd.Close()

	content, err := os.ReadFile(outFile)
	require.NoError(t, err)

	text := string(content)
	assert.Contains(t, text, "This is foo.")
	assert.NotContains(t, text, "***", "single scene should have no scene break")
}

func TestWriteMarkdownScenes_WithSceneHeadings(t *testing.T) {
	outFile := filepath.Join(t.TempDir(), "test.md")
	fd, err := os.OpenFile(outFile, os.O_CREATE|os.O_WRONLY, 0644)
	require.NoError(t, err)

	scenes := []string{
		"testdata/manuscript/foo.md",
		"testdata/manuscript/bar.md",
	}
	err = WriteMarkdownScenes(fd, scenes, true)
	require.NoError(t, err)
	fd.Close()

	content, err := os.ReadFile(outFile)
	require.NoError(t, err)

	text := string(content)
	assert.Contains(t, text, "## foo\n\nThis is foo.")
	assert.Contains(t, text, "## bar\n\nThis is bar.")
}

func TestAssembleMarkdown_WithSceneHeadings(t *testing.T) {
	outdir := t.TempDir()

	config := AssemblyConfig{
		InputFile:     "testdata/valid_book.yaml",
		OutputDir:     outdir,
		SceneHeadings: true,
	}
	_, _, err := AssembleMarkdown(config)
	require.NoError(t, err)

	chapterOne, err := os.ReadFile(filepath.Join(outdir, "002-chapter-one.md"))
	require.NoError(t, err)
	text := string(chapterOne)
	assert.Contains(t, text, "# Chapter One")
	assert.Contains(t, text, "## foo")
	assert.Contains(t, text, "## baz")
}

func TestWriteMetadata(t *testing.T) {
	outdir := t.TempDir()
	fm := &FrontMatter{
		Title:  "Test Title",
		Author: "Test Author",
	}

	err := WriteMetadata(fm, outdir)
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(outdir, "metadata.yaml"))
	require.NoError(t, err)

	text := string(content)
	assert.Contains(t, text, "---\n")
	assert.Contains(t, text, "title: Test Title")
	assert.Contains(t, text, "author: Test Author")
}

func TestSceneWordCount(t *testing.T) {
	// "This is foo." = 3 words
	count, err := SceneWordCount("testdata/manuscript/foo.md")
	require.NoError(t, err)
	assert.Equal(t, 3, count)
}

func TestSceneWordCount_FileNotFound(t *testing.T) {
	_, err := SceneWordCount("testdata/manuscript/nonexistent.md")
	require.Error(t, err)
}

func TestFormatWordCount(t *testing.T) {
	result := WordCountResult{Scene: "foo.md", Count: 42}
	assert.Equal(t, "foo.md: 42 words", FormatWordCount(result))
}

func TestOutputFiles(t *testing.T) {
	outdir := t.TempDir()

	config := AssemblyConfig{
		InputFile: "testdata/valid_book.yaml",
		OutputDir: outdir,
		WordCount: false,
	}
	_, _, err := AssembleMarkdown(config)
	require.NoError(t, err)

	files, err := OutputFiles(outdir)
	require.NoError(t, err)
	assert.Len(t, files, 4)
}
