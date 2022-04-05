package lib

import (
	"os"
	"testing"

	pb "github.com/Maddosaurus/gotp/proto/gotp"
	"github.com/stretchr/testify/assert"
)

func Test_Getenv_Key(t *testing.T) {
	key_name := "gotp_ge_key"
	key_value := "testing_123"
	os.Setenv(key_name, key_value)
	assert.Equal(t, Getenv(key_name, "fallback"), key_value)
}

func Test_Getenv_Fallback(t *testing.T) {
	key_name := "gotp_ge_key"
	os.Unsetenv(key_name)
	assert.Equal(t, Getenv(key_name, "fallback"), "fallback")
}
func Test_ValidateEntry_Valid(t *testing.T) {
	err := ValidateEntry(&pb.OTPEntry{
		Uuid:        "38518e4a-0b71-4d85-b925-1abdf3b56b03",
		SecretToken: "1234567890ABCDEF",
	})
	assert.Empty(t, err)
}

func Test_ValidateEntry_Invalid_UUID(t *testing.T) {
	err := ValidateEntry(&pb.OTPEntry{
		Uuid:        "42!!324&",
		SecretToken: "1234567890ABCDEF",
	})
	assert.ErrorContains(t, err, "failed to verify UUID: uuid:")
}

func Test_ValidateEntry_Invalid_Secret_Len(t *testing.T) {
	err := ValidateEntry(&pb.OTPEntry{
		Uuid:        "38518e4a-0b71-4d85-b925-1abdf3b56b03",
		SecretToken: "1234",
	})
	assert.ErrorContains(t, err, "ValidateEntry: error while verifying token! Ensure it is 16 chars of upper case ASCII!")
}

func Test_ValidateEntry_Invalid_Secret_Case(t *testing.T) {
	err := ValidateEntry(&pb.OTPEntry{
		Uuid:        "38518e4a-0b71-4d85-b925-1abdf3b56b03",
		SecretToken: "abcdef1234567890",
	})
	assert.ErrorContains(t, err, "ValidateEntry: error while verifying token! Ensure it is 16 chars of upper case ASCII!")
}
