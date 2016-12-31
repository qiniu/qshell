package qshell

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"qiniu/log"
	"strings"
	"time"
)

type DirCache struct {
}

func (this *DirCache) Cache(cacheRootPath string, cacheResultFile string) (fileCount int) {
	cacheResultFh, err := os.Create(cacheResultFile)
	if err != nil {
		log.Errorf("Failed to open cache file `%s`", cacheResultFile)
		return
	}
	defer cacheResultFh.Close()

	bWriter := bufio.NewWriter(cacheResultFh)

	//walk start
	walkStart := time.Now()
	log.Infof("Walk `%s` start from %s", cacheRootPath, walkStart.String())
	filepath.Walk(cacheRootPath, func(path string, fi os.FileInfo, err error) error {
		var retErr error
		if fi.IsDir() {
			log.Infof("Walking through `%s`", path)
		}

		//check error
		if err != nil {
			log.Errorf("Walk through `%s` error, %s", path, err)
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

	if err := bWriter.Flush(); err != nil {
		log.Errorf("Failed to flush to cache file `%s`", cacheResultFile)
	}

	walkEnd := time.Now()
	log.Infof("Walk `%s` end at %s", cacheRootPath, walkEnd.String())
	log.Infof("Walk `%s` last for %s", cacheRootPath, time.Since(walkStart))
	log.Infof("Total file count %d", fileCount)
	return
}
