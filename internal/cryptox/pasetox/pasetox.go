package pasetox

import (
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/sembraniteam/setetes/internal/cryptox"
)

const (
	audience = "com.sembraniteam.setetes"
	issuer   = "https://setetes.sembraniteam.com"
	kid      = "ed25519-v1"
)

type (
	Claims struct {
		Platform        string
		Subject         string
		Expiration      time.Time
		NotBefore       time.Time
		IssuedAt        time.Time
		TokenIdentifier string
	}

	TokenPair struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int64  `json:"expires_in"`
	}

	Config struct {
		keypair *cryptox.Keypair
		claims  *Claims
	}

	Verifier struct {
		keypair *cryptox.Keypair
	}
)

func New(keypair *cryptox.Keypair, claims Claims) *Config {
	return &Config{
		keypair: keypair,
		claims:  &claims,
	}
}

func NewVerifier(keypair *cryptox.Keypair) *Verifier {
	return &Verifier{keypair: keypair}
}

func (c *Config) Signed() (string, error) {
	token := paseto.NewToken()
	token.SetIssuer(issuer)
	token.SetSubject(c.claims.Subject)
	token.SetAudience(audience)
	token.SetExpiration(c.claims.Expiration)
	token.SetNotBefore(c.claims.NotBefore)
	token.SetIssuedAt(c.claims.IssuedAt)
	token.SetJti(c.claims.TokenIdentifier)
	token.SetString("platform", c.claims.Platform)

	secretKey, err := paseto.NewV4AsymmetricSecretKeyFromEd25519(
		c.keypair.PrivateKey(),
	)
	if err != nil {
		return "", err
	}

	return token.V4Sign(secretKey, []byte(kid)), nil
}

func (v *Verifier) Verify(token string) (*paseto.Token, error) {
	publicKey, err := paseto.NewV4AsymmetricPublicKeyFromEd25519(
		v.keypair.PublicKey(),
	)
	if err != nil {
		return nil, err
	}

	parser := paseto.NewParser()
	parser.AddRule(
		paseto.NotExpired(),
		paseto.NotBeforeNbf(),
		paseto.ForAudience(audience),
		paseto.IssuedBy(issuer),
	)

	return parser.ParseV4Public(publicKey, token, []byte(kid))
}
