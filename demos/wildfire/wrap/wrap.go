package wrap

import (
  "../../../backend"
	"../../../sdlex"
)

type WildfireWrap struct {}

func MakeWildfireWrap() *WildfireWrap {
	return &WildfireWrap{}
}

func (wildfireWrap *WildfireWrap) Destroy() {}
func (wildfireWrap *WildfireWrap) Initialize(sclWrap *sdlex.SdlWrap, handle *backend.Handle) error { return nil }
func (wildfireWrap *WildfireWrap) IsReady() bool { return true }

func (wildfireWrap *WildfireWrap) Render(sdlWrap *sdlex.SdlWrap, handle *backend.Handle) error { 
	var err error 

	err = wildfireWrap.RenderParticles(sdlWrap, handle)

	return err 
}