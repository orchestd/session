package sessionresolver

import (
	"bitbucket.org/HeilaSystems/session"
	"context"
	"fmt"
	"time"
)

type sessionWrapper struct {
	repo session.SessionRepo
}

const Token = "token"
const TimeLayoutYYYYMMDD_HHMMSS = "2006-01-02 15:04:05"

type ActiveOrder struct {
	SubServiceType string    `json:"subServiceType"`
	StoreId        string    `json:"storeId"`
	TimeTo         time.Time `json:"timeTo"`
	Tags           []string  `json:"tags"`
}

type CurrentSession struct {
	CustomerId             string            `json:"customerId"`
	ActiveOrderId          string            `json:"activeOrderId"`
	FakeNow                *string           `json:"fakeNow"`
	CacheVersions          map[string]string `json:"cacheVersions"`
	ActiveOrder            *ActiveOrder      `json:"activeOrder"`
	getLatestCacheVersions func(time.Time) (map[string]string, error)
}

func (c CurrentSession) GetActiveOrder() (hasActiveOrder bool, subServiceType string, storeId string, timeTo time.Time, tags []string) {
	if c.ActiveOrder != nil {
		return true, c.ActiveOrder.SubServiceType, c.ActiveOrder.StoreId, c.ActiveOrder.TimeTo, c.ActiveOrder.Tags
	} else {
		return false, "", "", time.Time{}, nil
	}
}

func (c CurrentSession) GetActiveOrderId() string {
	return c.ActiveOrderId
}

func (c CurrentSession) GetCurrentCustomerId() string {
	return c.CustomerId
}

func (c CurrentSession) GetNow() (time.Time, error) {
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

func (c CurrentSession) GetCacheVersions() (map[string]string, error) {
	now, err := c.GetNow()
	if err != nil {
		return nil, err
	}
	versions, err := c.getLatestCacheVersions(now)
	if err != nil {
		return nil, err
	}
	if c.CacheVersions == nil {
		c.CacheVersions = make(map[string]string)
	}
	for collection, version := range versions {
		if _, ok := c.CacheVersions[collection]; !ok {
			c.CacheVersions[collection] = version
		}
	}
	return c.CacheVersions, nil
}

func (s *sessionWrapper) GetCurrentSession(c context.Context) (session.Session, error) {
	return s.getCurrentSessionInt(c)
}

func (s *sessionWrapper) getCurrentSessionInt(c context.Context) (CurrentSession, error) {
	var order CurrentSession
	if val, ok := c.Value(Token).(string); !ok {
		return order, fmt.Errorf("tokenNotFound")
	} else if err := s.repo.GetUserSessionByTokenToStruct(c, val, &order); err != nil {
		return order, err
	} else {
		order.getLatestCacheVersions = func(now time.Time) (map[string]string, error) {
			return s.repo.GetCacheVersions(c, now)
		}
		return order, nil
	}
}

func (s *sessionWrapper) SetActiveOrder(c context.Context, subServiceType string, storeId string, timeTo time.Time, tags []string) error {
	cSession, err := s.getCurrentSessionInt(c)
	if err != nil {
		return err
	}
	cSession.ActiveOrder = &ActiveOrder{
		SubServiceType: subServiceType,
		StoreId:        storeId,
		TimeTo:         timeTo,
		Tags:           tags,
	}

	return s.repo.InsertOrUpdate(c, cSession.GetCurrentCustomerId(), cSession)
}

func (s *sessionWrapper) SetCurrentSession(c context.Context, customerId string, activeOrderId string,
	fakeNow *string, cacheVersions map[string]string) error {
	cSession := CurrentSession{
		CustomerId:    customerId,
		ActiveOrderId: activeOrderId,
		FakeNow:       fakeNow,
		CacheVersions: cacheVersions,
		getLatestCacheVersions: func(now time.Time) (map[string]string, error) {
			return s.repo.GetCacheVersions(c, now)
		},
	}
	return s.repo.InsertOrUpdate(c, cSession.GetCurrentCustomerId(), cSession)
}
