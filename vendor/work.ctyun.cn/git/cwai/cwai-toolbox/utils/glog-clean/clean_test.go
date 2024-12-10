package glog_clean

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListOutdatedFiles(t *testing.T) {
	dirname := "/tmp"
	interval := time.Hour * 24
	files, err := listOutdatedFiles(dirname, interval)
	assert.Nil(t, err)
	for i, file := range files {
		t.Logf("files[%d]=%s", i, file.Name())
	}
}

func TestCleaner(t *testing.T) {
	dirname := "/tmp"
	cleaner := NewCleaner(dirname, WithRotate(1))
	cleaner.prefixGetter = func() string {
		return fmt.Sprintf("ut.hostname.zhangsan")
	}
	names := []string{
		// vc-scheduler.volcano-scheduler-5b7b6bcd7c-slqnp.root.log.WARNING.20210303-183058.1
		fmt.Sprintf("%s.%s", cleaner.prefixGetter(), "log.INFO.1"),
		fmt.Sprintf("%s.%s", cleaner.prefixGetter(), "log.INFO.2"),
		fmt.Sprintf("%s.%s", cleaner.prefixGetter(), "log.INFO.3"),
	}
	for _, name := range names {
		os.Create(filepath.Join(dirname, name))
	}

	// getOutdatedFiles
	files, err := cleaner.getOutdatedFiles()
	require.Nil(t, err)
	t.Logf("outdated files: %+v", cleaner.fileNames(files))
	assert.Equal(t, []string{names[0], names[1]}, cleaner.fileNames(files))
	// clean
	err = cleaner.clean()
	require.Nil(t, err)
}
