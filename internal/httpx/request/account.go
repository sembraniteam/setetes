package request

import "github.com/megalodev/setetes/internal/ent/account"

type Account struct {
	NationalID     string `json:"national_id"      validate:"required,len=16"`
	FullName       string `json:"full_name"        validate:"required,min=3,max=164"`
	Gender         string `json:"gender"           validate:"required,oneof=FEMALE MALE" reason:"oneof=gender must be one of FEMALE, MALE"`
	Email          string `json:"email"            validate:"required,email"`
	Password       string `json:"password"         validate:"required,min=8,password"    reason:"password=password must include uppercase, lowercase, number, and special characters"`
	CountryISOCode string `json:"country_iso_code" validate:"required,iso3166_1_alpha2"  reason:"iso3166_1_alpha2=country_iso_code must be in ISO 3166-1 alpha-2 format"`
	DialCode       string `json:"dial_code"        validate:"required,min=1,max=6"`
	PhoneNumber    string `json:"phone_number"     validate:"required,min=11,max=13"`
}

func (a *Account) GetGender() account.Gender {
	switch a.Gender {
	case "FEMALE":
		return account.GenderFemale
	default:
		return account.GenderMale
	}
}
