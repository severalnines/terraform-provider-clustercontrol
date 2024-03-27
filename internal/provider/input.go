package provider

import (
	"errors"
	"fmt"
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

	if *chkVal == "" {
		defaultVal, ok := someMap[mapKey]
		if ok {
			*chkVal = defaultVal
		} else {
			errStr := fmt.Sprintf("%s: Unsupported key: %s", *chkVal)
			slog.Warn(errStr)
			return errors.New(errStr)
		}
	}

	return nil
}
