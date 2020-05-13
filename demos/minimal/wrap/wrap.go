package wrap

import (
  "../../../backend"
	"../../../sdlex"
)

type MinimalWrap struct {}

func MakeMinimalWrap() *MinimalWrap {
	return &MinimalWrap{}
}

func (minimalWrap *MinimalWrap) Destroy() {}
func (minimalWrap *MinimalWrap) Initialize(sclWrap *sdlex.SdlWrap, handle *backend.Handle) error { return nil }
func (minimalWrap *MinimalWrap) IsReady() bool { return true }
func (minimalWrap *MinimalWrap) Render(sdlWrap *sdlex.SdlWrap, handle *backend.Handle) error { return nil }