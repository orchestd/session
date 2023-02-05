package sessionresolver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/orchestd/session"
	"github.com/orchestd/session/models"
	"github.com/orchestd/tokenauth"
	"time"
)

type sessionWrapper struct {
	repo session.SessionRepo
}

const DataVersionsKey = "versions"
const DataNowKey = "dateNow"

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

type deviceInfo struct {
	Hardware    string
	Runtime     string
	OS          string
	DeviceModel string
	BrowserType string
	AppVersion  string
	OSVersion   string
}

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
	Referrer             string
	DeviceInfo           deviceInfo
	TermsApproval        bool
}

func (di deviceInfo) GetHardware() string {
	return di.Hardware
}

func (di deviceInfo) GetRuntime() string {
	return di.Runtime
}

func (di deviceInfo) GetOS() string {
	return di.OS
}

func (di deviceInfo) GetDeviceModel() string {
	return di.DeviceModel
}

func (di deviceInfo) GetBrowserType() string {
	return di.BrowserType
}

func (di deviceInfo) GetAppVersion() string {
	return di.AppVersion
}

func (di deviceInfo) GetOSVersion() string {
	return di.OSVersion
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
		// remove timezone
		l := "2006-01-02 15:04:05"
		d, _ := time.Parse(l, time.Now().Format(l))
		return d
	}
}

func (c currentSession) GetId() string {
	return c.Id
}

func (c *currentSession) SetLang(lang string) {
	c.Lang = lang
}

func (c currentSession) GetLang() string {
	return c.Lang
}

func (c *currentSession) SetTermsApproval(termsApproval bool) {
	c.TermsApproval = termsApproval
}

func (c currentSession) GetTermsApproval() bool {
	return c.TermsApproval
}

func (c *currentSession) SetDeviceInfo(hardware, runtime, os, deviceModel, browserType, appVersion, osVersion string) {
	c.DeviceInfo = deviceInfo{
		Hardware:    hardware,
		Runtime:     runtime,
		OS:          os,
		DeviceModel: deviceModel,
		BrowserType: browserType,
		AppVersion:  appVersion,
		OSVersion:   osVersion,
	}
}

func (c currentSession) GetDeviceInfo() session.DeviceInfoResolver {
	return c.DeviceInfo
}

func (c *currentSession) SetReferrer(referrer string) {
	c.Referrer = referrer
}

func (c currentSession) GetReferrer() string {
	return c.Referrer
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

func (s sessionWrapper) SetDataToContext(curSession session.Session, c context.Context) (context.Context, error) {
	c, err := s.versionsToContext(c, curSession)
	if err != nil {
		return c, err
	}

	c, err = s.nowToContext(c, curSession)
	if err != nil {
		return c, err
	}

	versions := make(map[string]string)
	fixedVersions := curSession.GetFixedCacheVersions()
	versionsForDate, err := s.repo.GetCacheVersions(c, curSession.GetNow())
	if err != nil {
		return c, err
	}
	for collection, ver := range versionsForDate {
		versions[collection] = ver
	}
	for collection, ver := range fixedVersions {
		versions[collection] = ver
	}
	curSession.SetCurrentCacheVersions(versions)

	c, err = s.versionsToContext(c, curSession)
	if err != nil {
		return c, err
	}

	return c, nil
}

func (s sessionWrapper) SetDataFromSessionToContext(c context.Context) (context.Context, error) {
	curSession, err := s.GetCurrentSession(c)
	if err != nil {
		return nil, err
	}
	c, err = s.versionsToContext(c, curSession)
	if err != nil {
		return c, err
	}

	c, err = s.nowToContext(c, curSession)
	if err != nil {
		return c, err
	}
	return c, nil
}

func (s sessionWrapper) nowToContext(c context.Context, curSession session.Session) (context.Context, error) {
	now := curSession.GetNow()
	b, err := json.Marshal(now)
	if err != nil {
		return nil, err
	}
	c = context.WithValue(c, DataNowKey, string(b))
	return c, nil
}

func (s sessionWrapper) versionsToContext(c context.Context, curSession session.Session) (context.Context, error) {
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
