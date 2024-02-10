package domainerror

import "errors"

// ErrNotValidURL missing godoc.
var ErrNotValidURL = errors.New("not valid URL")

// ErrInternal missing godoc.
var ErrInternal = errors.New("internal error")

// ErrURLExists missing godoc.
var ErrURLExists = errors.New("url already exists")

// ErrURLNotExist missing godoc.
var ErrURLNotExist = errors.New("url does not exists")

// ErrURLDeleted missing godoc.
var ErrURLDeleted = errors.New("url deleted")

var ErrNotValidRequest = errors.New("request not valid")

var ErrUserUnauthorized = errors.New("user not authorized")
