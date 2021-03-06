package wrap

import (
  "../../../backend"
	"../../../sdlex"

	"math/rand"
	"time"
)

type WildfireWrap struct {
    particles *Particles
}

func MakeWildfireWrap() *WildfireWrap {
	return &WildfireWrap{ particles: &Particles{ randomizer: rand.New(rand.NewSource(time.Now().UTC().UnixNano())) } }
}

func (wildfireWrap *WildfireWrap) Destroy() {}
func (wildfireWrap *WildfireWrap) Initialize(sdlWrap *sdlex.SdlWrap, handle *backend.Handle) error { return nil }
func (wildfireWrap *WildfireWrap) IsReady() bool { return true }

func (wildfireWrap *WildfireWrap) Render(sdlWrap *sdlex.SdlWrap, handle *backend.Handle) error { 
	return wildfireWrap.RenderParticles(sdlWrap)
}

func (wildfireWrap *WildfireWrap) Particles() *Particles {
	return wildfireWrap.particles
}