package password

import (
	PwGen "github.com/sethvargo/go-password/password"
)

func Generate(length, digits, symbols int, allowRepeat bool) (string, error) {
	return PwGen.Generate(length, digits, symbols, false, allowRepeat)
}
