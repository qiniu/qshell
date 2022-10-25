package file

import (
	"fmt"
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

const (
	rotateFileNameSep = "-"
	returnChar        = "\n"
)

type RotateOption func(file *rotateFile)

func RotateOptionMaxLine(maxLine int64) RotateOption {
	return func(file *rotateFile) {
		file.maxLine = maxLine
	}
}

func RotateOptionMaxSize(maxSize int64) RotateOption {
	return func(file *rotateFile) {
		file.maxSize = maxSize
	}
}

func RotateOptionAppendMode(mode bool) RotateOption {
	return func(file *rotateFile) {
		file.appendMode = mode
	}
}

func RotateOptionFileHeader(header string) RotateOption {
	return func(file *rotateFile) {
		file.fileHeader = header
	}
}

func RotateOptionOnOpenFile(f func(filename string)) RotateOption {
	return func(file *rotateFile) {
		file.onOpenNewFile = f
	}
}

func NewRotateFile(name string, options ...RotateOption) (io.WriteCloser, *data.CodeError) {

	filenameWithSuffix := filepath.Base(name)
	fileExt := filepath.Ext(filenameWithSuffix)
	fileName := strings.TrimSuffix(filenameWithSuffix, fileExt)
	r := &rotateFile{
		mu:              sync.Mutex{},
		maxLine:         0,
		maxSize:         0,
		currentFileLine: 0,
		currentFileSize: 0,
		fileDir:         filepath.Dir(name),
		fileName:        fileName,
		fileExt:         fileExt,
		fileIndex:       0,
		appendMode:      false,
		file:            nil,
	}

	for _, option := range options {
		option(r)
	}

	if r.appendMode {
		if index, gErr := r.getFileIndex(); gErr != nil {
			return nil, gErr
		} else {
			r.fileIndex = index
		}
	}

	if err := r.createFile(); err != nil {
		return nil, data.ConvertError(err)
	}

	return r, nil
}

type rotateFile struct {
	mu              sync.Mutex            //
	maxLine         int64                 // 最大行数
	maxSize         int64                 // 最大文件大小
	currentFileLine int64                 // 当前文件最后一行下表，从 1 开始
	currentFileSize int64                 //
	fileDir         string                //
	fileName        string                // 文件名称，不带扩展名
	fileExt         string                // 文件扩展名
	fileIndex       int                   // 文件的下表
	appendMode      bool                  // 是否为拼接模式
	fileHeader      string                // 新文件的头，创建新文件时自动添加
	file            *os.File              //
	onOpenNewFile   func(filename string) // 打开某个文件后的回调
}

func (r *rotateFile) Write(p []byte) (n int, err error) {
	// 不用 rotate
	if !r.needRotate() {
		return r.file.Write(p)
	}

	return r.writeByRotateWithLock(p)
}

func (r *rotateFile) Close() error {
	if r.file == nil {
		return nil
	}

	return r.file.Close()
}

func (r *rotateFile) needRotate() bool {
	return r.maxSize > 0 || r.maxLine > 0
}

func (r *rotateFile) writeByRotateWithLock(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.writeByRotate(p)
}

func (r *rotateFile) writeByRotate(p []byte) (n int, err error) {

	// 不滚动
	if r.maxLine <= 0 && r.maxSize <= 0 {
		return r.file.Write(p)
	}

	items := strings.Split(string(p), returnChar)
	if len(items) == 0 {
		return len(p), nil
	}

	for i, line := range items {
		if cn, cErr := r.writeLine(i != 0, line); cErr != nil {
			return cn + n, cErr
		} else {
			n += cn
		}
	}
	return len(p), nil
}

func (r *rotateFile) writeLine(isNew bool, line string) (n int, err error) {
	if !isNew && len(line) == 0 {
		return 0, nil
	}

	needCreateNewFile := false

	// 检测行限制，需要则创建新文件
	if isNew {
		r.currentFileLine++

		// 切分的数量比 \n 的数量多 1
		if r.maxLine > 0 && r.currentFileLine > (r.maxLine-1) {
			needCreateNewFile = true
		}
	}

	// 检测文件大小限制，需要则创建新文件
	if !needCreateNewFile && r.maxSize > 0 &&
		r.currentFileSize > 0 && (r.currentFileSize+int64(len(line))) > r.maxSize {
		needCreateNewFile = true
	}

	// 创建新的文件
	if needCreateNewFile {
		if cErr := r.createFile(); cErr != nil {
			return len(returnChar), cErr
		}

		r.currentFileLine = 0
		r.currentFileSize = 0
	} else if isNew && r.currentFileSize != 0 {
		// 非新文件的新行 增加换行符
		line = returnChar + line
	}

	r.currentFileSize += int64(len(line))

	return r.file.WriteString(line)
}

func (r *rotateFile) createFile() error {
	if mErr := os.MkdirAll(r.fileDir, 0766); mErr != nil {
		return mErr
	}

	if cErr := r.Close(); cErr != nil {
		return cErr
	}

	r.file = nil
	r.fileIndex++

	// 打开或创建文件
	flag := os.O_WRONLY | os.O_CREATE
	if !r.appendMode {
		flag |= os.O_TRUNC
	}

	newFileName := fmt.Sprintf("%s%s", r.fileName, r.fileExt)
	if r.needRotate() {
		// 创建 rotate file 路径,
		newFileName = fmt.Sprintf("%s%s%d%s", r.fileName, rotateFileNameSep, r.fileIndex, r.fileExt)
	}
	newFileName = filepath.Join(r.fileDir, newFileName)

	if file, err := os.OpenFile(newFileName, flag, 0666); err != nil {
		return err
	} else {
		r.file = file
	}

	r.currentFileLine = 0
	r.currentFileSize = 0

	// 写头部
	if stat, err := r.file.Stat(); err != nil {
		return err
	} else if stat.Size() <= 0 {
		if _, wErr := r.writeByRotate([]byte(r.fileHeader + returnChar)); wErr != nil {
			return fmt.Errorf("rotate file write header error:%v", wErr)
		}
	}

	if r.onOpenNewFile != nil {
		r.onOpenNewFile(newFileName)
	}

	return nil
}

func (r *rotateFile) getFileIndex() (index int, err *data.CodeError) {

	// 找到最新的文件
	wErr := filepath.WalkDir(r.fileDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		fileNamePrefix := r.fileName + rotateFileNameSep
		if !strings.HasPrefix(path, fileNamePrefix) || !strings.HasSuffix(path, r.fileExt) {
			return nil
		}

		indexString := strings.TrimPrefix(path, fileNamePrefix)
		indexString = strings.TrimSuffix(indexString, r.fileExt)
		if i, aErr := strconv.Atoi(indexString); aErr != nil {
			return nil
		} else if i > index {
			index = i
		}

		return nil
	})
	if wErr != nil {
		return 0, data.ConvertError(wErr)
	}

	return index, data.ConvertError(wErr)
}
