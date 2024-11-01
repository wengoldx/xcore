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
	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/wengoldx/xcore/logger"
)

// Elasticsearch client
type ESClient struct {
	Conn *es.Client
}

// Elastic client singleton, setup when DID_ES_AGENTS config received or changed.
var esc *ESClient

// Object logger with [ESC] perfix for elastic module
var esclog = logger.NewLogger("ESC")

// Create a new elasticsearch client.
func NewEsClient(address []string, user, pwd, cfp string) error {
	cfg := es.Config{
		Addresses:              address, // A list of Elasticsearch nodes to use.
		Username:               user,    // Username for HTTP Basic Authentication.
		Password:               pwd,
		CertificateFingerprint: cfp,
	}

	conn, err := es.NewClient(cfg)
	if err != nil {
		return err
	} else if _, err := conn.Info(); err != nil {
		return err
	}

	esc = &ESClient{conn}
	return nil
}

// Return elastic singleton instance,  it may not connect before setuped.
func GetEs() *ESClient {
	if esc == nil {
		return &ESClient{} // Ensure return not nil instance
	}
	return esc
}
