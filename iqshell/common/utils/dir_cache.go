package utils

import (
	"bufio"
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/qiniu/qshell/v2/iqshell/common/log"
)

// DirCache
// generate the file list for the specified directory
// @param cacheRootPath - dir to generate cache file
// @param cacheResultFile - cache result file path
// @return (fileCount, retErr) - total file count and any error meets
func DirCache(cacheRootPath string, cacheResultFile string) (int64, *data.CodeError) {
	//check dir
	rootPathFileInfo, statErr := os.Stat(cacheRootPath)
	if statErr != nil {
		log.ErrorF("Failed to stat path `%s`, %s", cacheRootPath, statErr)
		return 0, data.NewEmptyError().AppendError(statErr)
	}

	if !rootPathFileInfo.IsDir() {
		log.ErrorF("Dir cache failed, `%s` should be a directory rather than a file", cacheRootPath)
		return 0, data.NewEmptyError().AppendDesc("dircache failed")
	}

	var cacheResultFh io.Writer
	if cacheResultFile == "stdout" {
		cacheResultFh = os.Stdout
	} else {
		catchDir := filepath.Dir(cacheResultFile)
		mkErr := os.MkdirAll(catchDir, os.ModePerm)
		if mkErr != nil {
			log.ErrorF("Failed to create cache dir `%s`, %s", catchDir, mkErr)
			return 0, data.NewEmptyError().AppendError(mkErr)
		}

		//create result file
		cResultFh, createErr := os.Create(cacheResultFile)
		if createErr != nil {
			log.ErrorF("Failed to open cache file `%s`, %s", cacheResultFile, createErr)
			return 0, data.NewEmptyError().AppendError(createErr)
		}
		defer cResultFh.Close()
		cacheResultFh = cResultFh
	}

	bWriter := bufio.NewWriter(cacheResultFh)
	defer bWriter.Flush()

	//walk start
	walkStart := time.Now()

	log.DebugF("Walk `%s` start from %s", cacheRootPath, walkStart.String())

	var fileCount int64 = 0
	filepath.Walk(cacheRootPath, func(path string, fi os.FileInfo, walkErr error) error {
		var retErr error
		//check error
		if walkErr != nil {
			log.ErrorF("Walk through `%s` error, %s", path, walkErr)

			//skip this dir
			retErr = filepath.SkipDir
		} else {
			if fi.IsDir() {
				log.DebugF("Walking through `%s`", path)
			} else {
				var relativePath string
				trimPrefix := cacheRootPath
				if strings.HasPrefix(trimPrefix, ".") {
					trimPrefix = strings.TrimPrefix(strings.TrimPrefix(trimPrefix, "."), string(os.PathSeparator))
				}
				relativePath = strings.TrimPrefix(strings.TrimPrefix(path, trimPrefix), string(os.PathSeparator))
				log.DebugF("cacheRootPath:`%s` path:`%s` relativePath:`%s`", cacheRootPath, path, relativePath)

				fsize := fi.Size()
				//Unit is 100ns
				flmd := fi.ModTime().UnixNano() / 100

				log.DebugF("Meet file `%s`, size: %d, modtime: %d", relativePath, fsize, flmd)
				fmeta := fmt.Sprintf("%s\t%d\t%d\n", relativePath, fsize, flmd)
				if _, err := bWriter.WriteString(fmeta); err != nil {
					log.ErrorF("Failed to write data `%s` to cache file `%s`", fmeta, cacheResultFile)
				} else {
					fileCount += 1
				}
			}
		}
		return retErr
	})

	if fErr := bWriter.Flush(); fErr != nil {
		log.ErrorF("Failed to flush to cache file `%s`", cacheResultFile)
		return 0, data.NewEmptyError().AppendError(fErr)
	}

	walkEnd := time.Now()
	log.DebugF("Walk `%s` end at %s", cacheRootPath, walkEnd.String())
	log.DebugF("Walk `%s` last for %s", cacheRootPath, time.Since(walkStart))
	log.DebugF("Total file count cached %d", fileCount)
	return fileCount, nil
}
