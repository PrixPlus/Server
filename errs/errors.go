package errs

import "errors"

var (
	ElemNotFound  = errors.New("Element not found")
	ElemNotUnique = errors.New("Element not unique")
)

var Status = map[error]int{
	ElemNotFound:  471,
	ElemNotUnique: 472,
}
