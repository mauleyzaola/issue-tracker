package tecutils

import (
	"io/ioutil"
	"math"
	"os"
	"path"
	"time"
)

const (
	DEFAULT_PERMISSION = 0777
)

//Returns true if the directory exists
func DirectoryExists(fullpath string) (ok bool) {
	fileInfo, err := os.Lstat(fullpath)
	if err != nil {
		ok = false
		return
	}
	ok = len(fileInfo.Name()) != 0
	return
}

//Verifies if the directory already exists, if it does not then it is created
func CreateDirectoryIfNotExist(fullpath string) (err error) {
	ok := DirectoryExists(fullpath)
	if ok {
		return
	}
	err = os.MkdirAll(fullpath, DEFAULT_PERMISSION)
	return
}

type FileInfo struct {
	Name  string
	Size  int64
	IsDir bool
}

type FileLambda func(fullpath string, info *FileInfo, level int)

//Reads a directory content recursively and process the lambda code inside
func ProcessDirectoryContents(fullpath string, recursive bool, fn FileLambda, level int) (err error) {
	files, err := ioutil.ReadDir(fullpath)
	if err != nil {
		return
	}
	for _, file := range files {
		if file.IsDir() && recursive {
			ProcessDirectoryContents(path.Join(fullpath, file.Name()), recursive, fn, level+1)
		}
		if fn != nil {
			fn(fullpath, &FileInfo{Name: file.Name(), Size: file.Size(), IsDir: file.IsDir()}, level)
		}
	}

	return
}

func FileDaysOld(f *os.FileInfo) int {
	file := *f
	interval := time.Now().Sub(file.ModTime()).Hours() / 24
	return int(math.Floor(interval))
}

type Subdirectory struct {
	Name      string `json:"name"`
	FileCount int64  `json:"fileCount"`
	Bytes     int64  `json:"bytes"`
}

func SubdirectoriesInfo(root string) (dirs []*Subdirectory, err error) {
	files, err := ioutil.ReadDir(root)
	var currDir *Subdirectory
	fn := func(localPath string, info *FileInfo, level int) {
		if !info.IsDir {
			currDir.FileCount++
		}
		currDir.Bytes += info.Size
	}
	for _, file := range files {
		if file.IsDir() {
			subdir := &Subdirectory{Name: file.Name()}
			dirs = append(dirs, subdir)
			currDir = &Subdirectory{}
			ProcessDirectoryContents(path.Join(root, subdir.Name), true, fn, 0)
			subdir.Bytes = currDir.Bytes
			subdir.FileCount = currDir.FileCount
		}
	}
	return
}
