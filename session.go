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
	SetCurrentSession(c context.Context, CustomerId string, ActiveOrderId string, FakeNow *string) error
}
type Session interface {
	GetActiveOrderId() string
	GetCurrentCustomerId() string
	GetNow() (*time.Time,error)
}
type SessionRepo interface {
	GetUserSessionByTokenToStruct(context context.Context,token string,dest interface{}) error
	Insert(ctx context.Context , id string, obj interface{}) error
}

