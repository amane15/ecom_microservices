package data

import (
	"unicode/utf8"

	"github.com/amane15/ecom_microservice/services/proudcts/internal/validator"
)

func ValidateSlug(v *validator.Validator, slug string) {
	v.Check(slug != "", "slug", "must be provided")
	v.Check(utf8.RuneCountInString(slug) <= 128, "slug", "must not be more that 128 characters long")
	v.Check(validator.Matches(slug, validator.HyphenatedRegex), "slug", "must be in a valid format. For e.g. lenovo-r5-16g")
}

func ValidateName(v *validator.Validator, name string) {
	v.Check(name != "", "name", "must be provided")
	v.Check(utf8.RuneCountInString(name) <= 255, "name", "must not be more than 255 characters long")
}

func ValidateProduct(v *validator.Validator, product *Product) {
	ValidateSlug(v, product.Slug)
	ValidateName(v, product.Name)
}

func ValidateUpdateProduct(v *validator.Validator, productInput *UpdateProductRow) {
	ValidateName(v, *productInput.Name)
}
