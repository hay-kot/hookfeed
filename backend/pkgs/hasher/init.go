package hasher

import (
	"fmt"
	"os"

	"github.com/alexedwards/argon2id"
)

var params = argon2id.DefaultParams

func init() { // nolint:gochecknoinits
	disableHas := os.Getenv("UNSAFE_PASSWORD_PROTECTION") == "yes_i_am_sure"

	if disableHas {
		fmt.Println("WARNING: Password protection is disabled. This is unsafe in production.")
		params = &argon2id.Params{
			Memory:      16 * 1024, // Very low memory
			Iterations:  1,         // Very low iterations
			Parallelism: 1,         // Very low parallelism
			SaltLength:  1,
			KeyLength:   1,
		}
	}
}
