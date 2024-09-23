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
	"fmt"

	es "github.com/elastic/go-elasticsearch/v8"
)

type ESClient struct {
	Conn *es.Client
}

func CreateNewClient(address []string, user, pwd, cfp string) (*ESClient, error) {
	cfg := es.Config{
		Addresses:              address, // A list of Elasticsearch nodes to use.
		Username:               user,    // Username for HTTP Basic Authentication.
		Password:               pwd,
		CertificateFingerprint: cfp,
	}

	c := &ESClient{}
	conn, err := es.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("dail with elastic search server err:%v", err)
	}

	// get cluster info
	if _, err := conn.Info(); err != nil {
		return nil, err
	}

	c.Conn = conn
	return c, nil
}
