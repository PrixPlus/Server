package models

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/prixplus/server/errs"
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

func (t Token) String() string {
	s, err := json.Marshal(t)
	if err != nil { // Just log the error
		errs.LogError(errors.Wrap(err, "encoding json"))
	}

	return string(s)
}
