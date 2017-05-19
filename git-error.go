/*
Package gitwrap implements a simple library for GIT interactions.
*/
package gitwrap

import "errors"

// Error codes returned by failures to parse an expression.
var (
	ErrNoRepo              = errors.New("GIT: No repository found in working directory")
	ErrNoGit               = errors.New("GIT: Could not find GIT library")
	ErrTagCreateNoPrevious = errors.New("GIT: No previous tag present, we need to compare")
	ErrTagCreateDirty      = errors.New("GIT: Repository is still dirty, please add and commit changes")
	ErrTagCreateNoChanges  = errors.New("GIT: No changes since latest tag")
	ErrTagAlreadyExists    = errors.New("GIT: Tag already exists")
	ErrTagCreating         = errors.New("GIT: Error creating tag")
	ErrTagEmpty            = errors.New("GIT: Did not find any tags matching your query")
	ErrDiff                = errors.New("GIT: Error while verifying git diff")
	ErrNoRemote            = errors.New("GIT: Did not find remote stream")
)
