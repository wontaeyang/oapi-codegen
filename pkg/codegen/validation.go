package codegen

import (
	"fmt"
	"math"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// adds validation tags
// https://github.com/go-playground/validator
func validationTags(s *openapi3.Schema, req bool) string {
	if s.Type == "array" {
		s = s.Items.Value
	}

	var values []string

	// bool value of false should not be required
	// integer with min 0 should not be required
	skipRequired := s.Type == "boolean" ||
		(s.Type == "integer" && s.Min != nil && *s.Min == 0)

	if req {
		if !skipRequired {
			values = append(values, "required")
		}
	} else {
		values = append(values, "omitempty")
	}

	addMin := func(min float64) {
		exclusiveMin := s.ExclusiveMin
		if s.Type == "integer" {
			// the "validator" panic when compare int with float
			if min != math.Ceil(min) {
				min = math.Ceil(min)
				exclusiveMin = false
			}
		}
		if exclusiveMin {
			values = append(values, fmt.Sprintf("gt=%g", min))
		} else {
			values = append(values, fmt.Sprintf("gte=%g", min))
		}
	}

	addMax := func(max float64) {
		exclusiveMax := s.ExclusiveMax
		if s.Type == "integer" {
			// the "validator" panic when compare int with float
			if max != math.Floor(max) {
				max = math.Floor(max)
				exclusiveMax = false
			}
		}
		if exclusiveMax {
			values = append(values, fmt.Sprintf("lt=%g", max))
		} else {
			values = append(values, fmt.Sprintf("lte=%g", max))
		}
	}

	if s.Min != nil {
		addMin(*s.Min)
	} else if s.MinLength != 0 {
		addMin(float64(s.MinLength))
	} else if s.MinItems != 0 {
		addMin(float64(s.MinItems))
	}

	if s.Max != nil {
		addMax(*s.Max)
	} else if s.MaxLength != nil {
		addMax(float64(*s.MaxLength))
	} else if s.MaxItems != nil {
		addMax(float64(*s.MaxItems))
	}

	if s.MultipleOf != nil {
		// todo: generate custom validation function MultipleOf
	}

	if s.Format != "" {
		switch s.Format {
		case "int32":
		case "int64":
		case "float":
		case "double":
		case "byte":
		case "binary":
		case "date":
		case "date-time":
		case "password":
		default:
			values = append(values, s.Format)
		}
	}

	if s.Pattern != "" {
		// todo: generate custom validation function with precompiled regex
	}

	if len(s.Enum) > 0 {
		var items []string
		for _, item := range s.Enum {
			typed := item.(string)
			if strings.Contains(typed, " ") {
				typed = fmt.Sprintf("'%s'", item)
			}
			items = append(items, typed)
		}
		values = append(values, "oneof="+strings.Join(items, " "))
	}

	if s.UniqueItems {
		values = append(values, "unique")
	}

	if s.MinProps > 0 {
		// todo
	}

	if s.MaxProps != nil {
		// todo
	}

	// prevent single omitempty validation
	if len(values) == 1 && values[0] == "omitempty" {
		values = []string{}
	}

	return strings.Join(values, ",")
}
