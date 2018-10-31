package qshell

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

/*
generate the file list for the specified directory

@param cacheRootPath - dir to generate cache file
@param cacheResultFile - cache result file path

@return (fileCount, retErr) - total file count and any error meets
*/
func DirCache(cacheRootPath string, cacheResultFile string) (fileCount int64, retErr error) {
	//check dir
	rootPathFileInfo, statErr := os.Stat(cacheRootPath)
	if statErr != nil {
		retErr = statErr
		logs.Error("Failed to stat path `%s`, %s", cacheRootPath, statErr)
		return
	}

	if !rootPathFileInfo.IsDir() {
		retErr = errors.New("dircache failed")
		logs.Error("Dir cache failed, `%s` should be a directory rather than a file", cacheRootPath)
		return
	}

	var cacheResultFh io.Writer
	if cacheResultFile == "stdout" {
		cacheResultFh = os.Stdout
	} else {
		//create result file
		cResultFh, createErr := os.Create(cacheResultFile)
		if createErr != nil {
			retErr = createErr
			logs.Error("Failed to open cache file `%s`, %s", cacheResultFile, createErr)
			return
		}
		defer cResultFh.Close()
		cacheResultFh = cResultFh
	}

	bWriter := bufio.NewWriter(cacheResultFh)
	defer bWriter.Flush()

	//walk start
	walkStart := time.Now()

	logs.Informational("Walk `%s` start from %s", cacheRootPath, walkStart.String())
	filepath.Walk(cacheRootPath, func(path string, fi os.FileInfo, walkErr error) error {
		var retErr error
		//check error
		if walkErr != nil {
			logs.Error("Walk through `%s` error, %s", path, walkErr)

			//skip this dir
			retErr = filepath.SkipDir
		} else {
			if fi.IsDir() {
				logs.Debug("Walking through `%s`", path)
			} else {
				var relPath string
				if cacheRootPath == "." {
					relPath = path
				} else {
					relPath = strings.TrimPrefix(strings.TrimPrefix(path, cacheRootPath), string(os.PathSeparator))
				}

				fsize := fi.Size()
				//Unit is 100ns
				flmd := fi.ModTime().UnixNano() / 100

				logs.Debug("Meet file `%s`, size: %d, modtime: %d", relPath, fsize, flmd)
				fmeta := fmt.Sprintf("%s\t%d\t%d\n", relPath, fsize, flmd)
				if _, err := bWriter.WriteString(fmeta); err != nil {
					logs.Error("Failed to write data `%s` to cache file `%s`", fmeta, cacheResultFile)
				} else {
					fileCount += 1
				}
			}
		}
		return retErr
	})

	if fErr := bWriter.Flush(); fErr != nil {
		logs.Error("Failed to flush to cache file `%s`", cacheResultFile)
		retErr = fErr
		return
	}

	walkEnd := time.Now()
	logs.Debug("Walk `%s` end at %s", cacheRootPath, walkEnd.String())
	logs.Debug("Walk `%s` last for %s", cacheRootPath, time.Since(walkStart))
	logs.Debug("Total file count cached %d", fileCount)
	return
}
