package handler

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type FileInfo struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	IsDir bool   `json:"isDir"`
}

// ListDirectories 列出目录内容
func (h *Handler) ListDirectories(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		path = "."
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "读取目录失败",
			"error":   err.Error(),
		})
		return
	}

	var files []FileInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		files = append(files, FileInfo{
			Name:  entry.Name(),
			Path:  filepath.Join(path, entry.Name()),
			IsDir: info.IsDir(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": files,
	})
}

// ListDrives 列出所有盘符（Windows）
func (h *Handler) ListDrives(c *gin.Context) {
	var drives []FileInfo
	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		path := string(drive) + ":\\"
		if _, err := os.Stat(path); err == nil {
			drives = append(drives, FileInfo{
				Name:  path,
				Path:  path,
				IsDir: true,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": drives,
	})
}
