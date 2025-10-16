package hasher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const ITERATIONS = 200

func Test_NewToken(t *testing.T) {
	t.Parallel()
	tokens := make([]Token, ITERATIONS)
	for i := 0; i < ITERATIONS; i++ {
		tokens[i] = GenerateToken()
	}

	// Check all tokens are unique
	for i := 0; i < ITERATIONS; i++ {
		for j := 0; j < ITERATIONS; j++ {
			if i == j {
				continue
			}
			assert.NotEqual(t, tokens[i].Raw, tokens[j].Raw)
			assert.NotEqual(t, tokens[i].Hash, tokens[j].Hash)
		}
	}
}

func Test_HashToken_CheckTokenHash(t *testing.T) {
	t.Parallel()
	for i := 0; i < ITERATIONS; i++ {
		token := GenerateToken()

		// Check raw text is relatively random
		for j := 0; j < 5; j++ {
			assert.NotEqual(t, token.Raw, GenerateToken().Raw)
		}

		// Check token length is less than 107 characters
		assert.Less(t, len(token.Raw), 108)

		// Check hash is the same
		assert.Equal(t, token.Hash, HashToken(token.Raw))
	}
}
