package sessionresolver

import (
	"bitbucket.org/HeilaSystems/session"
	"bitbucket.org/HeilaSystems/session/models"
	"bitbucket.org/HeilaSystems/tokenauth"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type sessionWrapper struct {
	repo session.SessionRepo
}

const Token = "token"
const TimeLayoutYYYYMMDD_HHMMSS = "2006-01-02 15:04:05"
const DataVersionsKey = "versions"

type ActiveOrder struct {
	Id             string            `json:"id"`
	SubServiceType string            `json:"subServiceType"`
	StoreId        string            `json:"storeId"`
	TimeTo         time.Time         `json:"timeTo"`
	Tags           []string          `json:"tags"`
	Versions       map[string]string `json:"versions"`
}

type Otp struct {
	UUID string `json:"uuid"`
}

type CurrentSession struct {
	CustomerId             string            `json:"customerId"`
	ActiveOrderId          string            `json:"activeOrderId"`
	FakeNow                *string           `json:"fakeNow"`
	CacheVersions          map[string]string `json:"cacheVersions"`
	ActiveOrder            *ActiveOrder      `json:"activeOrder"`
	OtpData                *Otp              `json:"otpData"`
	getLatestCacheVersions func(time.Time) (map[string]string, error)
}

func (c CurrentSession) GetActiveOrder() (hasActiveOrder bool, id, subServiceType string, storeId string, timeTo time.Time, tags []string, cacheVersions map[string]string) {
	if c.ActiveOrder != nil {
		return true,
			c.ActiveOrder.Id,
			c.ActiveOrder.SubServiceType,
			c.ActiveOrder.StoreId,
			c.ActiveOrder.TimeTo,
			c.ActiveOrder.Tags,
			c.ActiveOrder.Versions
	} else {
		return false, "", "", "", time.Time{}, nil, nil
	}
}

func (c CurrentSession) GetOtpData() string {
	if c.OtpData != nil {
		return c.OtpData.UUID
	} else {
		return ""
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

func (sw sessionWrapper) GetSessionById(c context.Context, id string) (bool, session.Session, error) {
	s := CurrentSession{}
	ok, err := sw.repo.GetUserSessionByTokenToStruct(c, id, &s)
	return ok, s, err
}

func (s *sessionWrapper) GetCurrentSession(c context.Context) (session.Session, error) {
	return s.getCurrentSessionInt(c)
}

func (s *sessionWrapper) GetTokenData(c context.Context) (map[string]interface{}, error) {
	tokenData := make(map[string]interface{})
	if tokenDataJson, ok := c.Value(tokenauth.TokenDataContextKey).(string); !ok {
		return nil, fmt.Errorf("tokenDataNotFound")
	} else if err := json.Unmarshal([]byte(tokenDataJson), &tokenData); err != nil {
		return nil, fmt.Errorf("tokenDataNotValidJSON")
	}
	return tokenData, nil
}

func (s *sessionWrapper) GetTokenDataValueAsString(c context.Context, key string) (string, error) {
	tokenData, err := s.GetTokenData(c)
	if err != nil {
		return "", err
	}
	val, ok := tokenData[key]
	if !ok {
		return "", fmt.Errorf("valueInTokenDataNotFound")
	}
	strVal, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("valueIsNotString")
	}
	return strVal, nil
}

func (s *sessionWrapper) getCurrentSessionInt(c context.Context) (CurrentSession, error) {
	var order CurrentSession
	if sessionId, err := s.GetTokenDataValueAsString(c, "sessionId"); err != nil {
		return order, err
	} else if _, err := s.repo.GetUserSessionByTokenToStruct(c, sessionId, &order); err != nil {
		return order, err
	} else {
		order.getLatestCacheVersions = func(now time.Time) (map[string]string, error) {
			return s.repo.GetCacheVersions(c, now)
		}
		return order, nil
	}
}

func (s *sessionWrapper) SetActiveOrder(c context.Context, id, subServiceType string, storeId string, timeTo time.Time, tags []string) error {
	cSession, err := s.getCurrentSessionInt(c)
	if err != nil {
		return err
	}
	versions, err := cSession.GetCacheVersions()
	if err != nil {
		return err
	}
	cSession.ActiveOrder = &ActiveOrder{
		Id:             id,
		SubServiceType: subServiceType,
		StoreId:        storeId,
		TimeTo:         timeTo,
		Tags:           tags,
		Versions:       versions,
	}

	return s.repo.InsertOrUpdate(c, cSession.GetCurrentCustomerId(), cSession)
}

func (s *sessionWrapper) SetOtpData(c context.Context, uuid string) error {
	cSession, err := s.getCurrentSessionInt(c)
	if err != nil {
		return err
	}
	cSession.OtpData = &Otp{UUID: uuid}
	return s.repo.InsertOrUpdate(c, cSession.GetCurrentCustomerId(), cSession)
}

func (s *sessionWrapper) SetCustomerId(c context.Context, customerId string) error {
	cSession, err := s.getCurrentSessionInt(c)
	if err != nil {
		return err
	}
	cSession.CustomerId = customerId
	return s.repo.InsertOrUpdate(c, cSession.GetCurrentCustomerId(), cSession)
}

func (s sessionWrapper) VersionsFromSessionToContext(c context.Context) (context.Context, error) {
	curSession, err := s.GetCurrentSession(c)
	if err != nil {
		return nil, err
	}

	versions, err := curSession.GetCacheVersions()
	if err != nil {
		return nil, err
	}

	ok, _, _, _, _, _, orderVersions := curSession.GetActiveOrder()
	if ok && versions != nil {
		for key, val := range orderVersions {
			versions[key] = val
		}
	}

	b, err := json.Marshal(versions)
	if err != nil {
		return nil, err
	}
	c = context.WithValue(c, DataVersionsKey, string(b))
	return c, nil
}

func (s sessionWrapper) GetVersionsFromContext(c context.Context) (models.Versions, bool, error) {
	v := c.Value(DataVersionsKey)
	if v == nil {
		return nil, false, nil
	}

	versions := make(models.Versions)
	err := json.Unmarshal([]byte(v.(string)), &versions)
	if err != nil {
		return nil, false, err
	}

	return versions, true, nil
}

func (s *sessionWrapper) SetCurrentSession(c context.Context, sessionId, customerId string, activeOrderId string,
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
	return s.repo.InsertOrUpdate(c, sessionId, cSession)
}
