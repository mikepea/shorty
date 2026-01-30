package cuid

import (
	"github.com/nrednav/cuid2"
)

// New generates a new CUID2 identifier (24 characters by default)
func New() string {
	return cuid2.Generate()
}

// IsValid checks if a string is a valid CUID2 format
func IsValid(id string) bool {
	return cuid2.IsCuid(id)
}
