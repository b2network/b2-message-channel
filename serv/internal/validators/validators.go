package validators

import "github.com/go-playground/validator/v10"

func RegisterValidators(v *validator.Validate) {
	RegisterBridgeValidators(v)
}
