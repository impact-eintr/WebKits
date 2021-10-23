package containers

type JSONSerializer interface {
	ToJSON() ([]byte, error)
}

type JSONDeserializer interface {
	FromJSON([]byte) error
}
