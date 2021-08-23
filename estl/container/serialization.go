package container

type JSONSerializer interface {
	ToJSON() ([]byte, error)
}

type JSONDeserializer interface {
	FromJSON([]byte) error
}
