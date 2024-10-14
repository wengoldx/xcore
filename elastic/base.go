// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2022/10/20   jidi           New version
// -------------------------------------------------------------------

package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/wengoldx/xcore/invar"
	"github.com/wengoldx/xcore/logger"
)

func (e *ESClient) SetupIndexs(indexs map[string]string) error {
	if e.Conn == nil || e.Conn.Indices == nil {
		return invar.ErrInvalidClient
	}

	for index, mapping := range indexs {
		if status, err := e.Conn.Indices.Get([]string{index}); err != nil {
			logger.E("Check index:", index, "if exist, err:", err)
			continue
		} else if status.StatusCode == invar.StatusOK {
			continue
		}
		e.CreateIndexMapping(index, mapping)
	}

	logger.I("Setup elastic mappings")
	return nil
}

// Create the new index and setting index mapping, it will return error when index is exist.
//
// The mapping param like :
//
//	mapping = `
//	{
//		"mappings": {
//			"properties": {
//				"title": {
//					"type": "text",
//					"analyzer": "ik_max_word",
//					"search_analyzer": "ik_smart"
//				},
//			}
//		}
//	}`
func (e *ESClient) CreateIndexMapping(index, mapping string) error {
	if e.Conn == nil || e.Conn.Indices == nil {
		return invar.ErrInvalidClient
	}

	createfunc := e.Conn.Indices.Create.WithBody(strings.NewReader(mapping))
	res, err := e.Conn.Indices.Create(index, createfunc)
	if err != nil {
		logger.E("Create index, err:", err)
		return err
	}
	return respError(res)
}

// Update the index mapping, for exmple: add new field or update one.
//
// The mapping param like :
//
//	mapping := `
//	{
//		"properties": {
//			"title": {
//			"type": "text",
//			"analyzer": "ik_max_word",
//			"search_analyzer": "ik_smart"
//			}
//		}
//	}`
func (e *ESClient) UpdateIndexMapping(index []string, mapping string) error {
	if e.Conn == nil || e.Conn.Indices == nil {
		return invar.ErrInvalidClient
	}

	res, err := e.Conn.Indices.PutMapping(index, strings.NewReader(mapping))
	if err != nil {
		logger.E("Update index, err:", err)
		return err
	}
	return respError(res)
}

// Create new doc, it will auto create index mapping when index unexist.
func (e *ESClient) CreateIndexDoc(index string, doc any, docid ...string) error {
	if e.Conn == nil || e.Conn.Indices == nil {
		return invar.ErrInvalidClient
	}

	id := ""
	if len(docid) > 0 {
		id = docid[0]
	}

	body, err := json.Marshal(doc)
	if err != nil {
		logger.E("Marshal index doc, err:", err)
		return err
	}

	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: id,
		Body:       bytes.NewReader(body),
		Refresh:    "wait_for",
	}

	res, err := req.Do(context.Background(), e.Conn)
	if err != nil {
		logger.E("Create index doc, err:", err)
		return err
	}
	return respError(res)
}

// Update the specified fields in the index.
//
// The doc string like :
//
//	doc := `
//	{
//		"doc": {
//			"fields":"value"
//		}
//	}
//	`
func (e *ESClient) UpdateIndexDoc(index, docid, doc string) error {
	if e.Conn == nil {
		return invar.ErrInvalidClient
	}

	updatefunc := e.Conn.Update.WithRefresh("wait_for")
	res, err := e.Conn.Update(index, docid, strings.NewReader(doc), updatefunc)
	if err != nil {
		logger.E("Update index doc, err:", err)
		return err
	}
	return respError(res)
}

// Search doc by query index, and set page, limit. by default page=0 and limit=10
func (e *ESClient) SearchIndex(index, query string, page int, limit ...int) (*Response, error) {
	if e.Conn == nil {
		return nil, invar.ErrInvalidClient
	}

	size := 10
	if len(limit) > 0 {
		size = limit[0]
	}

	res, err := e.Conn.Search(
		e.Conn.Search.WithIndex(index),
		e.Conn.Search.WithSize(size),
		e.Conn.Search.WithFrom(page),
		e.Conn.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		logger.E("Search index, err:", err)
		return nil, err
	}

	defer res.Body.Close()
	return readResponse(res)
}

// Check indexs whether exist
func (e *ESClient) IsExistIndex(index []string) (bool, error) {
	if e.Conn == nil || e.Conn.Indices == nil {
		return false, invar.ErrInvalidClient
	}

	exist, err := e.Conn.Indices.Get(index)
	if err != nil {
		return false, err
	}
	return (exist.StatusCode == invar.StatusOK), nil
}

// ----------------------------------------

// Read search response data or error reason
func readResponse(res *esapi.Response) (*Response, error) {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("Read resp err:" + err.Error())
	}

	if !res.IsError() {
		resp := &Response{}
		if err := json.Unmarshal(body, resp); err != nil {
			return nil, errors.New("Decode resp, err:" + err.Error())
		}
		return resp, nil
	} else {
		resp := &ErrorResp{}
		if err := json.Unmarshal(body, resp); err != nil {
			return nil, errors.New("Decode error, err:" + err.Error())
		}
		return nil, errors.New(resp.ErrorReason.Reason)
	}
}

// Response error if exist
func respError(res *esapi.Response) error {
	defer res.Body.Close()
	if !res.IsError() {
		return nil
	}

	_, err := readResponse(res)
	return err
}
