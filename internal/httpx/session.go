package httpx

import (
	"github.com/google/uuid"
	"github.com/sembraniteam/setetes/internal/cryptox/pasetox"
)

type (
	UserSession struct {
		ID        uuid.UUID
		Anonymous bool
	}

	UserSessionClaims struct {
		Claims pasetox.Claims
	}
)
