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

type CustomerStatus int

const (
	NoCustomer CustomerStatus = iota
	NewCustomer
	ExistingCustomer
)

type currentSession struct {
	Id                   string
	CustomerId           string
	ActiveOrderId        string
	FakeNow              *time.Time
	FixedCacheVersions   map[string]string
	CurrentCacheVersions map[string]string
	ActiveOrder          *ActiveOrder
	OtpData              *Otp
	CustomerStatus       CustomerStatus
	Lang                 string
}

func (c *currentSession) SetCustomerDetails(id string, isNew bool) {
	c.CustomerId = id
	if id == "" {
		c.CustomerStatus = NoCustomer
	} else if isNew {
		c.CustomerStatus = NewCustomer
	} else {
		c.CustomerStatus = ExistingCustomer
	}
}

func (c *currentSession) SetOtpData(uuid string) {
	c.OtpData = &Otp{UUID: uuid}
}

func (c *currentSession) SetFakeNow(fakeNow time.Time) {
	c.FakeNow = &fakeNow
}

func (c *currentSession) SetFixedCacheVersions(versions map[string]string) {
	c.FixedCacheVersions = make(map[string]string)
	for collection, version := range versions {
		c.FixedCacheVersions[collection] = version
	}
}

func (c *currentSession) SetCurrentCacheVersions(versions map[string]string) {
	c.CurrentCacheVersions = make(map[string]string)
	for collection, version := range versions {
		c.CurrentCacheVersions[collection] = version
	}
}

func (c currentSession) GetCurrentCacheVersions() map[string]string {
	return c.CurrentCacheVersions
}

func (c currentSession) GetFixedCacheVersions() map[string]string {
	return c.FixedCacheVersions
}

func (c currentSession) GetOtpData() string {
	if c.OtpData != nil {
		return c.OtpData.UUID
	} else {
		return ""
	}
}

func (c currentSession) GetIsNoCustomer() bool {
	return c.CustomerStatus == NoCustomer
}

func (c currentSession) GetIsCustomerNew() bool {
	return c.CustomerStatus == NewCustomer
}

func (c currentSession) GetCustomerId() string {
	return c.CustomerId
}

func (c currentSession) HasFakeNow() bool {
	return c.FakeNow != nil
}

func (c currentSession) GetNow() time.Time {
	if c.FakeNow != nil {
		return *c.FakeNow
	} else {
		return time.Now()
	}
}

func (c currentSession) GetId() string {
	return c.Id
}

func (c currentSession) SetLang(lang string) {
	c.Lang = lang
}

func (c currentSession) GetLang() string {
	return c.Lang
}

func (sw sessionWrapper) NewSession(id string) session.Session {
	newCurrentSession := &currentSession{Id: id, CustomerStatus: NoCustomer}
	return newCurrentSession
}

func (sw sessionWrapper) SaveSession(c context.Context, cSession session.Session) error {
	return sw.repo.InsertOrUpdate(c, cSession.GetId(), cSession)
}

func (sw sessionWrapper) FreezeCacheVersionsForSession(c context.Context, curSession session.Session) error {
	versions := make(map[string]string)
	fixedVersions := curSession.GetFixedCacheVersions()
	versionsForDate, err := sw.repo.GetCacheVersions(c, curSession.GetNow())
	if err != nil {
		return err
	}
	for collection, ver := range versionsForDate {
		versions[collection] = ver
	}
	for collection, ver := range fixedVersions {
		versions[collection] = ver
	}
	curSession.SetCurrentCacheVersions(versions)
	err = sw.SaveSession(c, curSession)
	if err != nil {
		return err
	}
	return nil
}

func (sw sessionWrapper) GetSessionById(c context.Context, id string) (bool, session.Session, error) {
	s := currentSession{}
	ok, err := sw.repo.GetUserSessionByTokenToStruct(c, id, &s)
	return ok, &s, err
}

func (s *sessionWrapper) GetCurrentSession(c context.Context) (session.Session, error) {
	var currentSession currentSession
	if sessionId, err := s.GetTokenDataValueAsString(c, "sessionId"); err != nil {
		return nil, err
	} else if _, err := s.repo.GetUserSessionByTokenToStruct(c, sessionId, &currentSession); err != nil {
		return nil, err
	} else {
		return &currentSession, nil
	}
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

func (s sessionWrapper) VersionsFromSessionToContext(c context.Context) (context.Context, error) {
	curSession, err := s.GetCurrentSession(c)
	if err != nil {
		return nil, err
	}
	versions := curSession.GetCurrentCacheVersions()

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

func (s sessionWrapper) GetVersionForCollectionFromContext(c context.Context, collectionName string) (string, error) {
	versions, _, err := s.GetVersionsFromContext(c)
	if err != nil {
		return "", err
	}
	ver, ok := versions[collectionName]
	if !ok {
		return "", fmt.Errorf("latest version for collection %v not  found", collectionName)
	}
	return ver, nil
}

func (s sessionWrapper) IsObsolete(c context.Context, sessionId string) (bool, error) {
	ok, _, err := s.GetSessionById(c, sessionId)
	return !ok, err
}
