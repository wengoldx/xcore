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
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/wengoldx/xcore/invar"
	"github.com/wengoldx/xcore/logger"
	"github.com/wengoldx/xcore/secure"
)

type FileValueTypes interface {
	string | *os.File
} 

// Check the file path (or file object) whether point to a file.
func IsFile[T FileValueTypes](file T) bool {
	switch ft := any(file).(type) {
	case string:
		if info, err := os.Stat(ft); err == nil {
			return !info.IsDir()
		}
	case *os.File:
		if info, err := ft.Stat(); err == nil {
			return !info.IsDir()
		}
	}
	return false
}

// Check the dir path (or file object) whether point to a directory.
func IsDir[T FileValueTypes](dir T) bool {
	switch pt := any(dir).(type) {
	case string:
		if info, err := os.Stat(pt); err == nil {
			return info.IsDir()
		}
	case *os.File:
		if info, err := pt.Stat(); err == nil {
			return info.IsDir()
		}
	}
	return false
}

// Check the file paths whether point to exist file or directory,
// and return the false when anyone file path unexist.
//
// # WARNING:
// 
// This method mark the path exist when case error except NotExist!
//
//	@See use IsFile() to check exist file exactly.
//	@See use IsDir()  to check exist directory exactly.
func IsExistFiles(fps ...string) bool {
	for _, fp := range fps {
		if _, err := os.Stat(fp); err != nil {
			if os.IsNotExist(err) {
				return false
			}
			logger.E("Stat file:", fp, "err:", err)
		}
	}
	return true
}

// Check the dirpath whether point to a exist directory, then
// create the directories if unexist, it maybe return error when
// dirpath point to a exist file.
func EnsurePath(dirpath string) error {
	if info, err := os.Stat(dirpath); err != nil {
		if os.IsNotExist(err) {
			return MakeDirs(dirpath)
		}
		return err
	} else if !info.IsDir() {
		return invar.NewError("Exist file, invalid directory!")
	}
	return nil
}

// Create a new directories with permission bits, by default perm is 0777.
//
// # WARNING:
//	- This function not validate dirpath, please ensure it valid.
//	- Use EnsurePath() to check and create directories.
func MakeDirs(dirpath string, perm ...os.FileMode) error {
	if len(perm) > 0 && perm[0] != 0 {
		return os.MkdirAll(dirpath, perm[0])
	}
	return os.MkdirAll(dirpath, os.ModePerm)
}

// Create a new writeonly file with permission bits, by default perm is 0666,
// it will append write datas to file tails.
//
// # WARNING:
// - The caller must call file.Close() after writing finished.
func OpenWriteFile(fp string, perm ...os.FileMode) (*os.File, error) {
	flag := os.O_CREATE | os.O_WRONLY | os.O_APPEND
	if len(perm) > 0 && perm[0] != 0 {
		return os.OpenFile(fp, flag, perm[0])
	}
	return os.OpenFile(fp, flag, 0666)
}

// Create a new writonly file with permission bits, by default perm is 0666,
// it will clear file content and write datas from file start.
//
// # WARNING:
//	- The caller must call file.Close() after writing finished.
func OpenTruncFile(fp string, perm ...os.FileMode) (*os.File, error) {
	flag := syscall.O_CREAT | os.O_WRONLY | syscall.O_TRUNC
	if len(perm) > 0 && perm[0] != 0 {
		return os.OpenFile(fp, flag, perm[0])
	}
	return os.OpenFile(fp, flag, 0666)
}

// Save the multipart file datas to given local file path.
func SaveMultipartFile(dirpath, filename string, file multipart.File) error {
	if !IsFile(dirpath) {
		if err := MakeDirs(dirpath); err != nil {
			logger.E("Make paths:", dirpath, "err:", err)
			return  err
		}
	}

	dstfile := filepath.Join(dirpath, filename)
	dst, err := OpenTruncFile(dstfile)
	if err != nil {
		return  err 
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		logger.E("Save file:", dstfile, "err:", err)
		return  err
	}

	logger.I("Saved file:", dstfile)
	return nil	
}

// Save the multipart file datas to given local file path from files header.
func SaveByFileHeader(dirpath, filename string, header *multipart.FileHeader)  error {
	partfile, err := header.Open()
	if err != nil {
		logger.E("Open multipart file by header, err:", err)
		return  err
	}
	defer partfile.Close()

	fn := header.Filename
	return SaveMultipartFile(dirpath, fn, partfile)
}

// Save file datas to target file on override or append mode, by default override
// the datas to file, the function will auto create the unexist directories.
func SaveFile(dirpath, filename string, datas []byte, append ...bool) error {
	dirpath, filename = filepath.Clean(dirpath), NormalizePath(filename)

	if len(datas) == 0 {
		return nil // non-need write anything.
	} else if filename == "" || filename == "." || filename == ".." {
		return invar.NewError("Invalid filename '" + filename + "'")
	} else if err := EnsurePath(dirpath); err != nil {
		return err
	}

	var err error
	var tagfile *os.File
	fp := filepath.Join(dirpath, filename)
	if Variable(append, false) {
		tagfile, err = OpenWriteFile(fp)
	} else {
		tagfile, err = OpenTruncFile(fp)
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

// Delete the target exist file, it will not do anything when file path
// point to a exist directoy, or the file unexist.
//
//	@See use DeleteFolder() to delete exist folder and anything it contained.
func DeleteFile(fp string) error {
	if info, err := os.Stat(fp); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	} else if !info.IsDir() {
		return os.Remove(fp)
	}
	return nil
}

// Delete the target exist folder, it will not do anything when dirpath
// point to a exist files, or the folder unexist.
func DeleteFolder(dirpath string) error {
	if info, err := os.Stat(dirpath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	} else if info.IsDir() {
		return os.RemoveAll(dirpath)
	}
	return nil
}

// Read file content and encode datas as md5 abstract string.
func FileAbstract(fp string) (string, error) {
	h := md5.New()
	if file, err := os.Open(fp); err != nil {
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

// Return file suffix without prefix . char, it maybe return empty if
// the filename param is '', '.', '..', and to lower by default.
//
//	Data translate result: `filename.pdf` -> `pdf`
func FileSuffix(filename string, orignl ...bool) string {
	if !Variable(orignl, false) {
		return strings.ToLower(strings.TrimPrefix(filepath.Ext(filename), "."))
	}
	return strings.TrimPrefix(filepath.Ext(filename), ".")
}

// Normalize the given path, it will remove the separator both perfix and tail.
//
//	path: '  /  1//2\\3/..///4/./5/6\\\\' -> '1/2/4/5/6'
//	path: ''                              -> '.'
func NormalizePath(fp string) string {
	separator := string(os.PathSeparator)
	normailze := strings.Trim(filepath.Clean(strings.TrimSpace(fp)), separator)
	return strings.TrimSpace(normailze) 
}

// Split the given file path to dir and base name, the dir called
// clean to trim the / suffix.
//
// For File:
//
//	path: '/1/2/3.doc' -> ['/1/2', '3.doc']
//	path: '1/2/3.doc'  -> ['1/2',  '3.doc']
//	path: '3.doc'      -> ['.',    '3.doc']
//	path: ''           -> ['.',    '.']
//
// For Directory:
//
//	path: '/1/2/3_dir' -> ['/1/2', '3_dir']
//	path: '1/2/3_dir'  -> ['1/2',  '3_dir']
//	path: '3_doc'      -> ['.',    '3_dir']
//
// Use filepath.Split() tail / suffix.
func SplitPath(fp string) (string, string) {
	return filepath.Dir(fp), filepath.Base(fp)
}

// Retrn file simple name without suffix, and trim spaces both start and tails.
//
//	path: '/1/2/   123  .pdf' -> ['123', 'pdf']
//	path: '/1/2/123.pdf'      -> ['123', 'pdf']
//	path: '123.PDF'           -> ['123', 'pdf']
//	path: '123'               -> ['123', ''   ]
//	path: '.pdf'              -> ['',    'pdf']
//	path: ''                  -> ['',    ''   ]
func SplitSuffix(fp string) (string, string) {
	base := filepath.Base(fp)
	suffix := filepath.Ext(base)

	filename := strings.TrimSpace(strings.TrimSuffix(base, suffix))
	suffix = strings.ToLower(strings.TrimPrefix(suffix, "."))
	return filename, suffix
}

// Copy source file to traget file.
func CopyFile(src string, dest string) error {
	srcfile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcfile.Close()

	dstfile, err := OpenTruncFile(dest)
	if err != nil {
		return err
	}
	defer dstfile.Close()

	// start copying
	if _, err := io.Copy(dstfile, srcfile); err != nil {
		logger.E("Copy file err:", err)
		return  invar.ErrCopyFile
	}
	return  nil
}

// Copy source file to given dir.
func CopyFileTo(src string, dir string) error {
	_, fn := SplitPath(src)
	dst := filepath.Join(dir, fn)
	return CopyFile(src, dst)
}

// Read json file content and unmarshal to json object.
func ReadJsonFile(jsonfile string, out any) error {
	if !IsFile(jsonfile) {
		return invar.ErrFileNotFound
	}

	buf, err := os.ReadFile(jsonfile)
	if err != nil {
		logger.E("Read", jsonfile, "err:", err)
		return err
	} else if err := json.Unmarshal(buf, out); err != nil {
		logger.E("Unmarshal err:", err)
		return err
	}
	return nil
}

// Read file content and output as string.
func ReadFileString(txtfile string) string {
	if IsFile(txtfile) {
		buf, err := os.ReadFile(txtfile);
		if err != nil {
			logger.E("Read", txtfile, "err:", err)
			return ""
		}

		content := strings.Trim(string(buf), "\n")
		return strings.TrimSpace(content)
	}
	return ""
}

// Write struct data to json file.
func WriteJsonFile(jsonfile string, data any) error {
	buf, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(jsonfile, buf, 0755)
}

/* ------------------------------------------------------------------- */
/* Deprecated Methods                                                  */
/* ------------------------------------------------------------------- */

// Deprecated: format the size number of len.
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

// Deprecated: verify upload file and size, it support jpg/jpeg/JPG/JPEG/png/PNG/mp3/mp4 suffix.
func VerifyFile(fh *multipart.FileHeader, maxBytes ...int64) (string, error) {
	suffix := filepath.Ext(fh.Filename)
	maxSizeInByte := Variable(maxBytes, 0)

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

// Deprecated: verify upload file and size in MB.
func VerifyFileFormat(fh *multipart.FileHeader, format string, size int64) (string, error) {
	if len(format) == 0 || size <= 0 {
		return "", invar.ErrInvalidParams
	}

	suffix := filepath.Ext(fh.Filename)
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
