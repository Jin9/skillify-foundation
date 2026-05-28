package codec

import (
	"testing"
)

func TestJSONCoder_Marshal(t *testing.T) {
	coder := NewJSONCoder()

	t.Run("valid input", func(t *testing.T) {
		input := map[string]string{"key": "value"}
		expected := `{"key":"value"}` + "\n"

		data, err := coder.Marshal(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(data) != expected {
			t.Errorf("expected %q, got %q", expected, string(data))
		}
	})

	t.Run("error encoding", func(t *testing.T) {
		// make chan which cannot be marshalled to JSON
		input := make(chan int)
		_, err := coder.Marshal(input)
		if err == nil {
			t.Error("expected error for channel input")
		}
	})
}

func TestJSONCoder_Unmarshal(t *testing.T) {
	coder := NewJSONCoder()

	t.Run("valid json", func(t *testing.T) {
		data := []byte(`{"key":"value"}`)
		var output map[string]string

		err := coder.Unmarshal(data, &output)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if output["key"] != "value" {
			t.Errorf("expected 'value', got %q", output["key"])
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		data := []byte(`{"key":"value"`)
		var output map[string]string

		err := coder.Unmarshal(data, &output)
		if err == nil {
			t.Error("expected error for invalid json")
		}
	})
}
