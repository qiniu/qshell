package dockerfile

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ComputeFilesHash 计算 COPY 指令文件的 SHA-256 哈希。
// 匹配 e2b 的 calculateFilesHash 算法：
//  1. hash.Update("COPY src dest")
//  2. 对每个文件（按相对 POSIX 路径排序）：
//     hash.Update(相对 POSIX 路径)
//     hash.Update(文件模式)
//     hash.Update(文件大小)
//     hash.Update(文件内容)
func ComputeFilesHash(src, dest, contextPath string, ignorePatterns []string) (string, error) {
	files, err := collectFiles(src, contextPath, ignorePatterns)
	if err != nil {
		return "", fmt.Errorf("collect files for hash: %w", err)
	}

	h := sha256.New()
	fmt.Fprintf(h, "COPY %s %s", src, dest)

	for _, f := range files {
		h.Write([]byte(f.relPath))
		fmt.Fprintf(h, "%d", f.info.Mode())
		fmt.Fprintf(h, "%d", f.info.Size())

		data, err := os.ReadFile(f.absPath)
		if err != nil {
			return "", fmt.Errorf("read file %s: %w", f.absPath, err)
		}
		h.Write(data)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// CollectAndUpload 收集源文件，创建 gzip tar 归档包，并上传到指定 URL。
func CollectAndUpload(ctx context.Context, uploadURL, src, contextPath string, ignorePatterns []string) error {
	files, err := collectFiles(src, contextPath, ignorePatterns)
	if err != nil {
		return fmt.Errorf("collect files: %w", err)
	}

	// 在内存中创建 tar.gz
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	for _, f := range files {
		header, err := tar.FileInfoHeader(f.info, "")
		if err != nil {
			return fmt.Errorf("create tar header for %s: %w", f.relPath, err)
		}
		header.Name = f.relPath

		if err := tw.WriteHeader(header); err != nil {
			return fmt.Errorf("write tar header for %s: %w", f.relPath, err)
		}

		data, err := os.ReadFile(f.absPath)
		if err != nil {
			return fmt.Errorf("read file %s: %w", f.absPath, err)
		}
		if _, err := tw.Write(data); err != nil {
			return fmt.Errorf("write file %s to tar: %w", f.relPath, err)
		}
	}

	if err := tw.Close(); err != nil {
		return fmt.Errorf("close tar writer: %w", err)
	}
	if err := gw.Close(); err != nil {
		return fmt.Errorf("close gzip writer: %w", err)
	}

	// 通过 PUT 上传
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uploadURL, &buf)
	if err != nil {
		return fmt.Errorf("create upload request: %w", err)
	}
	req.Header.Set("Content-Type", "application/gzip")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("upload files: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// ReadDockerignore 从构建上下文目录读取 .dockerignore 规则。
func ReadDockerignore(contextPath string) []string {
	path := filepath.Join(contextPath, ".dockerignore")
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	var patterns []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, line)
	}
	return patterns
}

// fileEntry 表示一个收集到的文件，用于哈希计算或上传。
type fileEntry struct {
	absPath string
	relPath string // POSIX 风格的相对路径
	info    fs.FileInfo
}

// collectFiles 遍历构建上下文目录，匹配源 glob 模式，
// 过滤忽略规则。结果按相对路径排序。
func collectFiles(src, contextPath string, ignorePatterns []string) ([]fileEntry, error) {
	// 将源路径解析为相对于上下文的绝对路径
	srcAbs := src
	if !filepath.IsAbs(src) {
		srcAbs = filepath.Join(contextPath, src)
	}

	info, err := os.Stat(srcAbs)
	if err != nil {
		// 尝试作为 glob 模式
		return collectGlob(srcAbs, contextPath, ignorePatterns)
	}

	var files []fileEntry
	if info.IsDir() {
		// 遍历目录
		err := filepath.Walk(srcAbs, func(path string, fi fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if fi.IsDir() {
				return nil
			}
			rel, _ := filepath.Rel(contextPath, path)
			posixRel := filepath.ToSlash(rel)
			if isIgnored(posixRel, ignorePatterns) {
				return nil
			}
			files = append(files, fileEntry{
				absPath: path,
				relPath: posixRel,
				info:    fi,
			})
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		rel, _ := filepath.Rel(contextPath, srcAbs)
		posixRel := filepath.ToSlash(rel)
		if !isIgnored(posixRel, ignorePatterns) {
			files = append(files, fileEntry{
				absPath: srcAbs,
				relPath: posixRel,
				info:    info,
			})
		}
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].relPath < files[j].relPath
	})
	return files, nil
}

// collectGlob 使用 glob 模式匹配收集文件。
func collectGlob(pattern, contextPath string, ignorePatterns []string) ([]fileEntry, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid glob pattern %q: %w", pattern, err)
	}

	var files []fileEntry
	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil {
			continue
		}
		if info.IsDir() {
			continue
		}
		rel, _ := filepath.Rel(contextPath, match)
		posixRel := filepath.ToSlash(rel)
		if isIgnored(posixRel, ignorePatterns) {
			continue
		}
		files = append(files, fileEntry{
			absPath: match,
			relPath: posixRel,
			info:    info,
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].relPath < files[j].relPath
	})
	return files, nil
}

// isIgnored 检查路径是否匹配任何忽略规则。
func isIgnored(relPath string, patterns []string) bool {
	for _, p := range patterns {
		negated := false
		if strings.HasPrefix(p, "!") {
			negated = true
			p = p[1:]
		}

		matched, err := filepath.Match(p, relPath)
		if err != nil {
			continue
		}
		if !matched {
			// 也尝试匹配文件基名
			matched, _ = filepath.Match(p, filepath.Base(relPath))
		}
		if matched && !negated {
			return true
		}
	}
	return false
}
