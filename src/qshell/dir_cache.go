package qshell

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"qiniu/log"
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
		log.Errorf("Failed to stat path `%s`, %s", cacheRootPath, statErr)
		return
	}

	if !rootPathFileInfo.IsDir() {
		retErr = errors.New("dircache failed")
		log.Errorf("Dir cache failed, `%s` should be a directory rather than a file", cacheRootPath)
		return
	}

	//create result file
	cacheResultFh, createErr := os.Create(cacheResultFile)
	if createErr != nil {
		retErr = createErr
		log.Errorf("Failed to open cache file `%s`, %s", cacheResultFile, createErr)
		return
	}
	defer cacheResultFh.Close()

	bWriter := bufio.NewWriter(cacheResultFh)
	defer bWriter.Flush()

	//walk start
	walkStart := time.Now()

	log.Infof("Walk `%s` start from %s", cacheRootPath, walkStart.String())
	filepath.Walk(cacheRootPath, func(path string, fi os.FileInfo, walkErr error) error {
		var retErr error
		if fi.IsDir() {
			log.Infof("Walking through `%s`", path)
		}

		//check error
		if walkErr != nil {
			log.Errorf("Walk through `%s` error, %s", path, walkErr)

			//skip this dir
			retErr = filepath.SkipDir
		} else {
			if !fi.IsDir() {
				var relPath string
				if cacheRootPath == "." {
					relPath = path
				} else {
					relPath = strings.TrimPrefix(strings.TrimPrefix(path, cacheRootPath), string(os.PathSeparator))
				}

				fsize := fi.Size()
				//Unit is 100ns
				flmd := fi.ModTime().UnixNano() / 100

				log.Debugf("Meet file `%s`, size: %d, modtime: %d", relPath, fsize, flmd)
				fmeta := fmt.Sprintf("%s\t%d\t%d\n", relPath, fsize, flmd)
				if _, err := bWriter.WriteString(fmeta); err != nil {
					log.Errorf("Failed to write data `%s` to cache file `%s`", fmeta, cacheResultFile)
				} else {
					fileCount += 1
				}
			}
		}
		return retErr
	})

	if fErr := bWriter.Flush(); fErr != nil {
		log.Errorf("Failed to flush to cache file `%s`", cacheResultFile)
		retErr = fErr
		return
	}

	walkEnd := time.Now()
	log.Infof("Walk `%s` end at %s", cacheRootPath, walkEnd.String())
	log.Infof("Walk `%s` last for %s", cacheRootPath, time.Since(walkStart))
	log.Infof("Total file count cached %d", fileCount)
	return
}
