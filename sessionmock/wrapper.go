package sessionmock

import (
	"bitbucket.org/HeilaSystems/session"
	"context"
	"time"
)
type currentSessionMock struct {
	CustomerId             string            `json:"customerId"`
	ActiveOrderId          string            `json:"activeOrderId"`
	FakeNow                *string           `json:"fakeNow"`
	CacheVersions          map[string]string `json:"cacheVersions"`
	LatestCacheVersions func() map[string]string
}

func (c currentSessionMock) GetActiveOrderId() string {
	return c.ActiveOrderId
}

func (c currentSessionMock) GetCurrentCustomerId() string {
	return c.CustomerId
}

func (c currentSessionMock) GetNow() (time.Time, error) {
	if c.FakeNow != nil {
		if t, err := time.Parse(TimeLayoutYYYYMMDD_HHMMSS, *c.FakeNow); err != nil {
			return time.Time{}, err
		} else {
			return t, nil
		}
	} else {
		t := time.Now()
		return t, nil
	}}

func (c currentSessionMock) GetCacheVersions() (map[string]string, error) {
	return c.LatestCacheVersions() , nil
}

type sessionMockWrapper struct {
	CurrentSessionToken string
	CurrentSession currentSessionMock
}

func NewSessionMockWrapper(currentSessionToken string,customerId string, activeOrderId string, fakeNow *string, cacheVersions map[string]string ) session.SessionResolver {
	return &sessionMockWrapper{CurrentSessionToken: currentSessionToken, CurrentSession: currentSessionMock{
		CustomerId:          customerId,
		ActiveOrderId:       activeOrderId,
		FakeNow:             fakeNow,
		CacheVersions:       cacheVersions,
		LatestCacheVersions: func() map[string]string {
			return cacheVersions
		},
	}}
}

func (s *sessionMockWrapper) GetCurrentSession(context context.Context) (session.Session, error) {
	return s.CurrentSession, nil
}

func (s *sessionMockWrapper) SetCurrentSession(c context.Context, customerId string, activeOrderId string, fakeNow *string, cacheVersions map[string]string) error {
	ms := currentSessionMock{
		CustomerId:          customerId,
		ActiveOrderId:      activeOrderId,
		FakeNow:             fakeNow,
		CacheVersions:       cacheVersions,
		LatestCacheVersions: func() map[string]string {
			return cacheVersions
		},
	}
	s.CurrentSession = ms
	return nil
}

const TimeLayoutYYYYMMDD_HHMMSS = "2006-01-02 15:04:05"

