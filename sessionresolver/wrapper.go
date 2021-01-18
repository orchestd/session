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

type CurrentSession struct {
	CustomerId    string `json:"customerId"`
	ActiveOrderId string `json:"activeOrderId"`
	FakeNow       *string `json:"fakeNow"`
}

func (c CurrentSession) GetActiveOrderId() string {
	return c.ActiveOrderId
}

func (c CurrentSession) GetCurrentCustomerId() string {
	return c.CustomerId
}

func (c CurrentSession) GetNow() (*time.Time, error) {
	if c.FakeNow != nil {
		if t, err := time.Parse(TimeLayoutYYYYMMDD_HHMMSS, *c.FakeNow); err != nil {
			return nil, err
		} else {
			return &t, nil
		}
	} else {
		t := time.Now()
		return &t, nil
	}
}

func (s *sessionWrapper) GetCurrentSession(c context.Context) (session.Session, error) {
	var order CurrentSession
	if val, ok := c.Value(Token).(string); !ok {
		return nil, fmt.Errorf("tokenNotFound")
	} else if err := s.repo.GetUserSessionByTokenToStruct(c, val, &order); err != nil {
		return nil, err
	} else {
		return order, nil
	}
}

func (s *sessionWrapper) SetCurrentSession(c context.Context, CustomerId string, ActiveOrderId string, FakeNow *string) error {
	cSession :=  CurrentSession{
		CustomerId:    CustomerId,
		ActiveOrderId: ActiveOrderId,
		FakeNow:       FakeNow,
	}
	return s.repo.Insert(c, cSession.GetCurrentCustomerId(), cSession)
}
