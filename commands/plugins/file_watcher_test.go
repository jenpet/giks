package plugins

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
	"time"
)

func TestFileWatcher(t *testing.T) {
	fileWatcherTests := []struct {
		name           string
		pattern        string
		files          string
		okExpected     bool
		errExpected    bool
		cmdExcExpected bool
	}{
		{
			"files match pattern",
			".*.go",
			"../foo/bar.go bar.go foo.go",
			true,
			false,
			true,
		},
		{
			"files dont match pattern",
			".*.go",
			"../foo/bar.js bar.js foo.ts",
			true,
			false,
			false,
		},
		{
			"empty pattern",
			"",
			"foo.go bar.ts",
			false,
			true,
			false,
		},
		{
			"empty files",
			".*.go",
			"",
			true,
			false,
			false,
		},
	}
	fw, _ := Get("file-watcher")
	for _, tt := range fileWatcherTests {
		t.Run(tt.name, func(t *testing.T) {
			file := testFileName()
			_ = os.MkdirAll(path.Dir(file), 0777)
			vars := map[string]string{
				"FILE_WATCHER_PATTERN":    tt.pattern,
				"FILE_WATCHER_COMMAND":    "touch " + file,
				"FILE_WATCHER_FILES_LIST": tt.files,
			}
			ok, err := fw.Run("", "pre-commit", vars, nil)
			assert.Equal(t, tt.okExpected, ok, "expected bool does not match executed plugin result")
			assert.Equal(t, tt.errExpected, err != nil, "error expectation and result does not match")
			fh, err := os.Stat(file)
			if tt.cmdExcExpected {
				assert.NotNilf(t, fh, "command should have been executed successfully")
				assert.NoError(t, err, "command should have been executed successfully")
				_ = os.Remove(file)
			} else {
				assert.True(t, os.IsNotExist(err), "command should not have been executed")
			}
		})
	}
}

func testFileName() string {
	return fmt.Sprintf("../../test/output/%d.out", time.Now().UnixNano())
}
