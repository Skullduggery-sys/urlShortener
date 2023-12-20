package linkShortening

type Hasher interface {
	Hash() (string, error)
}
