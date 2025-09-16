package validation




import (
    "github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(s interface{}) map[string]string {
    errors := make(map[string]string)

    if err := validate.Struct(s); err != nil {
        for _, e := range err.(validator.ValidationErrors) {
            field := e.Field()
            switch e.Tag() {
            case "required":
                errors[field] = field + " is required"
            case "min":
                errors[field] = field + " must be at least " + e.Param() + " characters"
            case "max":
                errors[field] = field + " must be at most " + e.Param() + " characters"
            case "email":
                errors[field] = "invalid email format"
            case "gt":
                errors[field] = field + " must be greater than " + e.Param()
            case "gte":
                errors[field] = field + " must be greater or equal to " + e.Param()
            case "oneof":
                errors[field] = field + " must be one of: " + e.Param()
            default:
                errors[field] = "invalid " + field
            }
        }
    }

    return errors
}
