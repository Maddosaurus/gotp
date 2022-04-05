package lib

import (
	"errors"
	"fmt"
	"os"
	"strings"

	pb "github.com/Maddosaurus/pallas/proto/pallas"
	"github.com/gofrs/uuid"
)

func Getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func ValidateEntry(candiate *pb.OTPEntry) error {
	if _, err := uuid.FromString(candiate.Uuid); err != nil {
		return fmt.Errorf("ValidateEntry: failed to verify UUID: %w", err)
	}
	if len(candiate.SecretToken) != 16 || strings.Compare(candiate.SecretToken, strings.ToUpper(candiate.SecretToken)) != 0 {
		return errors.New("ValidateEntry: error while verifying token! Ensure it is 16 chars of upper case ASCII!")
	}
	// FIXME: Add OTP verification, but no error handling :/
	// https://github.com/xlzd/gotp/issues/18
	return nil
}
