package responsetypes

import (
	"github.com/google/uuid"
	"github.com/sembraniteam/setetes/internal/ent"
	"github.com/sembraniteam/setetes/internal/ent/account"
)

type (
	Account struct {
		*ent.Account
	}

	AccountResponse struct {
		ID               uuid.UUID      `json:"id"`
		NationalIDMasked string         `json:"national_id_masked"`
		FullName         string         `json:"full_name"`
		Gender           account.Gender `json:"gender"`
		Email            string         `json:"email"`
		CountryIsoCode   string         `json:"country_iso_code"`
		DialCode         string         `json:"dial_code"`
		PhoneNumber      string         `json:"phone_number"`
		Activated        bool           `json:"activated"`
		Locked           bool           `json:"locked"`
		TempLockedAt     int64          `json:"temp_locked_at"`
		CreatedAt        int64          `json:"created_at"`
		UpdatedAt        int64          `json:"updated_at"`
		DeletedAt        int64          `json:"deleted_at"`
	}
)

func (a Account) ToResponse() AccountResponse {
	return AccountResponse{
		ID:               a.ID,
		NationalIDMasked: a.NationalIDMasked,
		FullName:         a.FullName,
		Gender:           a.Gender,
		Email:            a.Email,
		CountryIsoCode:   a.CountryIsoCode,
		DialCode:         a.DialCode,
		PhoneNumber:      a.PhoneNumber,
		Activated:        a.Activated,
		Locked:           a.Locked,
		TempLockedAt:     a.TempLockedAt,
		CreatedAt:        a.CreatedAt,
		UpdatedAt:        a.UpdatedAt,
		DeletedAt:        a.DeletedAt,
	}
}
