package metadata

// Metadata interface
type Metadata interface {
	ToXML() ([]byte, error)
	Save(filename string) error
}

// MovieNfo interface
type MovieNfo interface {
	Metadata
	SetPoster(filename string)
}
