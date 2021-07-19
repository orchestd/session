package sessionmock

import (
	"bitbucket.org/HeilaSystems/session"
	"bitbucket.org/HeilaSystems/session/models"
	"context"
	"time"
)

// Session
type ActiveOrder struct {
	SubServiceType string    `json:"subServiceType"`
	StoreId        string    `json:"storeId"`
	TimeTo         time.Time `json:"timeTo"`
	Tags           []string  `json:"tags"`
}

type currentSessionMock struct {
	CustomerId    string       `json:"customerId"`
	ActiveOrderId string       `json:"activeOrderId"`
	FakeNow       *string      `json:"fakeNow"`
	ActiveOrder   *ActiveOrder `json:"activeOrder"`
}

func (c currentSessionMock) GetActiveOrder() (hasActiveOrder bool, subServiceType string, storeId string, timeTo time.Time, tags []string) {
	if c.ActiveOrder != nil {
		return true, c.ActiveOrder.SubServiceType, c.ActiveOrder.StoreId, c.ActiveOrder.TimeTo, c.ActiveOrder.Tags
	} else {
		return false, "", "", time.Time{}, nil
	}
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
	}
}

func (c currentSessionMock) GetCacheVersions() (map[string]string, error) {
	return nil, nil
}

// Session resolver
type sessionMockWrapper struct {
	CurrentSessionToken string
	CurrentSession      currentSessionMock
}

func (s *sessionMockWrapper) SetActiveOrder(c context.Context, subServiceType string, storeId string, timeTo time.Time, tags []string) error {
	s.CurrentSession.ActiveOrder = &ActiveOrder{
		SubServiceType: subServiceType,
		StoreId:        storeId,
		TimeTo:         timeTo,
		Tags:           tags,
	}
	return nil
}

func (s *sessionMockWrapper) VersionsFromSessionToContext(c context.Context) (context.Context, error) {
	return c, nil
}

func (s *sessionMockWrapper) GetVersionsFromContext(c context.Context) (models.Versions, bool, error) {
	versions := make(models.Versions)
	return versions, true, nil
}

func Test(ao ActiveOrder)error{
	return nil
}
func NewSessionMockWrapper(currentSessionToken string, customerId string, activeOrderId string, fakeNow *string,
	storeId string,timeTo time.Time,tags []string,
) session.SessionResolver {
	return &sessionMockWrapper{CurrentSessionToken: currentSessionToken, CurrentSession: currentSessionMock{
		CustomerId:    customerId,
		ActiveOrderId: activeOrderId,
		FakeNow:       fakeNow,
	}}
}

func (s *sessionMockWrapper) GetSessionById(c context.Context, id string) (bool, session.Session, error) {
	return true, s.CurrentSession, nil
}

func (s *sessionMockWrapper) GetCurrentSession(context context.Context) (session.Session, error) {
	return s.CurrentSession, nil
}

func (s *sessionMockWrapper) SetCurrentSession(c context.Context, customerId string, activeOrderId string, fakeNow *string, cacheVersions map[string]string) error {
	ms := currentSessionMock{
		CustomerId:    customerId,
		ActiveOrderId: activeOrderId,
		FakeNow:       fakeNow,
	}
	s.CurrentSession = ms
	return nil
}

const TimeLayoutYYYYMMDD_HHMMSS = "2006-01-02 15:04:05"
const DataVersionsKey = "versions"
