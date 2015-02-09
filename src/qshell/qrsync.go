package qshell

import (
	"bufio"
	"fmt"
	"github.com/qiniu/log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func CheckQrsync(dirCacheResultFile string, listBucketResultFile string, ignoreLocalDir bool, prefix string) {
	dirCacheFp, err := os.Open(dirCacheResultFile)
	if err != nil {
		log.Error("Open DirCacheResultFile failed,", err)
		return
	}
	defer dirCacheFp.Close()
	listBucketFp, err := os.Open(listBucketResultFile)
	if err != nil {
		log.Error("Open ListBucketResultFile failed,", err)
		return
	}
	defer dirCacheFp.Close()
	//read all list result
	listResultDataMap := make(map[string]int64)
	lbfScanner := bufio.NewScanner(listBucketFp)
	lbfScanner.Split(bufio.ScanLines)
	for lbfScanner.Scan() {
		line := strings.TrimSpace(lbfScanner.Text())
		items := strings.Split(line, "\t")
		if len(items) >= 2 {
			fname := items[0]
			fsize, err := strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				continue
			} else {
				listResultDataMap[fname] = fsize
			}
		} else {
			continue
		}
	}
	allUploaded := true
	//iter each local file and compare
	dcrScanner := bufio.NewScanner(dirCacheFp)
	dcrScanner.Split(bufio.ScanLines)
	for dcrScanner.Scan() {
		line := strings.TrimSpace(dcrScanner.Text())
		items := strings.Split(line, "\t")
		if len(items) >= 2 {
			localFname := items[0]
			remoteFname := localFname
			if ignoreLocalDir {
				ldx := strings.LastIndex(remoteFname, string(os.PathSeparator))
				remoteFname = remoteFname[ldx+1:]
			}
			if prefix != "" {
				remoteFname = prefix + remoteFname
			}
			//convert \ to / under windows
			if runtime.GOOS == "windows" {
				remoteFname = strings.Replace(remoteFname, "\\", "/", -1)
			}
			fsize, err := strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				continue
			}
			if rFsize, ok := listResultDataMap[remoteFname]; ok {
				if rFsize != fsize {
					allUploaded = false
					log.Error("Uploaded but size not ok for file ", localFname, "=>", remoteFname)
				}
			} else {
				allUploaded = false
				log.Error("Not uploaded for file ", localFname, "=>", remoteFname)
			}
		} else {
			continue
		}
	}
	if allUploaded {
		fmt.Println("All Uploaded!")
	}
}
