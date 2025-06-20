package layout

import "errors"

var (
	ErrLayoutNotFound      = errors.New("layout not found")
	ErrInvalidLayoutData   = errors.New("invalid layout data")
	ErrTemplateRender      = errors.New("template render error")
	ErrLayoutAlreadyExists = errors.New("layout with this name already exists")
)
