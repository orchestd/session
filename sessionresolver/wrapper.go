package sessionresolver

import (
	"bitbucket.org/HeilaSystems/session"
	"context"
	"fmt"
)

type sessionWrapper struct {
	repo session.SessionRepo
}

const Token = "token"
type ActiveOrder struct {
	ActiveOrder string `json:"activeOrder"`
}
func (s *sessionWrapper) GetActiveOrder(c context.Context) (string,error) {
	var order ActiveOrder
	if val , ok := c.Value(Token).(string);!ok {
		return "", fmt.Errorf("tokenNotFound")
	}else if err := s.repo.GetUserSessionByTokenToStruct(c,val,&order);err != nil {
		return "", err
	} else {
		return order.ActiveOrder , nil
	}
}


