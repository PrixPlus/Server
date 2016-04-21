package models

import (
	"time"
)

type Token struct {
	Raw       string                 `json:"raw"`
	Expire    time.Time              `json:"expire"`
	Method    string                 `json:"-"`
	Header    map[string]interface{} `json:"-"`
	Claims    map[string]interface{} `json:"-"`
	Signature string                 `json:"-"`
	Valid     bool                   `json:"-"`
}
