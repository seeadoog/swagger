package swagger

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

func PtrVal[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

type Validator interface {
	Validate(ctx ValidateCtx) ValidateFuncs
}

type ValidateFunc func() error
type ValidateFuncs []func() error

type ValidateCtx struct {
}
type validateField struct {
}

func (c ValidateCtx) Maximum(name string, max int, v int) ValidateFunc {
	return func() error {
		if v > max {
			return fmt.Errorf("%s value should be less or equal than %d", name, max)
		}
		return nil
	}
}
func (c ValidateCtx) Minimum(name string, min int, v int) ValidateFunc {
	return func() error {
		if v < min {
			return fmt.Errorf("%s value should be greater or equal than %d", name, min)
		}
		return nil
	}
}

func (c ValidateCtx) NotEmpty(name string, v string) ValidateFunc {
	return func() error {
		if v == "" {
			return fmt.Errorf("%s value should not be empty", name)
		}
		return nil
	}
}

func (c ValidateCtx) StrIn(name string, v string, enums ...string) ValidateFunc {
	return func() error {
		vv := v
		for _, enum := range enums {
			if vv == enum {
				return nil
			}
		}
		return fmt.Errorf("%s value should be one of %v", name, enums)
	}
}

func (c ValidateCtx) IntIn(name string, v int, enums ...int) ValidateFunc {
	return func() error {
		vv := v
		for _, enum := range enums {
			if vv == enum {
				return nil
			}
		}
		return fmt.Errorf("%s value should be one of %v", name, enums)
	}
}

func (c ValidateCtx) LessThan(name string, param int64, lessThen int64) ValidateFunc {
	return func() error {
		if param >= lessThen {
			return fmt.Errorf("%s value should be less than %d", name, lessThen)
		}
		return nil
	}
}
func (c ValidateCtx) GreaterThan(name string, param int64, greaterThen int64) ValidateFunc {
	return func() error {
		if param <= greaterThen {
			return fmt.Errorf("%s value should be greater than %d", name, greaterThen)
		}
		return nil
	}
}

func (c ValidateCtx) MaxLength(name string, param string, maxLength int) ValidateFunc {
	return func() error {
		if len(param) > maxLength {
			return fmt.Errorf("%s length should be less or equal than %d", name, maxLength)
		}
		return nil
	}
}

func (c ValidateCtx) MinLength(name string, param string, minLength int) ValidateFunc {
	return func() error {
		if len(param) < minLength {
			return fmt.Errorf("%s length should be large or equal than %d", name, minLength)
		}
		return nil
	}
}

func FormatValidatorError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("'%s' is required", e.Field())
	case "max":
		return fmt.Sprintf("'%s' must be less or eq than '%s'", e.Field(), e.Param())
	case "min":
		return fmt.Sprintf("'%s' must be greater or eq than '%s'", e.Field(), e.Param())
	case "email":
		return fmt.Sprintf("'%s' is not a valid email address", e.Field())
	case "gt":
		return fmt.Sprintf("'%s' must be greater than '%s'", e.Field(), e.Param())
	case "lt":
		return fmt.Sprintf("'%s' must be less than '%s'", e.Field(), e.Param())
	case "gte":
		return fmt.Sprintf("'%s' must be less than or equal '%s'", e.Field(), e.Param())
	case "lte":
		return fmt.Sprintf("'%s' must be less than or equal '%s'", e.Field(), e.Param())
	case "oneof":
		return fmt.Sprintf("'%s' must be oneof '%s'", e.Field(), e.Param())
	case "startswith":
		return fmt.Sprintf("'%s' must startswith '%s'", e.Field(), e.Param())
	case "endswith":
		return fmt.Sprintf("'%s' must  endswith '%s'", e.Field(), e.Param())
	case "contains":
		return fmt.Sprintf("'%s' must  contains '%s'", e.Field(), e.Param())
	case "excludes":
		return fmt.Sprintf("'%s' must  excludes '%s'", e.Field(), e.Param())
	case "uuid":
		return fmt.Sprintf("'%s' must be valid uuid", e.Field())
	default:
		if strings.HasPrefix(e.Tag(), "required") {
			return fmt.Sprintf("'%s' is required", e.Field())
		}
		return e.Error()
	}
}
