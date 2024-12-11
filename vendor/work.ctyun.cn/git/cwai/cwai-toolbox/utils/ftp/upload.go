package ftp

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/jlaffaye/ftp"

	"work.ctyun.cn/git/cwai/cwai-toolbox/logger"
)

// FTPConfig 包含 FTP 连接的配置信息
type FTPConfig struct {
	Host             string
	Port             int
	Username         string
	Password         string
	DialTimeout      int
	RetryIntervalSec int
}

// UploadFile 上传文件到 FTP 服务器
func UploadFile(ctx context.Context, config FTPConfig, localFilePath, remoteFilePath string, maxRetries int) error {
	var err error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		err = tryUploadFile(ctx, config, localFilePath, remoteFilePath)
		if err == nil {
			return nil // 上传成功
		}

		logger.Errorf(ctx, "attempt %d failed: %v", attempt, err)
		if attempt < maxRetries {
			logger.Infof(ctx, "retrying in %d seconds...", config.RetryIntervalSec)
			time.Sleep(time.Duration(config.RetryIntervalSec) * time.Second)
		}
	}

	return fmt.Errorf("failed to upload file after %d attempts: %v", maxRetries, err)
}

// tryUploadFile 尝试上传文件到 FTP 服务器
func tryUploadFile(ctx context.Context, config FTPConfig, localFilePath string, remoteFilePath string) error {
	// 打开本地文件
	file, err := os.Open(localFilePath)
	if err != nil {
		return fmt.Errorf("failed to open local file: %v", err)
	}
	defer file.Close()

	// 创建一个同步的管道
	pr, pw := io.Pipe()
	defer pr.Close()

	// 连接到 FTP 服务器
	conn, err := ftp.Dial(fmt.Sprintf("%s:%d", config.Host, config.Port),
		ftp.DialWithContext(ctx), ftp.DialWithTimeout(time.Duration(config.DialTimeout)*time.Second))
	if err != nil {
		return fmt.Errorf("failed to connect to FTP server: %v", err)
	}
	defer conn.Quit()

	// 登录 FTP 服务器
	err = conn.Login(config.Username, config.Password)
	if err != nil {
		return fmt.Errorf("failed to login to FTP server: %v", err)
	}

	// 上传文件
	go func() {
		defer pw.Close()
		_, err := io.Copy(pw, file)
		if err != nil {
			logger.Errorf(context.Background(), "failed to copy file content to pipe: %v", err)
		}
	}()

	// 上传文件
	if err := conn.Stor(remoteFilePath, pr); err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}

	return nil
}
