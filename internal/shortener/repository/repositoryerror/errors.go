package repositoryerror

import "errors"

var ErrConflict = errors.New("conflict error")
var ErrNotFound = errors.New("not found error")
