package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitRatio_Value_Success(t *testing.T) {
	sr := SplitRatio{
		Payer:    60,
		Receiver: 40,
	}

	value, err := sr.Value()
	assert.NoError(t, err)
	assert.NotNil(t, value)

	// Verify the JSON output
	expected := `{"payer":60,"receiver":40}`
	assert.JSONEq(t, expected, string(value.([]byte)))
}

func TestSplitRatio_Scan_Success(t *testing.T) {
	var sr SplitRatio
	jsonData := []byte(`{"payer":70,"receiver":30}`)

	err := sr.Scan(jsonData)
	assert.NoError(t, err)
	assert.Equal(t, 70, sr.Payer)
	assert.Equal(t, 30, sr.Receiver)
}

func TestSplitRatio_Scan_InvalidType(t *testing.T) {
	var sr SplitRatio

	err := sr.Scan("invalid string")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "type assertion to []byte failed")
}

func TestSplitRatio_Scan_InvalidJSON(t *testing.T) {
	var sr SplitRatio
	invalidJSON := []byte(`{"payer":70,"receiver":}`)

	err := sr.Scan(invalidJSON)
	assert.Error(t, err)
}

func TestSplitRatio_RoundTrip(t *testing.T) {
	original := SplitRatio{
		Payer:    80,
		Receiver: 20,
	}

	// Convert to value (serialize)
	value, err := original.Value()
	assert.NoError(t, err)

	// Convert back from value (deserialize)
	var restored SplitRatio
	err = restored.Scan(value)
	assert.NoError(t, err)

	// Verify they are equal
	assert.Equal(t, original.Payer, restored.Payer)
	assert.Equal(t, original.Receiver, restored.Receiver)
}

func TestSplitRatio_Scan_NilValue(t *testing.T) {
	var sr SplitRatio

	err := sr.Scan(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "type assertion to []byte failed")
}

func TestSplitRatio_Value_ZeroValues(t *testing.T) {
	sr := SplitRatio{
		Payer:    0,
		Receiver: 0,
	}

	value, err := sr.Value()
	assert.NoError(t, err)
	assert.NotNil(t, value)

	expected := `{"payer":0,"receiver":0}`
	assert.JSONEq(t, expected, string(value.([]byte)))
}
