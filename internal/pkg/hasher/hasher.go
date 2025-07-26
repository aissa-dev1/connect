package hasher

type Strategy interface {
	Hash(password string) (string, error)
	Compare(password string, hash string) bool
}

type Hasher struct {
	strategy Strategy
}

func NewHasher(s Strategy) *Hasher {
	return &Hasher{strategy: s}
}

func (h *Hasher) Hash(password string) (string, error) {
	return h.strategy.Hash(password)
}

func (h *Hasher) Compare(password string, hash string) bool {
	return h.strategy.Compare(password, hash)
}
