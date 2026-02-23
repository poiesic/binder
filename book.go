package binder

import (
	"errors"
	"fmt"
	"iter"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/divan/num2words"
	"gopkg.in/yaml.v3"
)

type FrontMatter struct {
	Title               string `yaml:"title"`
	ShortTitle          string `yaml:"short_title"`
	Author              string `yaml:"author"`
	AuthorLastName      string `yaml:"author_lastname"`
	ContactName         string `yaml:"contact_name"`
	ContactAddress      string `yaml:"contact_address"`
	ContactCityStateZip string `yaml:"contact_city_state_zip"`
	ContactPhone        string `yaml:"contact_phone"`
	ContactEmail        string `yaml:"contact_email"`
}

type Chapter struct {
	Name      string   `yaml:"name,omitempty"`
	Interlude bool     `yaml:"interlude,omitempty"`
	Subdir    string   `yaml:"subdir,omitempty"`
	Scenes    []string `yaml:"scenes"`
}

type Book struct {
	BaseDir  string `yaml:"base_dir"`
	Chapters []Chapter
}

type IteratedChapter struct {
	Filename string
	Heading  string
	Scenes   []string
}

func (ic IteratedChapter) Validate() error {
	for _, scene := range ic.Scenes {
		if _, err := os.Stat(scene); err != nil {
			return err
		}
	}
	return nil
}

func (ic IteratedChapter) HeadingToFilename() string {
	if ic.Heading == "" {
		return "interlude"
	}
	return strings.ToLower(strings.ReplaceAll(ic.Heading, " ", "-"))
}

func (b *Book) GetChapters() iter.Seq[IteratedChapter] {
	caser := cases.Title(language.English)
	return func(yield func(IteratedChapter) bool) {
		cn := 1
		for _, chapter := range b.Chapters {
			ic := &IteratedChapter{
				Scenes: make([]string, len(chapter.Scenes)),
			}
			var chapterBaseDir string
			if chapter.Subdir != "" {
				chapterBaseDir = filepath.Join(b.BaseDir, chapter.Subdir)
			} else {
				chapterBaseDir = b.BaseDir
			}
			if chapter.Name != "" {
				ic.Heading = caser.String(chapter.Name)
			} else {
				if !chapter.Interlude {
					ic.Heading = fmt.Sprintf("Chapter %s", caser.String(num2words.Convert(cn)))
					cn += 1
				}
			}
			for i, s := range chapter.Scenes {
				ic.Scenes[i] = fmt.Sprintf("%s.md", filepath.Join(chapterBaseDir, s))
			}
			if !yield(*ic) {
				return
			}
		}
	}
}

type BookSpec struct {
	Book Book `yaml:"book"`
}

func LoadBook(fileName string) (*FrontMatter, *Book, error) {
	fd, err := os.Open(fileName)
	if err != nil {
		return nil, nil, err
	}
	fm := &FrontMatter{}
	bs := &BookSpec{}
	decoder := yaml.NewDecoder(fd)
	if err := decoder.Decode(fm); err != nil {
		return nil, nil, err
	}
	if err := decoder.Decode(bs); err != nil {
		return nil, nil, err
	}
	inputDir := filepath.Dir(fileName)
	book := &(bs.Book)
	relativeDir := filepath.Join(inputDir, book.BaseDir)
	info, err := os.Stat(relativeDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fm, book, nil
		}
		return nil, nil, err
	}
	if info.IsDir() {
		book.BaseDir = relativeDir
	}
	return fm, book, nil
}
