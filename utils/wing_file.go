// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// 00002       2019/06/30   zhaixing       Add function from godfs
// -------------------------------------------------------------------

package utils

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"strings"
	"syscall"

	"github.com/wengoldx/xcore/invar"
	"github.com/wengoldx/xcore/logger"
	"github.com/wengoldx/xcore/secure"
)

// truncate buffer size for file copy
const _truncateBufferSize = 1024 * 30

// Check the filepath whether point to a file.
func IsFile(filepath string) bool {
	if fileinfo, err := os.Stat(filepath); err == nil {
		return !fileinfo.IsDir()
	}
	return false
}

// Check the dirpath whether point to a directory.
func IsDir(dirpath string) bool {
	if fileinfo, err := os.Stat(dirpath); err == nil {
		return fileinfo.IsDir()
	}
	return false
}

// Check the file whether point to a file.
func IsFile2(file *os.File) bool {
	if fileinfo, err := file.Stat(); err == nil {
		return !fileinfo.IsDir()
	}
	return false
}

// Check the file whether point to a directory.
func IsDir2(file *os.File) bool {
	if fileinfo, err := file.Stat(); err == nil {
		return fileinfo.IsDir()
	}
	return false
}

// Check the filepath whether point to a exist file or directory.
//
//	@See use IsFile(), IsFile2() to check exist file exactly.
//	@See use IsDir(), IsDir2() to check exist directory exactly.
func IsExistFile(filepath string) bool {
	fileinfo, err := os.Stat(filepath)
	if err != nil || fileinfo == nil {
		return false
	}
	return true
}

// Check the dirpath whether point to a exist directory, then
// create the directories if unexist, it maybe return error when
// dirpath point to a exist file.
func EnsurePath(dirpath string) error {
	if fileinfo, err := os.Stat(dirpath); err != nil {
		if os.IsNotExist(err) {
			return MakeDirs(dirpath)
		}
		return err
	} else if !fileinfo.IsDir() {
		return invar.NewError("Exist file, invalid directory!")
	}
	return nil
}

// Create a new directories with permission bits, by default perm is 0777.
//
//	Warning: This function not validate dirpath, please ensure it valid.
//	@See use EnsurePath() to check and create directories.
func MakeDirs(dirpath string, perm ...os.FileMode) error {
	if len(perm) > 0 && perm[0] != 0 {
		return os.MkdirAll(dirpath, perm[0])
	}
	return os.MkdirAll(dirpath, os.ModePerm)
}

// Create a new writeonly file with permission bits, by default perm is 0666,
// it will append write datas to file tails.
//
//	Warning: The caller must call file.Close() after writing finished.
func OpenWriteFile(filepath string, perm ...os.FileMode) (*os.File, error) {
	flag := os.O_CREATE | os.O_WRONLY | os.O_APPEND
	if len(perm) > 0 && perm[0] != 0 {
		return os.OpenFile(filepath, flag, perm[0])
	}
	return os.OpenFile(filepath, flag, 0666)
}

// Create a new writonly file with permission bits, by default perm is 0666,
// it will clear file content and write datas from file start.
//
//	Warning: The caller must call file.Close() after writing finished.
func OpenTruncFile(filepath string, perm ...os.FileMode) (*os.File, error) {
	flag := syscall.O_CREAT | os.O_WRONLY | syscall.O_TRUNC
	if len(perm) > 0 && perm[0] != 0 {
		return os.OpenFile(filepath, flag, perm[0])
	}
	return os.OpenFile(filepath, flag, 0666)
}

// Save file datas to target file on override or append mode, by default override
// the datas to file, the function will auto create the unexist directories.
func SaveFile(dirpath, filename string, datas []byte, append ...bool) error {
	dirpath = strings.TrimSuffix(dirpath, "/")
	filename = strings.Trim(filename, "/")

	if len(datas) == 0 {
		return nil // non-need write anything.
	} else if filename == "" || filename == "." || filename == ".." {
		return invar.NewError("Invalid filename '" + filename + "'")
	} else if err := EnsurePath(dirpath); err != nil {
		return err
	}

	var err error
	var tagfile *os.File
	filepath := fmt.Sprintf("%s/%s", dirpath, filename)
	if VarBool(append, false) {
		tagfile, err = OpenTruncFile(filepath)
	} else {
		tagfile, err = OpenWriteFile(filepath)
	}
	if err != nil {
		return err
	}

	// write content to file.
	defer tagfile.Close()
	_, err = tagfile.Write(datas)
	return err
}

// Decode the base64 datas and override the plaintext datas to file.
func SaveB64File(dirpath, filename string, b64datas string) error {
	datas, err := secure.Base64ToByte(b64datas)
	if err != nil {
		return err
	}
	return SaveFile(dirpath, filename, datas)
}

// Delete the target exist file, it will not do anything when filepath
// point to a exist directoy, or the file unexist.
//
//	@See use DeleteFolder() to delete exist folder and anything it contained.
func DeleteFile(filepath string) error {
	if fileinfo, err := os.Stat(filepath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	} else if !fileinfo.IsDir() {
		return os.Remove(filepath)
	}
	return nil
}

// Delete the target exist folder, it will not do anything when dirpath
// point to a exist files, or the folder unexist.
func DeleteFolder(dirpath string) error {
	if fileinfo, err := os.Stat(dirpath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	} else if fileinfo.IsDir() {
		return os.RemoveAll(dirpath)
	}
	return nil
}

// Read file content and encode datas as md5 abstract string.
func FileAbstract(filepath string) (string, error) {
	h := md5.New()
	if file, err := os.Open(filepath); err != nil {
		return "", err
	} else {
		defer file.Close()
		if _, err = io.Copy(h, file); err != nil {
			return "", err
		}
	}

	cipher := h.Sum(nil)
	return hex.EncodeToString(cipher), nil
}

// -----------------------------------------------------------------

// Deprecated: CopyFile Copy source file to traget file.
func CopyFile(src string, dest string) (bool, error) {
	srcfile, err := os.Open(src)
	if err != nil {
		return false, err
	}
	defer srcfile.Close()

	// create or truncate dest file
	destfile, err := OpenTruncFile(dest)
	if err != nil {
		return false, err
	}
	defer destfile.Close()

	// start copying
	result, err := false, nil
	Try(func() {
		buff := make([]byte, _truncateBufferSize)
		for {
			len, state := srcfile.Read(buff)
			if state == io.EOF {
				break
			}
			destfile.Write(buff[0:len])
		}
		result = true
	}, func(execption error) {
		err = execption
	})
	return result, err
}

// Deprecated: CopyFileTo copy source file to given dir.
func CopyFileTo(src string, dir string) (bool, error) {
	srcfile, err := os.Open(src)
	if err != nil {
		return false, err
	}
	defer srcfile.Close()

	// create or truncate dest file
	fileInfo, _ := srcfile.Stat()
	filepath := FixPath(dir) + string(os.PathSeparator) + fileInfo.Name()
	destfile, err := OpenTruncFile(filepath)
	if err != nil {
		return false, err
	}
	defer destfile.Close()

	// start copying
	len, err := io.Copy(destfile, srcfile)
	if err != nil || len != fileInfo.Size() {
		logger.E("copy file err:", err)
		return false, invar.ErrCopyFile
	}
	return true, nil
}

// Deprecated: FixPath fix path, example:
//
// ---
//
//	/aaa/aa\\bb\\cc/d/////     -> /aaa/aa/bb/cc/d
//	E:/aaa/aa\\bb\\cc/d////e/  -> E:/aaa/aa/bb/cc/d/e
//	""                         -> .
//	/                          -> /
func FixPath(input string) string {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return "."
	}

	// replace windows path separator '\' to '/'
	replaceMent := strings.Replace(input, "\\", "/", -1)
	for {
		if strings.Contains(replaceMent, "//") {
			replaceMent = strings.Replace(replaceMent, "//", "/", -1)
			continue
		}

		if replaceMent == "/" {
			return replaceMent
		}

		len := len(replaceMent)
		if len <= 0 {
			break
		}

		if replaceMent[len-1:] == "/" {
			replaceMent = replaceMent[0 : len-1]
		} else {
			break
		}
	}
	return replaceMent
}

// Deprecated: ReadPropFile read properties file on filesystem.
func ReadPropFile(path string) (map[string]string, error) {
	f, e := os.Open(path)
	if e == nil {
		if IsFile2(f) {
			propMap := make(map[string]string)
			reader := bufio.NewReader(f)
			for {
				line, e1 := reader.ReadString('\n')
				if e1 == nil || e1 == io.EOF {
					line = strings.TrimSpace(line)
					if len(line) != 0 && line[0] != '#' {
						// li := strings.Split(line, "=")
						eIndex := strings.Index(line, "=")
						if eIndex == -1 {
							return nil, errors.New("error parameter: '" + line + "'")
						}
						li := []string{line[0:eIndex], line[eIndex+1:]}
						if len(li) > 1 {
							k := strings.TrimSpace(li[0])
							v := strings.TrimSpace(joinLeft(li[1:]))
							propMap[k] = v
						} else {
							return nil, errors.New("error parameter: '" + li[0] + "'")
						}
					}
					if e1 == io.EOF {
						break
					}
				} else {
					// real read error.
					return nil, errors.New("error read from configuration file")
				}
			}
			return propMap, nil
		} else {
			return nil, errors.New("expect file path not directory path")
		}
	} else {
		return nil, e
	}
}

// Deprecated: joinLeft only for ReadPropFile()
func joinLeft(g []string) string {
	if len(g) == 0 {
		return ""
	}
	var bf bytes.Buffer
	for i := range g {
		c := strings.Index(g[i], "#")
		if c == -1 {
			bf.WriteString(g[i])
		} else {
			bf.WriteString(g[i][0:c])
			break
		}
	}
	return bf.String()
}

// HumanReadable format the size number of len.
func HumanReadable(len int64, during int64) string {
	if len < 1024 {
		return strconv.FormatInt(len*1000/during, 10) + "B       "
	} else if len < 1048576 {
		return strconv.FormatInt(len*1.0/1024*1000/during, 10) + "KB       "
	} else if len < 1073741824 {
		return fmt.Sprintf("%.2f", float64(len)/1048576*1000/float64(during)) + "MB       "
	} else {
		return fmt.Sprintf("%.2f", float64(len)/1073741824*1000/float64(during)) + "GB       "
	}
}

// VerifyFile verify upload file and size, it support jpg/jpeg/JPG/JPEG/png/PNG/mp3/mp4 suffix.
func VerifyFile(fh *multipart.FileHeader, maxBytes ...int64) (string, error) {
	suffix := path.Ext(fh.Filename)
	maxSizeInByte := VarInt64(maxBytes, 0)

	switch suffix {
	case ".jpg", ".jpeg", ".JPG", ".JPEG", ".png", ".PNG", ".mp3":
		if maxSizeInByte == 0 {
			maxSizeInByte = 10 << 20 // set default max size
		}

		// image file, must less than 10MB or given size
		if fh.Size > maxSizeInByte {
			if suffix == ".mp3" {
				return "", invar.ErrAudioOverSize
			}
			return "", invar.ErrImgOverSize
		}
	case ".mp4":
		if maxSizeInByte == 0 {
			maxSizeInByte = 500 << 20 // set default max size
		}

		// vedio file, must less than 500MB or given size
		if fh.Size > maxSizeInByte {
			return "", invar.ErrVideoOverSize
		}
	default:
		return "", invar.ErrUnsupportedFile
	}
	return suffix, nil
}

// VerifyFileFormat verify upload file and size in MB.
func VerifyFileFormat(fh *multipart.FileHeader, format string, size int64) (string, error) {
	if len(format) == 0 || size <= 0 {
		return "", invar.ErrInvalidParams
	}

	suffix := path.Ext(fh.Filename)
	switch suffix {
	case format:
		if fh.Size > int64(size<<20) {
			return "", invar.ErrImgOverSize
		}
	default:
		return "", invar.ErrUnsupportedFile
	}
	return suffix, nil
}
