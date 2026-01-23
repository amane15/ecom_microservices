package data

import (
	"unicode/utf8"

	"github.com/amane15/ecom_microservice/pkg/validator"
	"github.com/shopspring/decimal"
)

func ValidateVariantName(v *validator.Validator, name string) {
	v.Check(name != "", "name", "must be provided")
	v.Check(utf8.RuneCountInString(name) <= 255, "name", "must be at most 255 characters long")
}

func ValidateVariantSlug(v *validator.Validator, slug string) {
	v.Check(slug != "", "slug", "must be provided")
	v.Check(validator.Matches(slug, validator.HyphenatedRegex), "slug", "must be valid slug")
	v.Check(utf8.RuneCountInString(slug) <= 128, "slug", "must be at most 128 character long")
}

func ValidateVariantPrice(v *validator.Validator, price decimal.Decimal) {
	v.Check(price.GreaterThanOrEqual(decimal.NewFromInt(1)), "price", "must be greater than 0")
}

func ValidateVariant(v *validator.Validator, variant *ProductVariant) {
	ValidateVariantName(v, variant.Name)
	ValidateVariantSlug(v, variant.Slug)

	ValidateVariantPrice(v, variant.Price)
	v.Check(variant.ProductID >= 1, "product_id", "must be a valid product_id")
}

func ValidateUpdateVariantInput(v *validator.Validator, input *UpdateVariantInput) {
	if input.Name != nil {
		ValidateVariantName(v, *input.Name)
	}

	if input.Price != nil {
		ValidateVariantPrice(v, *input.Price)
	}
}
