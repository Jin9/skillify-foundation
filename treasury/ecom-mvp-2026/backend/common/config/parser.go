package config

import env "github.com/caarlos0/env/v11"

func ParseEnv[T any](opts env.Options) (T, error) {
	var t T

	if err := env.Parse(&t); err != nil {
		return t, err
	}

	// Override with PREFIX_XXX if the prefixed env var exists.
	// This is intentionally best-effort: if the prefix is empty or missing,
	// we fall back to the non-prefixed value parsed above.
	//nolint:errcheck // intentional — prefix override is optional
	env.ParseWithOptions(&t, opts)

	return t, nil
}
