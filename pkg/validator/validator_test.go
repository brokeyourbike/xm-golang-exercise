package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatorCanReturnErrors(t *testing.T) {
	type dummy struct {
		Name string `validate:"required,gt=3"`
	}

	v := NewValidation()
	errs := v.Validate(dummy{Name: "jon"})

	assert.Len(t, errs, 1)
	assert.IsType(t, ValidationError{}, errs[0])

	assert.Equal(t, "Field validation for 'Name' failed on the 'gt' tag", errs[0].Error())
}
