package kafka

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.SetTagName("binding")
}

// BindMessage unmarshals JSON data into dest and validates it using "binding" struct tags.
// This is the Kafka consumer equivalent of Gin's c.ShouldBindJSON.
func BindMessage[T any](data json.RawMessage, dest *T) error {
	if err := json.Unmarshal(data, dest); err != nil {
		return err
	}

	return validate.Struct(dest)
}
