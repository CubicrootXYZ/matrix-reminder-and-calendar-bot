package format_test

import (
	"testing"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/format"
	"github.com/stretchr/testify/assert"
)

func Test_GetSuffixIntOnSuccess(t *testing.T) {
	result, err := format.GetSuffixInt("Well, this is a very long sentencen: please add some more characters and an integer like 203")
	assert.NoError(t, err)
	assert.Equal(t, 203, result)
}

func Test_GetSuffixIntOnFailure(t *testing.T) {
	t.Run("missing number", func(t *testing.T) {
		result, err := format.GetSuffixInt("Well, this is a very long sentencen: please add some more characters and an integer like two hundred")
		assert.Error(t, err)
		assert.Equal(t, 0, result)
	})

	t.Run("empty string", func(t *testing.T) {
		result, err := format.GetSuffixInt("")
		assert.Error(t, err)
		assert.Equal(t, 0, result)
	})
}
