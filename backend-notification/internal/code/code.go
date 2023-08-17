package code

import (
	"fmt"
	"math/rand"
)

func Generate() string {
	code := rand.Intn(10000)
	return fmt.Sprintf("%04d", code)
}
