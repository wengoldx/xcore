// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package wrpc

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/wengoldx/xcore/invar"
	"github.com/wengoldx/xcore/secure"
	"github.com/wengoldx/xcore/utils"
	acc "github.com/wengoldx/xcore/wrpc/accservice/proto"
	wss "github.com/wengoldx/xcore/wrpc/webss/proto"
)

// Signed upload url and MinIO bucket path.
type UrlPath struct {
	Url  string // Signed url to upload file
	Path string // Real bucket path of MinIO service
}

// ----------------------------------------
// For Wss GRPC agent
// ----------------------------------------

// Obtain a single pre-signed url for upload file to minio storage.
func (stub *GrpcStub) SignUrl(res, suffix string, addition ...string) (*UrlPath, error) {
	if stub.Wss == nil {
		return nil, invar.ErrInvalidClient
	} else if res == "" || suffix == "" {
		return nil, invar.ErrInvalidParams
	}

	add := utils.GetVariable(addition, "").(string)
	param := &wss.Sign{Res: res, Add: add, Suffix: suffix}
	signurl, err := stub.Wss.SignFileUrl(context.Background(), param)
	if err != nil {
		return nil, err
	}

	surl, _ := secure.DecodeBase64(signurl.Url)
	path, _ := secure.DecodeBase64(signurl.Path) // Format as '/bucket/file/paths'
	if surl == "" || path == "" {
		return nil, invar.ErrUnsupportFormat
	}
	return &UrlPath{surl, path}, nil
}

// Obtain multiple pre-signed urls for upload files to minio storage.
func (stub *GrpcStub) SignUrls(res string, suffixes []string, addition ...string) ([]*UrlPath, error) {
	if stub.Wss == nil {
		return nil, invar.ErrInvalidClient
	} else if res == "" || len(suffixes) == 0 {
		return nil, invar.ErrInvalidParams
	}

	add := utils.GetVariable(addition, "").(string)
	param := &wss.Signs{Res: res, Add: add, Suffixes: suffixes}
	signurls, err := stub.Wss.SignFileUrls(context.Background(), param)
	if err != nil {
		return nil, err
	}

	urlpaths := []*UrlPath{}
	for _, signurl := range signurls.Urls {
		surl, _ := secure.DecodeBase64(signurl.Url)
		path, _ := secure.DecodeBase64(signurl.Path) // Format as '/bucket/file/paths'
		if surl == "" || path == "" {
			return nil, invar.ErrUnsupportFormat
		}
		urlpaths = append(urlpaths, &UrlPath{surl, path})
	}
	return urlpaths, nil
}

// Obtain a single pre-signed url for upload file to minio storage, and keep the origin name.
func (stub *GrpcStub) NamedUrl(res, filepath string, addition ...string) (*UrlPath, error) {
	if stub.Wss == nil {
		return nil, invar.ErrInvalidClient
	} else if res == "" || filepath == "" {
		return nil, invar.ErrInvalidParams
	}

	fname := path.Base(filepath)
	filesuff := strings.Split(fname, ".")

	add := utils.GetVariable(addition, "").(string)
	param := &wss.FName{Res: res, Add: add, Name: filesuff[0], Suffix: path.Ext(fname)}
	oriurl, err := stub.Wss.OriginalUrl(context.Background(), param)
	if err != nil {
		return nil, err
	}

	surl, _ := secure.DecodeBase64(oriurl.Url)
	path, _ := secure.DecodeBase64(oriurl.Path) // Format as '/bucket/file/paths'
	if surl == "" || path == "" {
		return nil, invar.ErrUnsupportFormat
	}
	return &UrlPath{surl, path}, nil
}

// Obtain multiple pre-signed urls for upload files to minio storage, and keep the origin names.
func (stub *GrpcStub) NamedUrls(res string, files []*wss.NSuffix, addition ...string) ([]*UrlPath, error) {
	if stub.Wss == nil {
		return nil, invar.ErrInvalidClient
	} else if res == "" || len(files) == 0 {
		return nil, invar.ErrInvalidParams
	}

	add := utils.GetVariable(addition, "").(string)
	param := &wss.FNames{Res: res, Add: add, Files: files}
	oriurls, err := stub.Wss.OriginalUrls(context.Background(), param)
	if err != nil {
		return nil, err
	}

	urlpaths := []*UrlPath{}
	for _, oriurl := range oriurls.Urls {
		surl, _ := secure.DecodeBase64(oriurl.Url)
		path, _ := secure.DecodeBase64(oriurl.Path) // Format as '/bucket/file/paths'
		if surl == "" || path == "" {
			return nil, invar.ErrUnsupportFormat
		}
		urlpaths = append(urlpaths, &UrlPath{surl, path})
	}
	return urlpaths, nil
}

// Upload local file to minio storage, then return bucket relative path.
func (stub *GrpcStub) LocalUpload(filepath, res string, delete ...bool) (string, error) {
	if stub.Wss == nil {
		return "", invar.ErrInvalidClient
	} else if filepath == "" || res == "" {
		return "", invar.ErrInvalidParams
	}

	filename := path.Base(filepath)
	urlpath, err := stub.SignUrl(res, path.Ext(filename))
	if err != nil {
		return "", err
	}

	// Read local file data and post to minio storage
	if buff, err := os.ReadFile(filepath); err != nil {
		return "", err
	} else {
		req, err := http.NewRequest(http.MethodPut, urlpath.Url, bytes.NewBuffer(buff))
		if err != nil {
			return "", err
		}

		// Upload file datas and check response status
		if resp, err := http.DefaultClient.Do(req); err != nil {
			return "", err
		} else if resp.StatusCode != invar.StatusOK {
			return "", invar.ErrInvalidData
		}

		// Delete local file when upload success
		if len(delete) > 0 && delete[0] {
			if err := os.Remove(filepath); err != nil {
				return "", err
			}
		}
	}
	return urlpath.Path, nil
}

// Upload local file to minio storage with original name, then return bucket relative path.
func (stub *GrpcStub) NamedUpload(filepath string, res string, delete ...bool) (string, error) {
	if stub.Wss == nil {
		return "", invar.ErrInvalidClient
	} else if filepath == "" || res == "" {
		return "", invar.ErrInvalidParams
	}

	urlpath, err := stub.NamedUrl(res, filepath)
	if err != nil {
		return "", err
	}

	if buff, err := os.ReadFile(filepath); err != nil {
		return "", err
	} else {
		req, err := http.NewRequest(http.MethodPut, urlpath.Url, bytes.NewBuffer(buff))
		if err != nil {
			return "", err
		}

		if resp, err := http.DefaultClient.Do(req); err != nil {
			return "", err
		} else if resp.StatusCode != invar.StatusOK {
			return "", invar.ErrInvalidData
		}

		// Delete local file when upload success
		if len(delete) > 0 && delete[0] {
			if err := os.Remove(filepath); err != nil {
				return "", err
			}
		}
	}
	return urlpath.Path, nil
}

// Set bucket indicated files keeping save status forever.
func (stub *GrpcStub) MarkSave(bucket string, paths ...string) error {
	if stub.Wss == nil {
		return invar.ErrInvalidClient
	} else if bucket == "" || len(paths) == 0 {
		return invar.ErrInvalidParams
	}

	param := &wss.Tag{Bucket: bucket, Paths: paths, Status: "on"}
	if _, err := stub.Wss.SetFileLife(context.Background(), param); err != nil {
		return err
	}
	return nil
}

// Delete given file from minio storage server, the files param must
// contain less one file relative path in minio storage bucket.
func (stub *GrpcStub) Delete(bucket string, files ...string) error {
	if stub.Wss == nil {
		return invar.ErrInvalidClient
	} else if bucket == "" || len(files) == 0 {
		return invar.ErrInvalidParams
	}

	param := &wss.Files{Bucket: bucket, Files: files}
	if _, err := stub.Wss.DeleteFiles(context.Background(), param); err != nil {
		return err
	}
	return nil
}

// Filter out current keeping files and delete the remained old files,
// and the input params only support string or []string.
func (stub *GrpcStub) DiffDelete(bucket string, currs, olds any) error {
	if stub.Wss == nil {
		return invar.ErrInvalidClient
	} else if bucket == "" {
		return invar.ErrInvalidParams
	}

	diffs := []string{}
	switch oldfiles := olds.(type) {
	case []string:
		// Diff the old files to delete
		currfiles := currs.([]string)
		for _, old := range oldfiles {
			iskeeping := false
			for _, curr := range currfiles {
				if old == curr {
					iskeeping = true
					break
				}
			}
			if !iskeeping {
				diffs = append(diffs, old)
			}
		}

	case string:
		curr := currs.(string)
		if olds != curr {
			diffs = append(diffs, oldfiles)
		}
	}

	if len(diffs) > 0 {
		return stub.Delete(bucket, diffs...)
	}
	return nil
}

// ----------------------------------------
// For Acc GRPC agent
// ----------------------------------------

// Get accounts simple infos and avatars.
func (stub *GrpcStub) GetAvatars(uids []string) ([]*acc.Avatar, error) {
	if stub.Acc == nil {
		return nil, invar.ErrInvalidClient
	} else if len(uids) == 0 {
		return nil, invar.ErrInvalidParams
	}

	param := &acc.UIDS{Uids: uids}
	resp, err := stub.Acc.GetAvatars(context.Background(), param)
	if err != nil {
		return nil, err
	}
	return resp.Profs, nil
}
