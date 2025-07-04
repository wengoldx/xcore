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
	mea "github.com/wengoldx/xcore/wrpc/measure/proto"
	wss "github.com/wengoldx/xcore/wrpc/webss/proto"
	chat "github.com/wengoldx/xcore/wrpc/wgchat/proto"
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

	add := utils.Variable(addition, "")
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

	add := utils.Variable(addition, "")
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

	add := utils.Variable(addition, "")
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

	add := utils.Variable(addition, "")
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
		if utils.Variable(delete, false) {
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
		if utils.Variable(delete, false) {
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
	_, err := stub.Wss.SetFileLife(context.Background(), param)
	return err

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
	_, err := stub.Wss.DeleteFiles(context.Background(), param)
	return err
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

// Get file uploaded infos of minio storage server
func (stub *GrpcStub) GetFileInfo(filepath string) (*wss.Info, error) {
	if stub.Wss == nil {
		return nil, invar.ErrInvalidClient
	} else if filepath == "" {
		return nil, invar.ErrInvalidParams
	}

	param := &wss.File{Path: filepath}
	return stub.Wss.GetFileInfo(context.Background(), param)
}

// ----------------------------------------
// For Acc GRPC agent
// ----------------------------------------

// Get account request token by given uuid.
func (stub *GrpcStub) GetToken(uuid string) (string, error) {
	if stub.Acc == nil {
		return "", invar.ErrInvalidClient
	} else if uuid == "" {
		return "", invar.ErrInvalidParams
	}

	param := &acc.UUID{Uuid: uuid}
	token, err := stub.Acc.GetToken(context.Background(), param)
	if err != nil {
		return "", err
	}
	return token.Token, nil
}

// Get account profiles by given uuid.
func (stub *GrpcStub) GetProfile(uuid string) (*acc.Profile, error) {
	if stub.Acc == nil {
		return nil, invar.ErrInvalidClient
	} else if uuid == "" {
		return nil, invar.ErrInvalidParams
	}

	param := &acc.UUID{Uuid: uuid}
	return stub.Acc.GetProfile(context.Background(), param)
}

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

// Search accounts simple infos and avatars from given uuids and keyword.
func (stub *GrpcStub) SearchAvatars(uids []string, filter bool, key string) (*acc.Avatars, error) {
	if stub.Acc == nil {
		return nil, invar.ErrInvalidClient
	} else if len(uids) == 0 || (filter && key == "") {
		return nil, invar.ErrInvalidParams
	}

	param := &acc.SKeys{Uids: uids, Filterid: filter, Keyword: key}
	return stub.Acc.SearchAvatars(context.Background(), param)
}

// Send custom mail.
func (stub *GrpcStub) SendMail(mail, content, name, phone string, buy int32, links []string) error {
	if stub.Acc == nil {
		return invar.ErrInvalidClient
	} else if mail == "" || content == "" || (buy == 1 && name == "") {
		return invar.ErrInvalidParams
	}

	param := &acc.SugMail{
		Email: mail, Content: content, Name: name, Phone: phone,
		Isbuy: buy, Links: strings.Join(links, ","),
	}
	_, err := stub.Acc.SendCustomMail(context.Background(), param)
	return err
}

// ----------------------------------------
// For Mea GRPC agent
// ----------------------------------------

// Get measured body detail by request id.
func (stub *GrpcStub) GetBody(reqid string) (*mea.BodyDetail, error) {
	if stub.Mea == nil {
		return nil, invar.ErrInvalidClient
	} else if reqid == "" {
		return nil, invar.ErrInvalidParams
	}

	param := &mea.ReqID{Reqid: reqid}
	return stub.Mea.GetBody(context.Background(), param)
}

// Get measured bodys detail by request ids.
func (stub *GrpcStub) GetBodys(reqids []string) ([]*mea.BodyBase, error) {
	if stub.Mea == nil {
		return nil, invar.ErrInvalidClient
	} else if len(reqids) == 0 {
		return nil, invar.ErrInvalidParams
	}

	param := &mea.ReqIDs{Reqids: reqids}
	bodys, err := stub.Mea.GetBodys(context.Background(), param)
	if err != nil {
		return nil, err
	}
	return bodys.Body, nil
}

// Post a requst to start measure body by given datas.
func (stub *GrpcStub) Measure(sex, h, w, bust, waist, hipline, wrist int64, front, side, notifier string) (string, error) {
	if stub.Mea == nil {
		return "", invar.ErrInvalidClient
	} else if h <= 0 || w <= 0 || (sex != 1 /* male */ && sex != 2 /* female */) {
		return "", invar.ErrInvalidParams
	}

	param := &mea.BodyComplex{
		Sex: sex, Height: h, Weight: w, Bust: bust, Waist: waist, Hipline: hipline,
		Wrist: wrist, Fronturl: front, Sideurl: side, Notifyurl: notifier,
	}
	reqid, err := stub.Mea.Measure(context.Background(), param)
	if err != nil {
		return "", err
	}
	return reqid.Reqid, nil
}

// Post a requst to measure agine for exist body.
func (stub *GrpcStub) Remeasure(sex, h, w, bust, waist, hipline, wrist int64, reqid, front, side, notifier string) error {
	if stub.Mea == nil {
		return invar.ErrInvalidClient
	} else if reqid == "" || h <= 0 || w <= 0 || (sex != 1 /* male */ && sex != 2 /* female */) {
		return invar.ErrInvalidParams
	}

	param := &mea.UpComplex{
		Reqid: reqid, Sex: sex, Height: h, Weight: w, Bust: bust, Waist: waist,
		Hipline: hipline, Wrist: wrist, Fronturl: front, Sideurl: side, Notifyurl: notifier,
	}
	_, err := stub.Mea.Remeasure(context.Background(), param)
	return err
}

// Post a request to start capture body model.
func (stub *GrpcStub) BodyShot(reqid string) error {
	if stub.Mea == nil {
		return invar.ErrInvalidClient
	} else if reqid == "" {
		return invar.ErrInvalidParams
	}

	param := &mea.ReqID{Reqid: reqid}
	_, err := stub.Mea.BodyShot(context.Background(), param)
	return err
}

// Delete exist body by request id.
func (stub *GrpcStub) DelBody(reqid string) error {
	if stub.Mea == nil {
		return invar.ErrInvalidClient
	} else if reqid == "" {
		return invar.ErrInvalidParams
	}

	param := &mea.ReqID{Reqid: reqid}
	_, err := stub.Mea.DelBody(context.Background(), param)
	return err
}

// ----------------------------------------
// For Chat GRPC agent
// ----------------------------------------

// Add a new staff to company.
func (stub *GrpcStub) AddStaff(uuid, name, headurl, brand, client, old string) error {
	if stub.Chat == nil {
		return invar.ErrInvalidClient
	} else if uuid == "" || name == "" {
		return invar.ErrInvalidParams
	}

	param := &chat.Staff{
		Uuid: uuid, Nickname: name, Headurl: headurl,
		Company: brand, Client: client, Old: old,
	}
	_, err := stub.Chat.AddStaff(context.Background(), param)
	return err
}

// Update staff status.
func (stub *GrpcStub) UpdateStaff(brand, client string, status bool) error {
	if stub.Chat == nil {
		return invar.ErrInvalidClient
	} else if brand == "" || client == "" {
		return invar.ErrInvalidParams
	}

	param := &chat.Status{Company: brand, Client: client, Status: status}
	_, err := stub.Chat.UpdateStatus(context.Background(), param)
	return err
}
