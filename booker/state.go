package booker

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/drgo/core/ui"
	"github.com/drgo/mdson"
)
const mdson_ext = ".mdson"
//Algorithm:
// - read config.mdson 
// - parse templates from layout dir
// - parse mdson files from content dir
// -? do subdirs


type Config struct {
	// Dirs ignored in search and processing
	IgnoredDirs []string
	ContentPath  string // Relative or absolute location of article files and related content.
	TemplatePath string // Relative or absolute location of template files.

	BaseURL       string        // Absolute base URL (for permalinks; no trailing slash).
	BasePath      string        // Base URL path relative to server root (no trailing slash).
	Hostname      string        // Server host name, used for rendering ATOM feeds.
	// AnalyticsHTML template.HTML // Optional analytics HTML to insert at the beginning of <head>.

	HomeArticles int    // Articles to display on the home page.
	FeedArticles int    // Articles to include in Atom and JSON feeds.
	FeedTitle    string // The title of the Atom XML feed

	PlayEnabled     bool
	ServeLocalLinks bool // rewrite golang.org/{pkg,cmd} links to host-less, relative paths.
	Debug  ui.Debug
}

func newConfig() *Config{
	cfg := Config {
		Debug : ui.DebugAll,
	}	
	return &cfg
}
//Entry holds info on a book/site entry/page
type Entry struct {
	doc *mdson.Document
	Permalink string        // Canonical URL for this document.
	Path      string        // Path relative to server root (including base).
}

// State holds all info needed for a run
type State struct{
	// source filesystem
	sfs fs.FS
	cfg *Config
	ui.UI 
}

func newState(sfs fs.FS, cfg *Config) *State{
	return &State{
		sfs: sfs,
		cfg: cfg,
		UI : ui.NewUI(cfg.Debug),
	}
}
// Glob returns a list of all files 
// TODO: ignore dirs starting with .
// TODO: support filepath.glob 
func (s *State) Glob(pattern string) ([]string, error) {
	var paths []string

	err := fs.WalkDir(s.sfs, ".", func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		base := filepath.Base(path)
		for _, sp := range s.cfg.IgnoredDirs {
			// if the name of the folder has a prefix listed in SkipPaths
			// then we should skip the directory.
			// e.g. node_modules, testdata, _foo, .git
			if strings.HasPrefix(base, sp) {
				return filepath.SkipDir
			}
		}

		if filepath.Ext(path) != mdson_ext {
			return nil
		}

		paths = append(paths, path)

		return nil
	})

	return paths, err
}


func check(err error) {
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}



