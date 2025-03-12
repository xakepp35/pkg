package xuuid

import "github.com/gofrs/uuid/v5"

func Random() uuid.UUID {
	return uuid.Must(uuid.NewV4())
}
