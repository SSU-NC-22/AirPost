package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type KafkaData struct {
	NodeType  string	`json:"node_type"`
	NodeID    int       `json:"node_id"`
	Values    []float64 `json:"values"`
	Timestamp time.Time `json:"timestamp"`
}

type Document struct {
	Index string
	Doc   interface{}
}

func (d *Document) String() string {
	doc, err := json.Marshal(d.Doc)
	if err != nil {
		return ""
	}
	h := fmt.Sprintf("{\"index\":{\"_index\":\"%s\"}}\n", d.Index)
	return strings.Join([]string{h, string(doc), "\n"}, "")
}
