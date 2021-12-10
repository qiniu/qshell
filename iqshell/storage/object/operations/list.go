package operations

import (
	"bufio"
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/log"
	"github.com/qiniu/qshell/v2/iqshell/common/utils"
	"github.com/qiniu/qshell/v2/iqshell/common/workspace"
	"github.com/qiniu/qshell/v2/iqshell/storage/object/rs"
	"os"
	"strconv"
	"strings"
	"time"
)

type ListInfo struct {
	SaveToFile string
	StartDate  string
	EndDate    string
	Suffixes   string
	AppendMode bool
	Readable   bool
	ApiInfo    rs.ListApiInfo
}

func List(info ListInfo) {
	apiInfo := &info.ApiInfo

	startDate, err := info.getStartDate()
	if err != nil {
		log.ErrorF("date parse error: %v", err)
		os.Exit(data.STATUS_ERROR)
	}

	endDate, err := info.getEndDate()
	if err != nil {
		log.ErrorF("date parse error: %v", err)
		os.Exit(data.STATUS_ERROR)
	}

	var resultWitter *os.File
	if len(info.SaveToFile) == 0 {
		resultWitter = os.Stdout
	} else {
		var openErr error
		var mode int
		if info.AppendMode {
			mode = os.O_APPEND | os.O_RDWR
		} else {
			mode = os.O_CREATE | os.O_RDWR | os.O_TRUNC
		}
		resultWitter, openErr = os.OpenFile(info.SaveToFile, mode, 0666)
		if openErr != nil {
			log.Error("Failed to open list result file `%s`", info.SaveToFile)
			return
		}
		defer resultWitter.Close()
	}

	bWriter := bufio.NewWriter(resultWitter)
	objects, err := rs.List(workspace.GetContext(), apiInfo)
	if err != nil {
		log.Error(err)
		return
	}

	suffixes := info.getSuffixes()
	var fSizeValue interface{}
	for object := range objects {
		if filteredByTime(startDate, endDate, object) {
			continue
		}
		if filteredBySuffix(suffixes, object) {
			continue
		}

		if info.Readable {
			fSizeValue = utils.FormatFileSize(object.Fsize)
		} else {
			fSizeValue = object.Fsize
		}

		lineData := fmt.Sprintf("%s\t%v\t%s\t%d\t%s\t%d\t%s\r\n",
			object.Key, fSizeValue, object.Hash,
			object.PutTime, object.MimeType, object.Type, object.EndUser)
		_, err := bWriter.WriteString(lineData)
		if err != nil {
			log.ErrorF("marker: %s", apiInfo.Marker)
			log.ErrorF("listbucket Error: %v", err)
		}

		err = bWriter.Flush()
		if err != nil {
			log.ErrorF("marker: %s", apiInfo.Marker)
			log.ErrorF("listbucket flush Error: %v", err)
		}
	}

	if apiInfo.Marker != "" {
		log.ErrorF("Marker: %s\n", apiInfo.Marker)
	}
}

func parseDate(dateString string) (time.Time, error) {
	if len(dateString) == 0 {
		return time.Time{}, nil
	}

	fields := strings.Split(dateString, "-")
	if len(fields) > 6 {
		return time.Time{}, fmt.Errorf("date format must be year-month-day-hour-minute-second")
	}

	var dateItems [6]int
	for ind, field := range fields {
		field, err := strconv.Atoi(field)
		if err != nil {
			return time.Time{}, fmt.Errorf("date format must be year-month-day-hour-minute-second, each field must be integer")
		}
		dateItems[ind] = field
	}
	return time.Date(dateItems[0], time.Month(dateItems[1]), dateItems[2], dateItems[3], dateItems[4], dateItems[5], 0, time.Local), nil
}

func filteredByTime(startDate, endDate time.Time, item rs.ListItem) bool {
	putTime := time.Unix(item.PutTime/1e7, 0)
	switch {
	case startDate.IsZero() && endDate.IsZero():
		return false
	case !startDate.IsZero() && endDate.IsZero() && putTime.After(startDate):
		return false
	case !endDate.IsZero() && startDate.IsZero() && putTime.Before(endDate):
		return false
	case putTime.After(startDate) && putTime.Before(endDate):
		return false
	default:
		return true
	}
}

func filteredBySuffix(suffixes []string, item rs.ListItem) bool {
	if len(suffixes) == 0 {
		return false
	}

	hasSuffix := false
	for _, s := range suffixes {
		if strings.HasSuffix(item.Key, s) {
			hasSuffix = true
			break
		}
	}
	return hasSuffix
}

func (info ListInfo) getStartDate() (time.Time, error) {
	return parseDate(info.StartDate)
}

func (info ListInfo) getEndDate() (time.Time, error) {
	return parseDate(info.EndDate)
}

func (info ListInfo) getSuffixes() []string {
	sf := make([]string, 0)
	for _, s := range strings.Split(info.Suffixes, ",") {
		s = strings.TrimSpace(s)
		if len(s) > 0 {
			sf = append(sf, strings.TrimSpace(s))
		}
	}
	return nil
}

func (info ListInfo) getResultWitter() []string {
	sf := make([]string, 0)
	for _, s := range strings.Split(info.Suffixes, ",") {
		s = strings.TrimSpace(s)
		if len(s) > 0 {
			sf = append(sf, strings.TrimSpace(s))
		}
	}
	return nil
}