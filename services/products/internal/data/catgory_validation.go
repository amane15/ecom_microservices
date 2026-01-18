package data

import (
	"unicode/utf8"

	"github.com/amane15/ecom_microservice/services/proudcts/internal/validator"
)

func ValidateCategoryName(v *validator.Validator, name string) {
	v.Check(name != "", "name", "must be provided")
	v.Check(utf8.RuneCountInString(name) <= 255, "name", "must be at most 255 characters long")
}

func ValidateCategorySlug(v *validator.Validator, slug string) {
	v.Check(slug != "", "slug", "must be provided")
	v.Check(validator.Matches(slug, validator.HyphenatedRegex), "slug", "must be a valid slug")
	v.Check(utf8.RuneCountInString(slug) <= 128, "slug", "must be at most 128 characters long")
}

func ValidateCategory(v *validator.Validator, category *Category) {
	ValidateCategoryName(v, category.Name)
	ValidateCategorySlug(v, category.Slug)
}

func ValidateCategoryUpdate(v *validator.Validator, input *UpdateCategoryInput) {
	if input.Name != nil {
		ValidateCategoryName(v, *input.Name)
	}
}
