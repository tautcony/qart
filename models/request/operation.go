package request

import (
	"github.com/creasty/defaults"
	"math/rand"
	"rsc.io/qr/coding"
	"strconv"
)

type Operation struct {
	Image        string         `json:"image" default:"default"`
	Dx           int            `json:"dx" default:"4"`
	Dy           int            `json:"dy" default:"4"`
	Size         int            `json:"size" default:"0"`
	URL          string         `json:"url" default:"https://example.com"`
	Version      coding.Version `json:"version" default:"6"`
	Mask         coding.Mask    `json:"mask" default:"2"`
	RandControl  bool           `json:"randcontrol" default:"false"`
	Dither       bool           `json:"dither" default:"false"`
	OnlyDataBits bool           `json:"onlydatabits" default:"false"`
	SaveControl  bool           `json:"savecontrol" default:"false"`
	Seed         string         `json:"seed"`
	Scale        int            `json:"scale" default:"8"`
	Rotation     int            `json:"rotation" default:"0"` // range in [0,3]
}

func (op *Operation) SetDefaults() {
	if defaults.CanUpdate(op.Seed) {
		op.Seed = strconv.FormatInt(rand.Int63(), 10)
	}
}

func (op *Operation) GetVersion() coding.Version {
	if op.Version < 1 {
		return 1
	}
	if op.Version > 40 {
		return 40
	}
	return op.Version
}

func (op *Operation) GetMask() coding.Mask {
	if op.Mask < 0 {
		return 0
	}
	if op.Mask > 7 {
		return 7
	}
	return op.Mask
}

func (op *Operation) GetRotation() int {
	if op.Rotation < 0 {
		return 0
	}
	if op.Rotation > 3 {
		return 3
	}
	return op.Rotation
}

func (op *Operation) GetScale() int {
	if op.Version >= 12 && op.Scale >= 4 {
		return op.Scale / 2
	}
	return op.Scale
}

func (op *Operation) GetSeed() int64 {
	seed, err := strconv.ParseInt(op.Seed, 10, 64)
	if err != nil {
		return rand.Int63()
	}
	return seed
}

func NewOperation() (*Operation, error) {
	operation := &Operation{}
	var err error
	if err = defaults.Set(operation); err != nil {
		return nil, err
	}
	return operation, nil
}
