package nebulaorm

import (
	"errors"
	"github.com/haysons/nebulaorm/clause"
	"github.com/haysons/nebulaorm/resolver"
)

var (
	// ErrInvalidValue usually because a variable is passed in that doesn't match the expected type
	ErrInvalidValue = errors.New("invalid value")

	// ErrRecordNotFound methods that query only a single record return this error if they fail to get the record. eg: Take
	ErrRecordNotFound = errors.New("record not found")

	// ErrValueCannotSet usually because the variable is not passed in as a pointer and cannot be assigned a value
	ErrValueCannotSet = resolver.ErrValueCannotSet

	// ErrInvalidClauseParams usually because the arguments to the build clause are anomalous, causing the build to fail
	ErrInvalidClauseParams = clause.ErrInvalidClauseParams
)
