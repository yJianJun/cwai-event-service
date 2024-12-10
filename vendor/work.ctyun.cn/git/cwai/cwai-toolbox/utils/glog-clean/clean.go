package glog_clean

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/golang/glog"
)

// isSymlink 是否是符号链接
func isSymlink(file os.FileInfo) bool {
	return file.Mode()&os.ModeSymlink != 0
}

// CleanerOption: Cleaner选项
type CleanerOption func(c *Cleaner)

// WithInterval: 设置interval参数
func WithInterval(interval time.Duration) CleanerOption {
	return func(c *Cleaner) {
		c.interval = interval
	}
}

// WithRotate: 设置rotate参数
func WithRotate(rotate int) CleanerOption {
	return func(c *Cleaner) {
		c.rotate = rotate
	}
}

// Cleaner: Glog日志清理器
type Cleaner struct {
	dirname  string        // 清理的目标目录
	interval time.Duration // 保留日志文件的时间，即保留now ~ now-interval时间段内的日志
	rotate   int           // 最多保留多少个日志文件

	prefixGetter  func() string
	tagNameGetter func(string) string
}

func NewCleaner(dirname string, opts ...CleanerOption) *Cleaner {
	c := &Cleaner{
		dirname:       dirname,
		prefixGetter:  logPrefix,
		tagNameGetter: getTagName,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// listLogFiles: 获取目标目录下该程序的日志文件map，格式为：map[tag]map[string]os.FileInfo
func (c *Cleaner) listLogFiles() (map[string]map[string]os.FileInfo, error) {
	prefix := c.prefixGetter()
	files, err := ioutil.ReadDir(c.dirname)
	if err != nil {
		return nil, err
	}
	taggedFiles := map[string]map[string]os.FileInfo{}
	for _, file := range files {
		if !strings.HasPrefix(file.Name(), prefix) {
			continue
		}
		if !file.Mode().IsRegular() {
			continue
		}
		tag := c.tagNameGetter(file.Name())
		if len(tag) == 0 {
			continue
		}
		if _, exists := taggedFiles[tag]; !exists {
			taggedFiles[tag] = map[string]os.FileInfo{}
		}
		taggedFiles[tag][file.Name()] = file
	}
	return taggedFiles, nil
}

// getOutdatedFiles: 获取过期文件
func (c *Cleaner) getOutdatedFiles() ([]os.FileInfo, error) {
	taggedFiles, err := c.listLogFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to list log files: %v", err)
	}
	if c.rotate > 0 {
		return c.getOutdatedFilesByRotate(taggedFiles), nil
	} else if c.interval > 0 {
		return c.getOutdatedFilesByInterval(taggedFiles), nil
	}
	return nil, nil
}

// getOutdatedFilesByRotate: 根据rotate保留个数获取过期日志文件
func (c *Cleaner) getOutdatedFilesByRotate(taggedFiles map[string]map[string]os.FileInfo) []os.FileInfo {
	outdatedFiles := []os.FileInfo{}
	for _, files := range taggedFiles {
		if len(files) <= c.rotate {
			continue
		}
		fileNames := []string{}
		for _, f := range files {
			fileNames = append(fileNames, f.Name())
		}
		sort.Strings(fileNames)
		for idx := 0; idx < len(files)-c.rotate; idx++ {
			outdatedFiles = append(outdatedFiles, files[fileNames[idx]])
		}
	}
	return outdatedFiles
}

// getOutdatedFilesByInterval: 根据interval保留时间获取过期日志文件
func (c *Cleaner) getOutdatedFilesByInterval(taggedFiles map[string]map[string]os.FileInfo) []os.FileInfo {
	now := time.Now()
	interval := c.interval
	outdatedFiles := []os.FileInfo{}
	for _, files := range taggedFiles {
		for _, file := range files {
			if now.Sub(file.ModTime()) > interval {
				outdatedFiles = append(outdatedFiles, file)
			}
		}
	}
	return outdatedFiles
}

// cleanFiles: 清除文件
func (c *Cleaner) cleanFiles(files []os.FileInfo) {
	for _, file := range files {
		filePath := filepath.Join(c.dirname, file.Name())
		remove(filePath)
	}
}

func (c *Cleaner) fileNames(files []os.FileInfo) []string {
	names := []string{}
	for _, file := range files {
		names = append(names, file.Name())
	}
	return names
}

// clean: 清理过期glog日志
func (c *Cleaner) clean() error {
	files, err := c.getOutdatedFiles()
	if err != nil {
		glog.Errorf("failed to get outdated files: dirname=%v, interval=%v, rotate=%d, error=%v", c.dirname, c.interval, c.rotate, err)
		return err
	}
	glog.Infof("outdated files in %s: %v", c.dirname, c.fileNames(files))
	c.cleanFiles(files)
	return nil
}

// listOutdatedFiles 获取指定目录下过期的文件列表
func listOutdatedFiles(dirname string, interval time.Duration) ([]os.FileInfo, error) {
	now := time.Now()
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	// 找到过期的普通文件
	outdatedFiles := []os.FileInfo{}
	outdatedFilesMap := make(map[string]os.FileInfo)
	for _, file := range files {
		if !file.Mode().IsRegular() {
			continue
		}
		if now.Sub(file.ModTime()) > interval {
			outdatedFilesMap[file.Name()] = file
			outdatedFiles = append(outdatedFiles, file)
		}
	}
	// 找到过期的软连接
	// 如果软连接指向的文件过期，那么该软连接也过期
	for _, file := range files {
		if !isSymlink(file) {
			continue
		}
		filePath := filepath.Join(dirname, file.Name())
		targetPath, err := filepath.EvalSymlinks(filePath)
		if err != nil {
			glog.Errorf("failed to eval symlinks for '%s': %v", filePath, err)
			continue
		}
		if filepath.Dir(targetPath) != filepath.Clean(dirname) {
			// 外部软链
			continue
		}
		if _, exists := outdatedFilesMap[filepath.Base(targetPath)]; exists {
			outdatedFiles = append(outdatedFiles, file)
		}
	}
	glog.Infof("files in '%s': total=%d, outdated=%d", dirname, len(files), len(outdatedFiles))
	return outdatedFiles, nil
}

// remove 删除文件
func remove(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		glog.Errorf("failed to remove file '%s': %v", filePath, err)
	} else {
		glog.Infof("succeeded to remove file: %s", filepath.Base(filePath))
	}
	return err
}

// CleanGlog 清理glog日志目录
func CleanGlog(dirname string, interval time.Duration) error {
	cleaner := NewCleaner(dirname, WithInterval(interval))
	return cleaner.clean()
}

// CleanGlogWithOptions: 清理glog日志目录
func CleanGlogWithOptions(dirname string, opts ...CleanerOption) error {
	cleaner := NewCleaner(dirname, opts...)
	return cleaner.clean()
}
