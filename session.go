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
	VersionsFromSessionToContext(c context.Context) (context.Context, error)
	GetVersionsFromContext(c context.Context) (models.Versions, bool, error)
	GetSessionById(c context.Context, id string) (bool, Session, error)
	GetTokenDataValueAsString(c context.Context, key string) (string, error)
	NewSession(id string) Session
	SaveSession(c context.Context, cSession Session) error
	GetCurrentSession(c context.Context) (Session, error)
	FreezeCacheVersionsForSession(c context.Context, curSession Session) error
}

type Session interface {
	SetCustomerDetails(id string, isNew bool)
	SetOtpData(uuid string)
	SetFakeNow(fakeNow time.Time)
	GetFixedCacheVersions() map[string]string
	SetFixedCacheVersions(versions map[string]string)
	GetCurrentCacheVersions() map[string]string
	SetCurrentCacheVersions(versions map[string]string)
	GetOtpData() string
	GetIsNoCustomer() bool
	GetIsCustomerNew() bool
	GetCustomerId() string
	HasFakeNow() bool
	GetNow() time.Time
	GetId() string
}

type SessionRepo interface {
	GetUserSessionByTokenToStruct(context context.Context, token string, dest interface{}) (bool, error)
	InsertOrUpdate(ctx context.Context, id string, obj interface{}) error
	GetCacheVersions(ctx context.Context, now time.Time) (map[string]string, error)
}
