package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestFrontMatter_Unmarshal(t *testing.T) {
	yamlData := `
title: Benthamville
short_title: Benthamville
author: Kevin Smith
author_lastname: Smith
contact_name: Kevin Smith
contact_address: 6513 Rainbow Court
contact_city_state_zip: Raleigh, NC 27612
contact_phone: (919) 345-4521
contact_email: kevin@poiesic.com
`

	var fm FrontMatter
	err := yaml.Unmarshal([]byte(yamlData), &fm)
	require.NoError(t, err)

	assert.Equal(t, "Benthamville", fm.Title)
	assert.Equal(t, "Benthamville", fm.ShortTitle)
	assert.Equal(t, "Kevin Smith", fm.Author)
	assert.Equal(t, "Smith", fm.AuthorLastName)
	assert.Equal(t, "Kevin Smith", fm.ContactName)
	assert.Equal(t, "6513 Rainbow Court", fm.ContactAddress)
	assert.Equal(t, "Raleigh, NC 27612", fm.ContactCityStateZip)
	assert.Equal(t, "(919) 345-4521", fm.ContactPhone)
	assert.Equal(t, "kevin@poiesic.com", fm.ContactEmail)
}

func TestFrontMatter_UnmarshalFromActualYAML(t *testing.T) {
	// This test uses the format from front_matter.yaml
	yamlData := `
title: Benthamville
short_title: Benthamville
author: Kevin Smith
author_lastname: Smith
contact_name: Kevin Smith
contact_address: 6513 Rainbow Court
contact_city_state_zip: Raleigh, NC 27612
contact_phone: (919) 345-4521
contact_email: kevin@poiesic.com
`

	var fm FrontMatter
	err := yaml.Unmarshal([]byte(yamlData), &fm)
	require.NoError(t, err)

	assert.Equal(t, "Benthamville", fm.Title)
	assert.Equal(t, "Benthamville", fm.ShortTitle)
	assert.Equal(t, "Kevin Smith", fm.Author)
	// This will fail because the struct expects "author_last_name" but YAML has "author_lastname"
	assert.Equal(t, "Smith", fm.AuthorLastName, "AuthorLastName should be populated from author_lastname")
	assert.Equal(t, "Kevin Smith", fm.ContactName)
	assert.Equal(t, "6513 Rainbow Court", fm.ContactAddress)
	assert.Equal(t, "Raleigh, NC 27612", fm.ContactCityStateZip)
	assert.Equal(t, "(919) 345-4521", fm.ContactPhone)
	assert.Equal(t, "kevin@poiesic.com", fm.ContactEmail)
}

func TestFrontMatter_EmptyFields(t *testing.T) {
	yamlData := `title: Test Book`

	var fm FrontMatter
	err := yaml.Unmarshal([]byte(yamlData), &fm)
	require.NoError(t, err)

	assert.Equal(t, "Test Book", fm.Title)
	assert.Empty(t, fm.ShortTitle)
	assert.Empty(t, fm.Author)
	assert.Empty(t, fm.AuthorLastName)
	assert.Empty(t, fm.ContactName)
	assert.Empty(t, fm.ContactAddress)
	assert.Empty(t, fm.ContactCityStateZip)
	assert.Empty(t, fm.ContactPhone)
	assert.Empty(t, fm.ContactEmail)
}

func TestFrontMatter_PartialFields(t *testing.T) {
	yamlData := `
title: My Novel
author: Jane Doe
contact_email: jane@example.com
`

	var fm FrontMatter
	err := yaml.Unmarshal([]byte(yamlData), &fm)
	require.NoError(t, err)

	assert.Equal(t, "My Novel", fm.Title)
	assert.Empty(t, fm.ShortTitle)
	assert.Equal(t, "Jane Doe", fm.Author)
	assert.Empty(t, fm.AuthorLastName)
	assert.Empty(t, fm.ContactName)
	assert.Empty(t, fm.ContactAddress)
	assert.Empty(t, fm.ContactCityStateZip)
	assert.Empty(t, fm.ContactPhone)
	assert.Equal(t, "jane@example.com", fm.ContactEmail)
}

func TestBook_Unmarshal(t *testing.T) {
	// Test basic Book unmarshaling with proper YAML structure
	yamlData := `
base_dir: "manuscript"
chapters:
  - scenes:
      - "chapter1_scene1"
      - "chapter1_scene2"
  - scenes:
      - "chapter2_scene1"
`

	var book Book
	err := yaml.Unmarshal([]byte(yamlData), &book)
	require.NoError(t, err)

	assert.Equal(t, "manuscript", book.BaseDir)
	require.Len(t, book.Chapters, 2)
	assert.Equal(t, []string{"chapter1_scene1", "chapter1_scene2"}, book.Chapters[0].Scenes)
	assert.Equal(t, []string{"chapter2_scene1"}, book.Chapters[1].Scenes)
}

func TestBook_UnmarshalFromActualYAML(t *testing.T) {
	// This mimics the actual book.yaml format with top-level "book:" key
	// Note: The actual YAML has duplicate "chapter:" keys which is invalid YAML
	// Standard YAML parsers will only keep the last duplicate key
	yamlData := `
book:
  base_dir: "testdata/manuscript"
  chapters:
    - scenes:
        - "foo"
        - "baz"
    - scenes:
        - "bar"
        - "quux"
`

	// Need a wrapper struct for the top-level "book:" key
	type BookFile struct {
		Book Book `yaml:"book"`
	}

	var bf BookFile
	err := yaml.Unmarshal([]byte(yamlData), &bf)
	require.NoError(t, err)

	assert.Equal(t, "testdata/manuscript", bf.Book.BaseDir)
	require.Len(t, bf.Book.Chapters, 2)
	assert.Equal(t, []string{"foo", "baz"}, bf.Book.Chapters[0].Scenes)
	assert.Equal(t, []string{"bar", "quux"}, bf.Book.Chapters[1].Scenes)
}

func TestChapter_Unmarshal(t *testing.T) {
	yamlData := `
scenes:
  - "opening"
  - "conflict"
  - "resolution"
`

	var ch Chapter
	err := yaml.Unmarshal([]byte(yamlData), &ch)
	require.NoError(t, err)

	require.Len(t, ch.Scenes, 3)
	assert.Equal(t, "opening", ch.Scenes[0])
	assert.Equal(t, "conflict", ch.Scenes[1])
	assert.Equal(t, "resolution", ch.Scenes[2])
}

func TestChapter_EmptyScenes(t *testing.T) {
	yamlData := `scenes: []`

	var ch Chapter
	err := yaml.Unmarshal([]byte(yamlData), &ch)
	require.NoError(t, err)

	assert.Empty(t, ch.Scenes)
}

func TestBook_EmptyChapters(t *testing.T) {
	yamlData := `
base_dir: "manuscript"
chapters: []
`

	var book Book
	err := yaml.Unmarshal([]byte(yamlData), &book)
	require.NoError(t, err)

	assert.Equal(t, "manuscript", book.BaseDir)
	assert.Empty(t, book.Chapters)
}

func TestBook_CurrentYAMLFormat_DuplicateKeysAreInvalid(t *testing.T) {
	// This test documents that the current book.yaml format is INVALID YAML.
	// Duplicate "chapter:" keys cause a parse error in yaml.v3
	yamlData := `
book:
  base_dir: "testdata/manuscript"
  chapter:
    - "foo"
    - "baz"
  chapter:
    - "bar"
    - "quux"
`

	type BookFile struct {
		Book Book `yaml:"book"`
	}

	var bf BookFile
	err := yaml.Unmarshal([]byte(yamlData), &bf)
	require.Error(t, err, "Duplicate YAML keys should cause a parse error")
	assert.Contains(t, err.Error(), "already defined")
}

func TestBook_CorrectYAMLFormat(t *testing.T) {
	// This test shows the CORRECT YAML format that works with the Book struct
	yamlData := `
book:
  base_dir: "testdata/manuscript"
  chapters:
    - scenes:
        - "foo"
        - "baz"
    - scenes:
        - "bar"
        - "quux"
`

	type BookFile struct {
		Book Book `yaml:"book"`
	}

	var bf BookFile
	err := yaml.Unmarshal([]byte(yamlData), &bf)
	require.NoError(t, err)

	assert.Equal(t, "testdata/manuscript", bf.Book.BaseDir)
	require.Len(t, bf.Book.Chapters, 2)
	assert.Equal(t, []string{"foo", "baz"}, bf.Book.Chapters[0].Scenes)
	assert.Equal(t, []string{"bar", "quux"}, bf.Book.Chapters[1].Scenes)
}

func TestCombinedYAML_MultiDocument(t *testing.T) {
	// Combined YAML uses multi-document format with --- separators
	// First document: FrontMatter
	// Second document: Book (wrapped in "book:" key)
	yamlData := `---
title: Benthamville
short_title: Benthamville
author: Kevin Smith
author_lastname: Smith
contact_name: Kevin Smith
contact_address: 6513 Rainbow Court
contact_city_state_zip: Raleigh, NC 27612
contact_phone: (919) 345-4521
contact_email: kevin@poiesic.com
---
book:
  base_dir: "testdata/manuscript"
  chapters:
    - scenes:
        - "foo"
        - "baz"
    - scenes:
        - "bar"
        - "quux"
`

	decoder := yaml.NewDecoder(bytes.NewReader([]byte(yamlData)))

	// Decode first document: FrontMatter
	var fm FrontMatter
	err := decoder.Decode(&fm)
	require.NoError(t, err)

	assert.Equal(t, "Benthamville", fm.Title)
	assert.Equal(t, "Benthamville", fm.ShortTitle)
	assert.Equal(t, "Kevin Smith", fm.Author)
	assert.Equal(t, "Smith", fm.AuthorLastName)
	assert.Equal(t, "Kevin Smith", fm.ContactName)
	assert.Equal(t, "6513 Rainbow Court", fm.ContactAddress)
	assert.Equal(t, "Raleigh, NC 27612", fm.ContactCityStateZip)
	assert.Equal(t, "(919) 345-4521", fm.ContactPhone)
	assert.Equal(t, "kevin@poiesic.com", fm.ContactEmail)

	// Decode second document: Book
	type BookFile struct {
		Book Book `yaml:"book"`
	}
	var bf BookFile
	err = decoder.Decode(&bf)
	require.NoError(t, err)

	assert.Equal(t, "testdata/manuscript", bf.Book.BaseDir)
	require.Len(t, bf.Book.Chapters, 2)
	assert.Equal(t, []string{"foo", "baz"}, bf.Book.Chapters[0].Scenes)
	assert.Equal(t, []string{"bar", "quux"}, bf.Book.Chapters[1].Scenes)
}

func TestCombinedYAML_EmptyDocuments(t *testing.T) {
	// Test handling when one document is minimal
	yamlData := `---
title: Minimal Book
---
book:
  base_dir: "src"
  chapters: []
`

	decoder := yaml.NewDecoder(bytes.NewReader([]byte(yamlData)))

	var fm FrontMatter
	err := decoder.Decode(&fm)
	require.NoError(t, err)
	assert.Equal(t, "Minimal Book", fm.Title)
	assert.Empty(t, fm.Author)

	var bf BookSpec
	err = decoder.Decode(&bf)
	require.NoError(t, err)
	assert.Equal(t, "src", bf.Book.BaseDir)
	assert.Empty(t, bf.Book.Chapters)
}

// LoadBook tests

func TestLoadBook_Success(t *testing.T) {
	fm, book, err := LoadBook("testdata/valid_book.yaml")
	require.NoError(t, err)
	require.NotNil(t, fm)
	require.NotNil(t, book)

	// Verify FrontMatter
	assert.Equal(t, "Test Book", fm.Title)
	assert.Equal(t, "Test", fm.ShortTitle)
	assert.Equal(t, "Test Author", fm.Author)
	assert.Equal(t, "Author", fm.AuthorLastName)
	assert.Equal(t, "Test Contact", fm.ContactName)
	assert.Equal(t, "123 Test St", fm.ContactAddress)
	assert.Equal(t, "Test City, TS 12345", fm.ContactCityStateZip)
	assert.Equal(t, "(555) 123-4567", fm.ContactPhone)
	assert.Equal(t, "test@example.com", fm.ContactEmail)

	// Verify Book
	assert.Equal(t, "testdata/manuscript", book.BaseDir)
	require.Len(t, book.Chapters, 2)
	assert.Equal(t, []string{"foo", "baz"}, book.Chapters[0].Scenes)
	assert.Equal(t, []string{"bar", "quux"}, book.Chapters[1].Scenes)
}

func TestLoadBook_FileNotFound(t *testing.T) {
	fm, book, err := LoadBook("testdata/nonexistent.yaml")
	require.Error(t, err)
	assert.Nil(t, fm)
	assert.Nil(t, book)
	assert.Contains(t, err.Error(), "no such file")
}

func TestLoadBook_InvalidFrontMatter(t *testing.T) {
	fm, book, err := LoadBook("testdata/invalid_frontmatter.yaml")
	require.Error(t, err)
	assert.Nil(t, fm)
	assert.Nil(t, book)
}

func TestLoadBook_InvalidBook(t *testing.T) {
	fm, book, err := LoadBook("testdata/invalid_book.yaml")
	require.Error(t, err)
	assert.Nil(t, fm)
	assert.Nil(t, book)
}

// Book.GetChapters tests

func TestBook_GetChapters_AutoGeneratedNames(t *testing.T) {
	book := &Book{
		BaseDir: "manuscript",
		Chapters: []Chapter{
			{Scenes: []string{"scene1", "scene2"}},
			{Scenes: []string{"scene3"}},
			{Scenes: []string{"scene4", "scene5", "scene6"}},
		},
	}

	var chapters []IteratedChapter
	for ch := range book.GetChapters() {
		chapters = append(chapters, ch)
	}

	require.Len(t, chapters, 3)

	// Verify auto-generated chapter names use num2words
	assert.Equal(t, "Chapter one", chapters[0].Heading)
	assert.Equal(t, "Chapter two", chapters[1].Heading)
	assert.Equal(t, "Chapter three", chapters[2].Heading)

	// Verify scene paths
	assert.Equal(t, []string{"manuscript/scene1.md", "manuscript/scene2.md"}, chapters[0].Scenes)
	assert.Equal(t, []string{"manuscript/scene3.md"}, chapters[1].Scenes)
	assert.Equal(t, []string{"manuscript/scene4.md", "manuscript/scene5.md", "manuscript/scene6.md"}, chapters[2].Scenes)
}

func TestBook_GetChapters_CustomNames(t *testing.T) {
	book := &Book{
		BaseDir: "manuscript",
		Chapters: []Chapter{
			{Name: "Prologue", Scenes: []string{"intro"}},
			{Scenes: []string{"middle"}},
			{Name: "Epilogue", Scenes: []string{"outro"}},
		},
	}

	var chapters []IteratedChapter
	for ch := range book.GetChapters() {
		chapters = append(chapters, ch)
	}

	require.Len(t, chapters, 3)

	// Custom names should be used when provided
	assert.Equal(t, "Prologue", chapters[0].Heading)
	// Unnamed chapters get auto-generated names, but the counter only increments for unnamed
	assert.Equal(t, "Chapter one", chapters[1].Heading)
	assert.Equal(t, "Epilogue", chapters[2].Heading)
}

func TestBook_GetChapters_WithSubdirs(t *testing.T) {
	book := &Book{
		BaseDir: "base",
		Chapters: []Chapter{
			{Subdir: "part1", Scenes: []string{"scene1", "scene2"}},
			{Subdir: "part2", Scenes: []string{"scene3"}},
			{Scenes: []string{"scene4"}}, // No subdir, uses base_dir directly
		},
	}

	var chapters []IteratedChapter
	for ch := range book.GetChapters() {
		chapters = append(chapters, ch)
	}

	require.Len(t, chapters, 3)

	// Verify paths include subdirs when specified
	assert.Equal(t, []string{"base/part1/scene1.md", "base/part1/scene2.md"}, chapters[0].Scenes)
	assert.Equal(t, []string{"base/part2/scene3.md"}, chapters[1].Scenes)
	assert.Equal(t, []string{"base/scene4.md"}, chapters[2].Scenes)
}

func TestBook_GetChapters_EmptyChapters(t *testing.T) {
	book := &Book{
		BaseDir:  "manuscript",
		Chapters: []Chapter{},
	}

	var chapters []IteratedChapter
	for ch := range book.GetChapters() {
		chapters = append(chapters, ch)
	}

	assert.Empty(t, chapters)
}

func TestBook_GetChapters_EmptyScenes(t *testing.T) {
	book := &Book{
		BaseDir: "manuscript",
		Chapters: []Chapter{
			{Scenes: []string{}},
		},
	}

	var chapters []IteratedChapter
	for ch := range book.GetChapters() {
		chapters = append(chapters, ch)
	}

	require.Len(t, chapters, 1)
	assert.Equal(t, "Chapter one", chapters[0].Heading)
	assert.Empty(t, chapters[0].Scenes)
}

func TestBook_GetChapters_EarlyTermination(t *testing.T) {
	book := &Book{
		BaseDir: "manuscript",
		Chapters: []Chapter{
			{Scenes: []string{"scene1"}},
			{Scenes: []string{"scene2"}},
			{Scenes: []string{"scene3"}},
			{Scenes: []string{"scene4"}},
		},
	}

	// Only collect first 2 chapters to test early termination
	var chapters []IteratedChapter
	count := 0
	for ch := range book.GetChapters() {
		chapters = append(chapters, ch)
		count++
		if count >= 2 {
			break
		}
	}

	require.Len(t, chapters, 2)
	assert.Equal(t, "Chapter one", chapters[0].Heading)
	assert.Equal(t, "Chapter two", chapters[1].Heading)
}

func TestBook_GetChapters_LoadedFromFile(t *testing.T) {
	fm, book, err := LoadBook("testdata/book_with_named_chapters.yaml")
	require.NoError(t, err)
	require.NotNil(t, fm)
	require.NotNil(t, book)

	var chapters []IteratedChapter
	for ch := range book.GetChapters() {
		chapters = append(chapters, ch)
	}

	require.Len(t, chapters, 3)
	assert.Equal(t, "Prologue", chapters[0].Heading)
	assert.Equal(t, "Chapter one", chapters[1].Heading)
	assert.Equal(t, "Epilogue", chapters[2].Heading)

	assert.Equal(t, []string{"testdata/manuscript/foo.md"}, chapters[0].Scenes)
	assert.Equal(t, []string{"testdata/manuscript/bar.md", "testdata/manuscript/baz.md"}, chapters[1].Scenes)
	assert.Equal(t, []string{"testdata/manuscript/quux.md"}, chapters[2].Scenes)
}

// IteratedChapter.Validate tests

func TestIteratedChapter_Validate_AllScenesExist(t *testing.T) {
	ic := IteratedChapter{
		Heading: "Chapter one",
		Scenes: []string{
			"testdata/manuscript/foo.md",
			"testdata/manuscript/bar.md",
		},
	}

	err := ic.Validate()
	assert.NoError(t, err)
}

func TestIteratedChapter_Validate_EmptyScenes(t *testing.T) {
	ic := IteratedChapter{
		Heading: "Empty Chapter",
		Scenes:  []string{},
	}

	err := ic.Validate()
	assert.NoError(t, err)
}

func TestIteratedChapter_Validate_SceneNotFound(t *testing.T) {
	ic := IteratedChapter{
		Heading: "Chapter one",
		Scenes: []string{
			"testdata/manuscript/foo.md",
			"testdata/manuscript/nonexistent.md",
		},
	}

	err := ic.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no such file")
}

func TestIteratedChapter_Validate_FirstSceneNotFound(t *testing.T) {
	ic := IteratedChapter{
		Heading: "Chapter one",
		Scenes: []string{
			"testdata/manuscript/nonexistent.md",
			"testdata/manuscript/foo.md",
		},
	}

	err := ic.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no such file")
}

func TestIteratedChapter_Validate_AllScenesNotFound(t *testing.T) {
	ic := IteratedChapter{
		Heading: "Chapter one",
		Scenes: []string{
			"testdata/manuscript/missing1.md",
			"testdata/manuscript/missing2.md",
		},
	}

	err := ic.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no such file")
}

func TestIteratedChapter_Validate_IntegrationWithLoadBook(t *testing.T) {
	_, book, err := LoadBook("testdata/valid_book.yaml")
	require.NoError(t, err)

	for ch := range book.GetChapters() {
		err := ch.Validate()
		assert.NoError(t, err, "chapter %q should have all valid scenes", ch.Heading)
	}
}
