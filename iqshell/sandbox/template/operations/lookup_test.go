package operations

import (
	"fmt"
	"testing"

	"github.com/qiniu/go-sdk/v7/sandbox"
	"github.com/stretchr/testify/assert"
)

func TestIsTemplateAliasNotFound_APIError404(t *testing.T) {
	err := &sandbox.APIError{StatusCode: 404}

	assert.True(t, isTemplateAliasNotFound(err))
}

func TestIsTemplateAliasNotFound_OtherAPIError(t *testing.T) {
	err := &sandbox.APIError{StatusCode: 500}

	assert.False(t, isTemplateAliasNotFound(err))
}

func TestIsTemplateAliasNotFound_WrappedAPIError404(t *testing.T) {
	err := fmt.Errorf("lookup template: %w", &sandbox.APIError{StatusCode: 404})

	assert.True(t, isTemplateAliasNotFound(err))
}

func TestIsTemplateAliasNotFound_OtherError(t *testing.T) {
	err := fmt.Errorf("network error")

	assert.False(t, isTemplateAliasNotFound(err))
}
