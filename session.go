package session

import (
	"bitbucket.org/HeilaSystems/session/models"
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
	SetActiveOrder(c context.Context, id, subServiceType string, storeId string, timeTo time.Time, tags []string) error
	SetCustomerId(c context.Context,customerId string) error
	SetOtpData(c context.Context , uuid string)error
	VersionsFromSessionToContext(c context.Context) (context.Context, error)
	GetVersionsFromContext(c context.Context) (models.Versions, bool, error)
	GetSessionById(c context.Context, id string) (bool, Session, error)
}

type Session interface {
	GetActiveOrderId() string
	GetCurrentCustomerId() string
	GetNow() (time.Time, error)
	GetCacheVersions() (map[string]string, error)
	GetActiveOrder() (hasActiveOrder bool, id, subServiceType string, storeId string, timeTo time.Time, tags []string, cacheVersions map[string]string)
	GetOtpData() string
}

type SessionRepo interface {
	GetUserSessionByTokenToStruct(context context.Context, token string, dest interface{}) (bool, error)
	InsertOrUpdate(ctx context.Context, id string, obj interface{}) error
	GetCacheVersions(ctx context.Context, now time.Time) (map[string]string, error)
}
