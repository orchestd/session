package session

import (
	"context"
	"time"
)

type SessionResolverBuilder interface {
	SetRepo(repo SessionRepo) SessionResolverBuilder
	Build() (SessionResolver,error)
}

type SessionResolver interface {
	GetCurrentSession(context context.Context) (Session,error)
}
type Session interface {
	GetActiveOrderId() string
	GetCurrentCustomerId() string
	GetNow() (*time.Time,error)
}
type SessionRepo interface {
	GetUserSessionByTokenToStruct(context context.Context,token string,dest interface{}) error
}

