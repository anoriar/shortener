package repositoryerror

import "errors"

var ErrConflict = errors.New("conflict error")
var NotFound = errors.New("not found error")
