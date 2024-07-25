package validators

import (
	"github.com/go-playground/validator/v10"
	"strings"
)

type BridgeValidator struct {
	validator *validator.Validate
}

func (v *BridgeValidator) EthAddressValidate(fl validator.FieldLevel) bool {
	ethAddress := strings.ToLower(fl.Field().String())
	return strings.HasPrefix(ethAddress, "0x") && len(ethAddress) == 42
}

func (v *BridgeValidator) BtcAddressValidate(fl validator.FieldLevel) bool {
	btcAddress := strings.ToLower(fl.Field().String())
	return (strings.HasPrefix(btcAddress, "2") && len(btcAddress) == 35) ||
		(strings.HasPrefix(btcAddress, "m") && len(btcAddress) == 34) ||
		(strings.HasPrefix(btcAddress, "tb1") && (len(btcAddress) == 42 || len(btcAddress) == 62))
}

func (v *BridgeValidator) LogIndexValidate(fl validator.FieldLevel) bool {
	logIndex := fl.Field().Int()

	return logIndex >= 0
}

func RegisterBridgeValidators(v *validator.Validate) {
	err := v.RegisterValidation("logIndex", (&BridgeValidator{v}).LogIndexValidate)
	if err != nil {
		return
	}
	err = v.RegisterValidation("ethAddress", (&BridgeValidator{v}).EthAddressValidate)
	if err != nil {
		return
	}
	err = v.RegisterValidation("btcAddress", (&BridgeValidator{v}).BtcAddressValidate)
	if err != nil {
		return
	}
}
