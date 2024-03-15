package provider

import (
	"errors"
	"log/slog"
	"strings"
)

func CheckForEmptyAndSet(chkVal *string, toVal string) {
	if strings.EqualFold(*chkVal, "") {
		*chkVal = toVal
	}
}

func CheckForEmpty(chkVal *string) error {
	if strings.EqualFold(*chkVal, "") {
		return errors.New("Attribute value is empty string.")
	}

	return nil
}

func CheckForEmptyAndSetDefault(chkVal *string, someMap map[string]string, mapKey string) error {
	funcName := "CheckForEmptyAndSet"
	slog.Debug(funcName)

	if strings.EqualFold(*chkVal, "") {
		defaultVal, ok := someMap[mapKey]
		if ok {
			*chkVal = defaultVal
		} else {
			slog.Debug(funcName, "Unsupported key", mapKey, "attribute", *chkVal)
			return errors.New("Default DB admin user not set in config.")
		}
	}

	return nil
}
