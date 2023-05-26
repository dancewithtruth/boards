package types

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type UUID struct {
	uuid.UUID
}

func (u UUID) MarshalJSON() ([]byte, error) {
	fmt.Println("hey", u.UUID)
	if u.UUID == uuid.Nil {
		return []byte("null"), nil
	}
	return json.Marshal(u.UUID.String())
}
