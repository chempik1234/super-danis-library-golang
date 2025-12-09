package types

import (
	"errors"
)

var ErrEmptyText = errors.New("empty text")

// NotEmptyText is a value type for any text len>0
type NotEmptyText string

// NewNotEmptyText creates a new NotEmptyText from giving text
func NewNotEmptyText(text string) (NotEmptyText, error) {
	if len(text) > 0 {
		return NotEmptyText(text), nil
	}
	return NotEmptyText(""), ErrEmptyText
}

// String returns value of AnyText of type string
func (d NotEmptyText) String() string {
	return string(d)
}
