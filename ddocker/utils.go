package ddocker

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func DockerIngoreExclude(ignorePath string) []string {
	f, err := os.Open(ignorePath)
	if err != nil {
		log.Println("error fetching dockerignore file:", err)
		return []string{}
	}

	res, err := dockeringoreReadAll(f)
	if err != nil {
		log.Println("error parsing dockeringore file:", err)
		return []string{}
	}

	return res
}

// ReadAll reads a .dockerignore file and returns the list of file patterns
// to ignore. Note this will trim whitespace from each line as well
// as use GO's "clean" func to get the shortest/cleanest path for each.
// copied from: https://github.com/docker/cli/blob/master/vendor/github.com/moby/buildkit/frontend/dockerfile/dockerignore/dockerignore.go
func dockeringoreReadAll(reader io.Reader) ([]string, error) {
	if reader == nil {
		return nil, nil
	}

	scanner := bufio.NewScanner(reader)
	var excludes []string
	currentLine := 0

	utf8bom := []byte{0xEF, 0xBB, 0xBF}
	for scanner.Scan() {
		scannedBytes := scanner.Bytes()
		// We trim UTF8 BOM
		if currentLine == 0 {
			scannedBytes = bytes.TrimPrefix(scannedBytes, utf8bom)
		}
		pattern := string(scannedBytes)
		currentLine++
		// Lines starting with # (comments) are ignored before processing
		if strings.HasPrefix(pattern, "#") {
			continue
		}
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}
		// normalize absolute paths to paths relative to the context
		// (taking care of '!' prefix)
		invert := pattern[0] == '!'
		if invert {
			pattern = strings.TrimSpace(pattern[1:])
		}
		if len(pattern) > 0 {
			pattern = filepath.Clean(pattern)
			pattern = filepath.ToSlash(pattern)
			if len(pattern) > 1 && pattern[0] == '/' {
				pattern = pattern[1:]
			}
		}
		if invert {
			pattern = "!" + pattern
		}

		excludes = append(excludes, pattern)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading .dockerignore: %v", err)
	}
	return excludes, nil
}
