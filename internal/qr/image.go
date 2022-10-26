package qr

import (
	"bytes"
	"fmt"
	"github.com/tautcony/qart/internal/utils"
	"log"
	"math/rand"
	"rsc.io/qr"
	"rsc.io/qr/coding"
	"rsc.io/qr/gf256"
	"sort"
)

type Image struct {
	Name     string
	Target   [][]byte
	Dx       int
	Dy       int
	URL      string
	Version  coding.Version
	Mask     coding.Mask
	Level    coding.Level
	Scale    int
	Rotation int
	Size     int

	// RandControl says to pick the pixels randomly.
	RandControl bool
	Seed        int64

	// Dither says to dither instead of using threshold pixel layout.
	Dither bool

	// OnlyDataBits says to use only data bits, not check bits.
	OnlyDataBits bool

	// Code is the final QR code.
	Code *qr.Code

	// Control is a PNG showing the pixels that we controlled.
	// Pixels we don't control are grayed out.
	SaveControl bool
	Control     []byte
}

func (m *Image) target(x, y int) (targ byte, contrast int) {
	tx := x + m.Dx
	ty := y + m.Dy
	if ty < 0 || ty >= len(m.Target) || tx < 0 || tx >= len(m.Target[ty]) {
		return 255, -1
	}

	v0 := m.Target[ty][tx]
	if v0 < 0 {
		return 255, -1
	}
	targ = v0

	n := 0
	sum := 0
	sumsq := 0
	const del = 5
	for dy := -del; dy <= del; dy++ {
		for dx := -del; dx <= del; dx++ {
			if 0 <= ty+dy && ty+dy < len(m.Target) && 0 <= tx+dx && tx+dx < len(m.Target[ty+dy]) {
				v := int(m.Target[ty+dy][tx+dx])
				sum += v
				sumsq += v * v
				n++
			}
		}
	}

	avg := sum / n
	contrast = sumsq/n - avg*avg
	return
}

func (m *Image) rotate(p *coding.Plan, rot int) {
	utils.RotatePixel(p.Pixel, rot)
}

func (m *Image) Encode() error {
	p, err := coding.NewPlan(m.Version, m.Level, m.Mask)
	if err != nil {
		return err
	}

	m.rotate(p, m.Rotation)

	randNumber := rand.New(rand.NewSource(m.Seed))

	// QR parameters.
	nd := p.DataBytes / p.Blocks
	nc := p.CheckBytes / p.Blocks
	extra := p.DataBytes - nd*p.Blocks
	rs := gf256.NewRSEncoder(coding.Field, nc)

	// Build information about pixels, indexed by data/check bit number.
	pixByOff := make([]Pixinfo, (p.DataBytes+p.CheckBytes)*8)
	expect := make([][]bool, len(p.Pixel))
	for y, row := range p.Pixel {
		expect[y] = make([]bool, len(row))
		for x, pix := range row {
			targ, contrast := m.target(x, y)
			if m.RandControl && contrast >= 0 {
				contrast = randNumber.Intn(128) + 64*((x+y)%2) + 64*((x+y)%3%2)
			}
			expect[y][x] = pix&coding.Black != 0
			if r := pix.Role(); r == coding.Data || r == coding.Check {
				pixByOff[pix.Offset()] = Pixinfo{X: x, Y: y, Pix: pix, Targ: targ, Contrast: contrast}
			}
		}
	}

Again:
	// Count fixed initial data bits, prepare template URL.
	url := m.URL + "#"
	var b coding.Bits
	coding.String(url).Encode(&b, p.Version)
	coding.Num("").Encode(&b, p.Version)
	bbit := b.Bits()
	dbit := p.DataBytes*8 - bbit
	if dbit < 0 {
		return fmt.Errorf("cannot encode URL into available bits")
	}
	num := make([]byte, dbit/10*3)
	for i := range num {
		num[i] = '0'
	}
	b.Pad(dbit)
	b.Reset()
	coding.String(url).Encode(&b, p.Version)
	coding.Num(num).Encode(&b, p.Version)
	b.AddCheckBytes(p.Version, p.Level)
	data := b.Bytes()

	doff := 0 // data offset
	coff := 0 // checksum offset
	mbit := bbit + dbit/10*10

	// Choose pixels.
	bitblocks := make([]*BitBlock, p.Blocks)
	for blocknum := 0; blocknum < p.Blocks; blocknum++ {
		if blocknum == p.Blocks-extra {
			nd++
		}

		bdata := data[doff/8 : doff/8+nd]
		cdata := data[p.DataBytes+coff/8 : p.DataBytes+coff/8+nc]
		bb := newBlock(nd, nc, rs, bdata, cdata)
		bitblocks[blocknum] = bb

		// Determine which bits in this block we can try to edit.
		lo, hi := 0, nd*8
		if lo < bbit-doff {
			lo = bbit - doff
			if lo > hi {
				lo = hi
			}
		}
		if hi > mbit-doff {
			hi = mbit - doff
			if hi < lo {
				hi = lo
			}
		}

		// Preserve [0, lo) and [hi, nd*8).
		for i := 0; i < lo; i++ {
			if !bb.canSet(uint(i), (bdata[i/8]>>uint(7-i&7))&1) {
				return fmt.Errorf("cannot preserve required bits")
			}
		}
		for i := hi; i < nd*8; i++ {
			if !bb.canSet(uint(i), (bdata[i/8]>>uint(7-i&7))&1) {
				return fmt.Errorf("cannot preserve required bits")
			}
		}

		// Can edit [lo, hi) and checksum bits to hit target.
		// Determine which ones to try first.
		order := make([]Pixorder, (hi-lo)+nc*8)
		for i := lo; i < hi; i++ {
			order[i-lo].Off = doff + i
		}
		for i := 0; i < nc*8; i++ {
			order[hi-lo+i].Off = p.DataBytes*8 + coff + i
		}
		if m.OnlyDataBits {
			order = order[:hi-lo]
		}
		for i := range order {
			po := &order[i]
			po.Priority = pixByOff[po.Off].Contrast<<8 | randNumber.Intn(256)
		}
		sort.Sort(byPriority(order))

		const mark = false
		for i := range order {
			po := &order[i]
			pinfo := &pixByOff[po.Off]
			bval := pinfo.Targ
			if bval < 128 {
				bval = 1
			} else {
				bval = 0
			}
			pix := pinfo.Pix
			if pix&coding.Invert != 0 {
				bval ^= 1
			}
			if pinfo.HardZero {
				bval = 0
			}

			var bi int
			if pix.Role() == coding.Data {
				bi = po.Off - doff
			} else {
				bi = po.Off - p.DataBytes*8 - coff + nd*8
			}
			if bb.canSet(uint(bi), bval) {
				pinfo.Block = bb
				pinfo.Bit = uint(bi)
				if mark {
					p.Pixel[pinfo.Y][pinfo.X] = coding.Black
				}
			} else {
				if pinfo.HardZero {
					panic("hard zero")
				}
				if mark {
					p.Pixel[pinfo.Y][pinfo.X] = 0
				}
			}
		}
		bb.copyOut()

		const cheat = false
		for i := 0; i < nd*8; i++ {
			pinfo := &pixByOff[doff+i]
			pix := p.Pixel[pinfo.Y][pinfo.X]
			if bb.B[i/8]&(1<<uint(7-i&7)) != 0 {
				pix ^= coding.Black
			}
			expect[pinfo.Y][pinfo.X] = pix&coding.Black != 0
			if cheat {
				p.Pixel[pinfo.Y][pinfo.X] = pix & coding.Black
			}
		}
		for i := 0; i < nc*8; i++ {
			pinfo := &pixByOff[p.DataBytes*8+coff+i]
			pix := p.Pixel[pinfo.Y][pinfo.X]
			if bb.B[nd+i/8]&(1<<uint(7-i&7)) != 0 {
				pix ^= coding.Black
			}
			expect[pinfo.Y][pinfo.X] = pix&coding.Black != 0
			if cheat {
				p.Pixel[pinfo.Y][pinfo.X] = pix & coding.Black
			}
		}
		doff += nd * 8
		coff += nc * 8
	}

	// Pass over all pixels again, dithering.
	if m.Dither {
		for i := range pixByOff {
			pinfo := &pixByOff[i]
			pinfo.DTarg = int(pinfo.Targ)
		}
		for y, row := range p.Pixel {
			for x, pix := range row {
				if pix.Role() != coding.Data && pix.Role() != coding.Check {
					continue
				}
				pinfo := &pixByOff[pix.Offset()]
				if pinfo.Block == nil {
					// did not choose this pixel
					continue
				}

				pix := pinfo.Pix

				pval := byte(1) // pixel value (black)
				v := 0          // gray value (black)
				targ := pinfo.DTarg
				if targ >= 128 {
					// want white
					pval = 0
					v = 255
				}

				bval := pval // bit value
				if pix&coding.Invert != 0 {
					bval ^= 1
				}
				if pinfo.HardZero && bval != 0 {
					bval ^= 1
					pval ^= 1
					v ^= 255
				}

				// Set pixel value as we want it.
				pinfo.Block.reset(pinfo.Bit, bval)

				_, _ = x, y

				err := targ - v
				if x+1 < len(row) {
					addDither(pixByOff, row[x+1], err*7/16)
				}
				if false && y+1 < len(p.Pixel) {
					if x > 0 {
						addDither(pixByOff, p.Pixel[y+1][x-1], err*3/16)
					}
					addDither(pixByOff, p.Pixel[y+1][x], err*5/16)
					if x+1 < len(row) {
						addDither(pixByOff, p.Pixel[y+1][x+1], err*1/16)
					}
				}
			}
		}

		for _, bb := range bitblocks {
			bb.copyOut()
		}
	}

	noops := 0
	// Copy numbers back out.
	for i := 0; i < dbit/10; i++ {
		// Pull out 10 bits.
		v := 0
		for j := 0; j < 10; j++ {
			bi := uint(bbit + 10*i + j)
			v <<= 1
			v |= int((data[bi/8] >> (7 - bi&7)) & 1)
		}
		// Turn into 3 digits.
		if v >= 1000 {
			// Oops - too many 1 bits.
			// We know the 512, 256, 128, 64, 32 bits are all set.
			// Pick one at random to clear.  This will break some
			// checksum bits, but so be it.
			// log.Println("Oops - too many 1 bits", i, v)
			pinfo := &pixByOff[bbit+10*i+3] // TODO random
			pinfo.Contrast = 1e9 >> 8
			pinfo.HardZero = true
			noops++
		}
		num[i*3+0] = byte(v/100 + '0')
		num[i*3+1] = byte(v/10%10 + '0')
		num[i*3+2] = byte(v%10 + '0')
	}
	if noops > 0 {
		goto Again
	}

	var b1 coding.Bits
	coding.String(url).Encode(&b1, p.Version)
	coding.Num(num).Encode(&b1, p.Version)
	b1.AddCheckBytes(p.Version, p.Level)
	if !bytes.Equal(b.Bytes(), b1.Bytes()) {
		log.Printf("mismatch\n%d %x\n%d %x\n", len(b.Bytes()), b.Bytes(), len(b1.Bytes()), b1.Bytes())
		panic("byte mismatch")
	}

	cc, err := p.Encode(coding.String(url), coding.Num(num))
	if err != nil {
		return err
	}

	if !m.Dither {
		for y, row := range expect {
			for x, pix := range row {
				if cc.Black(x, y) != pix {
					log.Println("mismatch", x, y, p.Pixel[y][x].String())
				}
			}
		}
	}

	m.Code = &qr.Code{Bitmap: cc.Bitmap, Size: cc.Size, Stride: cc.Stride, Scale: m.Scale}

	if m.SaveControl {
		m.Control = utils.PngEncode(utils.MakeImage("", "", 0, cc.Size, 4, m.Scale, func(x, y int) (rgba uint32) {
			pix := p.Pixel[y][x]
			if pix.Role() == coding.Data || pix.Role() == coding.Check {
				pinfo := &pixByOff[pix.Offset()]
				if pinfo.Block != nil {
					if cc.Black(x, y) {
						return 0x000000ff
					}
					return 0xffffffff
				}
			}
			if cc.Black(x, y) {
				return 0x3f3f3fff
			}
			return 0xbfbfbfff
		}))
	}

	return nil
}
