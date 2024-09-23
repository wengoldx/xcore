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

import "encoding/json"

type Response struct {
	Took         int           `json:"took"`
	TimeOut      bool          `json:"timed_out"`
	Shards       *Shards       `json:"_shards"`
	Hits         *SearchHits   `json:"hits"`
	Aggregations *Aggregations `json:"aggregations"`
}

type Shards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

type SearchHits struct {
	Total    *Total       `json:"total"`
	MaxScore *float64     `json:"max_score"`
	Hits     []*SearchHit `json:"hits"`
}

type Total struct {
	Values   int    `json:"value"`
	Relation string `json:"relation"`
}

type SearchHit struct {
	Index  string          `json:"_index"`
	ID     string          `json:"_id"`
	Score  *float64        `json:"_score,omitempty"`  // computed score
	Source json.RawMessage `json:"_source,omitempty"` // stored document source

}

type Aggregations struct {
	TopN *TopN `json:"topN"`
}

type TopN struct {
	DocCEUB int        `json:"doc_count_error_upper_bound"`
	SumODC  int        `json:"sum_other_doc_count"`
	Buckets []*Buckets `json:"buckets"`
}

type Buckets struct {
	Key      string `json:"key"`
	DocCount string `json:"doc_count"`
}

type ErrorResp struct {
	Status      int     `json:"status"`
	ErrorReason *Reason `json:"error"`
}

type Reason struct {
	RootCause    []*RootCause `json:"root_cause"`
	Type         string       `json:"type"`
	Reason       string       `json:"reason"`
	ResourceType string       `json:"resource.type"`
	ResourceID   string       `json:"resource.id"`
	IndexUUID    string       `json:"index_uuid"`
	Index        string       `json:"index"`
}

type RootCause struct {
	Type         string `json:"type"`
	Reason       string `json:"reason"`
	ResourceType string `json:"resource.type"`
	ResourceID   string `json:"resource.id"`
	IndexUUID    string `json:"index_uuid"`
	Index        string `json:"index"`
}
