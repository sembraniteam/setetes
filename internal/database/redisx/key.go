package redisx

import "github.com/google/uuid"

type Key string

const AuthKey Key = "auth:"

func (k Key) String() string {
	return string(k)
}

func (k Key) WithSession(jti uuid.UUID) Key {
	return k + Key("session:"+jti.String())
}
