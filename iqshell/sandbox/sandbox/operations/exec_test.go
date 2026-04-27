package operations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEnvPairs_ValidPairs(t *testing.T) {
	result := parseEnvPairs([]string{"FOO=bar", "BAZ=qux"})
	assert.Equal(t, map[string]string{"FOO": "bar", "BAZ": "qux"}, result)
}

func TestParseEnvPairs_Empty(t *testing.T) {
	result := parseEnvPairs(nil)
	assert.Empty(t, result)

	result = parseEnvPairs([]string{})
	assert.Empty(t, result)
}

func TestParseEnvPairs_NoEquals(t *testing.T) {
	result := parseEnvPairs([]string{"INVALID"})
	assert.Empty(t, result)
}

func TestParseEnvPairs_EmptyKey(t *testing.T) {
	result := parseEnvPairs([]string{"=value"})
	assert.Empty(t, result)
}

func TestParseEnvPairs_EmptyValue(t *testing.T) {
	result := parseEnvPairs([]string{"KEY="})
	assert.Equal(t, map[string]string{"KEY": ""}, result)
}

func TestParseEnvPairs_ValueWithEquals(t *testing.T) {
	result := parseEnvPairs([]string{"DB_URL=postgres://host:5432/db?sslmode=disable"})
	assert.Equal(t, map[string]string{"DB_URL": "postgres://host:5432/db?sslmode=disable"}, result)
}

func TestParseEnvPairs_MixedValidInvalid(t *testing.T) {
	result := parseEnvPairs([]string{"GOOD=value", "BAD", "=empty_key", "ALSO_GOOD=123"})
	assert.Equal(t, map[string]string{"GOOD": "value", "ALSO_GOOD": "123"}, result)
}

func TestShellQuoteArgs_PreservesSimpleCommand(t *testing.T) {
	result := shellQuoteArgs([]string{"ls", "-la", "/tmp"})
	assert.Equal(t, "ls -la /tmp", result)
}

func TestShellQuoteArgs_QuotesShellCommandString(t *testing.T) {
	result := shellQuoteArgs([]string{"sh", "-lc", "cat /etc/os-release | head -5"})
	assert.Equal(t, "sh -lc 'cat /etc/os-release | head -5'", result)
}

func TestShellQuoteArgs_EscapesSingleQuote(t *testing.T) {
	result := shellQuoteArgs([]string{"printf", "it's ok"})
	assert.Equal(t, "printf 'it'\\''s ok'", result)
}

func TestShellQuoteArgs_QuotesEmptyArg(t *testing.T) {
	result := shellQuoteArgs([]string{"printf", ""})
	assert.Equal(t, "printf ''", result)
}
