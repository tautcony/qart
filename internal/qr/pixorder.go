package qr

type Pixorder struct {
	Off      int
	Priority int
}

type byPriority []Pixorder

func (x byPriority) Len() int           { return len(x) }
func (x byPriority) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x byPriority) Less(i, j int) bool { return x[i].Priority > x[j].Priority }
