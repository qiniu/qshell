package dockerfile

import (
	"archive/tar"
	"bufio"
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

		file, err := os.Open(f.absPath)
		if err != nil {
			return "", fmt.Errorf("read file %s: %w", f.absPath, err)
		}
		_, err = io.Copy(h, file)
		file.Close()
		if err != nil {
			return "", fmt.Errorf("hash file %s: %w", f.absPath, err)
		}
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// CollectAndUpload 收集源文件，通过流式方式创建 gzip tar 归档包并上传到指定 URL。
// 使用 io.Pipe 避免将整个归档加载到内存中。
func CollectAndUpload(ctx context.Context, uploadURL, src, contextPath string, ignorePatterns []string) error {
	files, err := collectFiles(src, contextPath, ignorePatterns)
	if err != nil {
		return fmt.Errorf("collect files: %w", err)
	}

	// 使用 io.Pipe 流式创建 tar.gz 并直接上传
	pr, pw := io.Pipe()

	go func() {
		gw := gzip.NewWriter(pw)
		tw := tar.NewWriter(gw)

		var writeErr error
		for _, f := range files {
			header, err := tar.FileInfoHeader(f.info, "")
			if err != nil {
				writeErr = fmt.Errorf("create tar header for %s: %w", f.relPath, err)
				break
			}
			header.Name = f.relPath

			if err := tw.WriteHeader(header); err != nil {
				writeErr = fmt.Errorf("write tar header for %s: %w", f.relPath, err)
				break
			}

			file, err := os.Open(f.absPath)
			if err != nil {
				writeErr = fmt.Errorf("open file %s: %w", f.absPath, err)
				break
			}
			_, err = io.Copy(tw, file)
			file.Close()
			if err != nil {
				writeErr = fmt.Errorf("write file %s to tar: %w", f.relPath, err)
				break
			}
		}

		tw.Close()
		gw.Close()
		if writeErr != nil {
			pw.CloseWithError(writeErr)
		} else {
			pw.Close()
		}
	}()

	// 通过 PUT 上传
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uploadURL, pr)
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
// 会验证源路径不会逃逸出构建上下文目录。
func collectFiles(src, contextPath string, ignorePatterns []string) ([]fileEntry, error) {
	contextAbs, err := filepath.Abs(contextPath)
	if err != nil {
		return nil, fmt.Errorf("resolve context path: %w", err)
	}

	// 将源路径解析为相对于上下文的绝对路径
	srcAbs := src
	if !filepath.IsAbs(src) {
		srcAbs = filepath.Join(contextPath, src)
	}
	srcAbs, err = filepath.Abs(srcAbs)
	if err != nil {
		return nil, fmt.Errorf("resolve source path: %w", err)
	}

	// 路径穿越检查：确保源路径在构建上下文内
	if !isWithinContext(srcAbs, contextAbs) {
		return nil, fmt.Errorf("path %q escapes build context %q", src, contextPath)
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
// 过滤掉逃逸出构建上下文的匹配结果。
func collectGlob(pattern, contextPath string, ignorePatterns []string) ([]fileEntry, error) {
	contextAbs, _ := filepath.Abs(contextPath)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid glob pattern %q: %w", pattern, err)
	}

	var files []fileEntry
	for _, match := range matches {
		absMatch, _ := filepath.Abs(match)
		if !isWithinContext(absMatch, contextAbs) {
			continue
		}
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

// isIgnored 检查路径是否匹配忽略规则。
// 使用 last-match-wins 语义，与 Docker 的 .dockerignore 规范一致：
// 否定规则（!pattern）可以重新包含之前被排除的文件。
func isIgnored(relPath string, patterns []string) bool {
	ignored := false
	for _, p := range patterns {
		negated := strings.HasPrefix(p, "!")
		if negated {
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
		if matched {
			ignored = !negated
		}
	}
	return ignored
}

// isWithinContext 检查路径是否在构建上下文目录内。
func isWithinContext(path, contextAbs string) bool {
	return strings.HasPrefix(path+string(filepath.Separator), contextAbs+string(filepath.Separator)) ||
		path == contextAbs
}
