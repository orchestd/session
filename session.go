package session

import (
	"context"
	"time"
)

type SessionResolverBuilder interface {
	SetRepo(repo SessionRepo) SessionResolverBuilder
	Build() (SessionResolver, error)
}

type SessionResolver interface {
	SetDataFromSessionToContext(c context.Context) (context.Context, error)
	SetDataFromCurrentSessionToContext(c context.Context, curSession Session) (context.Context, error)
	GetSessionById(c context.Context, id string) (bool, Session, error)
	GetTokenDataValueAsString(c context.Context, key string) (string, error)
	NewSession(id string) Session
	SaveSession(c context.Context, cSession Session) error
	GetCurrentSession(c context.Context) (Session, error)
	FreezeCacheVersionsForSession(c context.Context, curSession Session, action string) error
	UnFreezeCacheVersionsForSession(c context.Context, curSession Session, action string) error
	IsObsolete(c context.Context, sessionId string) (bool, error)
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
	SetLang(lang string)
	GetLang() string
	SetTermsApproval(termsApproval bool)
	GetTermsApproval() bool
	SetDeviceInfo(hardware, runtime, os, deviceModel, browserType, appVersion, osVersion string)
	GetDeviceInfo() DeviceInfoResolver
	SetReferrer(string)
	GetReferrer() string
}

type DeviceInfoResolver interface {
	GetHardware() string
	GetRuntime() string
	GetOS() string
	GetDeviceModel() string
	GetBrowserType() string
	GetAppVersion() string
	GetOSVersion() string
}

type SessionRepo interface {
	GetUserSessionByTokenToStruct(context context.Context, token string, dest interface{}) (bool, error)
	InsertOrUpdate(ctx context.Context, id string, obj interface{}) error
	GetCacheVersions(ctx context.Context, now time.Time, filterAction string) (map[string]string, error)
}
