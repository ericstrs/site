package render

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Content struct {
	Title     string
	Id        string
	UpdatedAt time.Time
}

// AllContent return all the content for the given content type
func AllContent(contentType string) ([]Content, error) {
	var items []Content
	baseDir := "docs/" + contentType

	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory's README
		if path == filepath.Join(baseDir, "README.md") {
			return nil
		}

		if !info.IsDir() && filepath.Ext(path) == ".md" {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			title := mdTitle(content)
			id := idFromPath(path)
			items = append(items, Content{
				Title:     title,
				Id:        id,
				UpdatedAt: info.ModTime(),
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Id < items[j].Id
	})

	return items, nil
}

// RecentContent returns n recent content for the given content type
func RecentContent(contentType string, n int) ([]Content, error) {
	var items []Content
	baseDir := "docs/" + contentType

	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory's README
		if path == filepath.Join(baseDir, "README.md") {
			return nil
		}

		if !info.IsDir() && filepath.Ext(path) == ".md" {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			title := mdTitle(content)
			id := idFromPath(path)
			items = append(items, Content{
				Title:     title,
				Id:        id,
				UpdatedAt: info.ModTime(),
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Id < items[j].Id
	})

	if len(items) > n {
		items = items[:n]
	}

	return items, nil
}

// mdTitle returns the markdown title for given markdown content.
func mdTitle(content []byte) string {
	scanner := bufio.NewScanner(bytes.NewReader(content))
	if scanner.Scan() {
		title := strings.TrimSpace(scanner.Text())
		title = strings.TrimLeft(title, "# ")
		return title
	}
	return "Untitled"
}

// idFromPath returns the page id for the given path.
func idFromPath(path string) string {
	// Get parent directory name
	dir := filepath.Dir(path)
	// Extract the last part of the path, which should be the ID
	return filepath.Base(dir)
}
