package utils

import (
	"fmt"
	"reflect"
	"go-ingest-service/internal/models"
	"github.com/go-playground/validator/v10"
)

// Create a single, reusable validator instance.
var validate = validator.New()

// ValidateStruct validates a single struct that has `validate` tags.
func ValidateStruct(s interface{}) []models.ValidationErrorDetail {
	var errors []models.ValidationErrorDetail 
	// Use the validator instance to validate the struct.
	if err := validate.Struct(s); err != nil {
		// If validation fails, format the errors into our custom ErrorDetail slice.
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, models.ValidationErrorDetail{ 
				Loc:  []string{err.Field()},
				Msg:  "Validation failed on tag '" + err.Tag() + "'",
				Type: "validation_error." + err.Tag(),
			})
		}
	}
	return errors
}

// ValidateBatch is now a generic function that validates a slice of any struct type.
func ValidateBatch(batch interface{}) []models.ValidationErrorDetail {
	var allErrors []models.ValidationErrorDetail 

	// Use reflection to check if the provided interface is a slice.
	slice := reflect.ValueOf(batch)
	if slice.Kind() != reflect.Slice {
		// This is a safeguard; should not happen if called correctly.
		allErrors = append(allErrors, models.ValidationErrorDetail{ 
			Loc:  []string{"batch"},
			Msg:  "Invalid type: expected a slice",
			Type: "internal_error",
		})
		return allErrors
	}

	// Iterate over the slice using reflection.
	for i := 0; i < slice.Len(); i++ {
		item := slice.Index(i).Interface()
		// Use the same validator instance to validate each item in the slice.
		if err := validate.Struct(item); err != nil {
			// If an item is invalid, add its errors to the list.
			for _, validationErr := range err.(validator.ValidationErrors) {
				allErrors = append(allErrors, models.ValidationErrorDetail{ // Use new struct name
					// Prepend the index to the location for clear error reporting.
					Loc:  []string{fmt.Sprintf("[%d].%s", i, validationErr.Field())},
					Msg:  "Validation failed on tag '" + validationErr.Tag() + "'",
					Type: "validation_error." + validationErr.Tag(),
				})
			}
		}
	}
	return allErrors
}