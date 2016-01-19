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

func (this *DirCache) Cache(cacheRootPath string, cacheResultFile string) (fileCount int64) {
	cacheResultFileH, err := os.Create(cacheResultFile)
	if err != nil {
		log.Errorf("Failed to open cache file `%s'", cacheResultFile)
		return
	}
	defer cacheResultFileH.Close()
	bWriter := bufio.NewWriter(cacheResultFileH)
	walkStart := time.Now()
	log.Debug(fmt.Sprintf("Walk `%s' start from `%s'", cacheRootPath, walkStart.String()))
	filepath.Walk(cacheRootPath, func(path string, fi os.FileInfo, err error) error {
		var retErr error
		//log.Debug(fmt.Sprintf("Walking through `%s'", cacheRootPath))
		if err != nil {
			retErr = err
		} else {
			if !fi.IsDir() {
				relPath := strings.TrimPrefix(strings.TrimPrefix(path, cacheRootPath), string(os.PathSeparator))
				fsize := fi.Size()
				//Unit is 100ns
				flmd := fi.ModTime().UnixNano() / 100
				//log.Debug(fmt.Sprintf("Hit file `%s' size: `%d' mode time: `%d`", relPath, fsize, flmd))
				fmeta := fmt.Sprintln(fmt.Sprintf("%s\t%d\t%d", relPath, fsize, flmd))
				if _, err := bWriter.WriteString(fmeta); err != nil {
					log.Errorf("Failed to write data `%s' to cache file", fmeta)
					retErr = err
				}
				fileCount += 1
			}
		}
		return retErr
	})
	if err := bWriter.Flush(); err != nil {
		log.Errorf("Failed to flush to cache file `%s'", cacheResultFile)
	}

	walkEnd := time.Now()
	log.Debug(fmt.Sprintf("Walk `%s' end at `%s'", cacheRootPath, walkEnd.String()))
	log.Debug(fmt.Sprintf("Walk `%s' last for `%s'", cacheRootPath, time.Since(walkStart)))
	return
}
