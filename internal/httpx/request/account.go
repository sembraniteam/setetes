package request

import (
	"github.com/sembraniteam/setetes/internal/ent/account"
	"github.com/sembraniteam/setetes/internal/ent/otp"
)

type (
	Account struct {
		NationalID     string `json:"national_id"      validate:"required,len=16"`
		FullName       string `json:"full_name"        validate:"required,min=3,max=164"`
		Gender         string `json:"gender"           validate:"required,oneof=FEMALE MALE" reason:"oneof=gender must be one of FEMALE, MALE"`
		Email          string `json:"email"            validate:"required,email"`
		CountryISOCode string `json:"country_iso_code" validate:"required,iso3166_1_alpha2"  reason:"iso3166_1_alpha2=country_iso_code must be in ISO 3166-1 alpha-2 format"`
		DialCode       string `json:"dial_code"        validate:"required,min=1,max=6"`
		PhoneNumber    string `json:"phone_number"     validate:"required,min=11,max=13"`
	}

	Activation struct {
		Code           string `json:"otp_code"        validate:"required,len=6"`
		Password       string `json:"password"        validate:"required,min=8,max=128,password"                  reason:"password=password must include uppercase, lowercase, number, and special characters"`
		RetypePassword string `json:"retype_password" validate:"required,min=8,max=128,password,eqfield=Password" reason:"password=retype_password must include uppercase, lowercase, number, and special characters"`
	}

	ResendOTP struct {
		Email string `json:"email" validate:"required,email"`
		Type  string `json:"type"  validate:"required,oneof=ACTIVATION RESET_PASSWORD CHANGE_PASSWORD" reason:"oneof=type must be one of ACTIVATION, RESET_PASSWORD, CHANGE_PASSWORD"`
	}
)

func (a *Account) GetGender() account.Gender {
	switch a.Gender {
	case "FEMALE":
		return account.GenderFemale
	default:
		return account.GenderMale
	}
}

func (r *ResendOTP) GetType() otp.Type {
	switch r.Type {
	case "ACTIVATION":
		return otp.TypeActivation
	case "RESET_PASSWORD":
		return otp.TypeResetPassword
	default:
		return otp.TypeChangePassword
	}
}
