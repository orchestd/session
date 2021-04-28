package session

import (
	"bitbucket.org/HeilaSystems/session/sessionresolver"
	"context"
	"time"
)

type SessionResolverBuilder interface {
	SetRepo(repo SessionRepo) SessionResolverBuilder
	Build() (SessionResolver, error)
}

type SessionResolver interface {
	GetCurrentSession(context context.Context) (Session, error)
	SetCurrentSession(c context.Context, customerId string, activeOrderId string, fakeNow *string, cacheVersions map[string]string) error
}

type Session interface {
	GetActiveOrderId() string
	GetCurrentCustomerId() string
	GetNow() (time.Time, error)
	GetCacheVersions() (map[string]string, error)
	GetActiveOrder() *sessionresolver.ActiveOrder
}

type SessionRepo interface {
	GetUserSessionByTokenToStruct(context context.Context, token string, dest interface{}) error
	InsertOrUpdate(ctx context.Context, id string, obj interface{}) error
	GetCacheVersions(ctx context.Context, now time.Time) (map[string]string, error)
}
