package qr

import (
	"bytes"
	"log"
	"rsc.io/qr/gf256"
)

type BitBlock struct {
	DataBytes  int
	CheckBytes int
	B          []byte
	M          [][]byte
	Tmp        []byte
	RS         *gf256.RSEncoder
	bdata      []byte
	cdata      []byte
}

func newBlock(nd, nc int, rs *gf256.RSEncoder, dat, cdata []byte) *BitBlock {
	b := &BitBlock{
		DataBytes:  nd,
		CheckBytes: nc,
		B:          make([]byte, nd+nc),
		Tmp:        make([]byte, nc),
		RS:         rs,
		bdata:      dat,
		cdata:      cdata,
	}
	copy(b.B, dat)
	rs.ECC(b.B[:nd], b.B[nd:])
	b.check()
	if !bytes.Equal(b.Tmp, cdata) {
		panic("cdata")
	}

	b.M = make([][]byte, nd*8)
	for i := range b.M {
		row := make([]byte, nd+nc)
		b.M[i] = row
		for j := range row {
			row[j] = 0
		}
		row[i/8] = 1 << (7 - uint(i%8))
		rs.ECC(row[:nd], row[nd:])
	}
	return b
}

func (b *BitBlock) check() {
	b.RS.ECC(b.B[:b.DataBytes], b.Tmp)
	if !bytes.Equal(b.B[b.DataBytes:], b.Tmp) {
		log.Printf("ecc mismatch\n%x\n%x\n", b.B[b.DataBytes:], b.Tmp)
		panic("mismatch")
	}
}

func (b *BitBlock) reset(bi uint, bval byte) {
	if (b.B[bi/8]>>(7-bi&7))&1 == bval {
		// already has desired bit
		return
	}
	// rows that have already been set
	m := b.M[len(b.M):cap(b.M)]
	for _, row := range m {
		if row[bi/8]&(1<<(7-bi&7)) != 0 {
			// Found it.
			for j, v := range row {
				b.B[j] ^= v
			}
			return
		}
	}
	panic("reset of unset bit")
}

func (b *BitBlock) canSet(bi uint, bval byte) bool {
	found := false
	m := b.M
	for j, row := range m {
		if row[bi/8]&(1<<(7-bi&7)) == 0 {
			continue
		}
		if !found {
			found = true
			if j != 0 {
				m[0], m[j] = m[j], m[0]
			}
			continue
		}
		for k := range row {
			row[k] ^= m[0][k]
		}
	}
	if !found {
		return false
	}

	targ := m[0]

	// Subtract from saved-away rows too.
	for _, row := range m[len(m):cap(m)] {
		if row[bi/8]&(1<<(7-bi&7)) == 0 {
			continue
		}
		for k := range row {
			row[k] ^= targ[k]
		}
	}

	// Found a row with bit #bi == 1 and cut that bit from all the others.
	// Apply to data and remove from m.
	if (b.B[bi/8]>>(7-bi&7))&1 != bval {
		for j, v := range targ {
			b.B[j] ^= v
		}
	}
	b.check()
	n := len(m) - 1
	m[0], m[n] = m[n], m[0]
	b.M = m[:n]

	for _, row := range b.M {
		if row[bi/8]&(1<<(7-bi&7)) != 0 {
			panic("did not reduce")
		}
	}

	return true
}

func (b *BitBlock) copyOut() {
	b.check()
	copy(b.bdata, b.B[:b.DataBytes])
	copy(b.cdata, b.B[b.DataBytes:])
}
