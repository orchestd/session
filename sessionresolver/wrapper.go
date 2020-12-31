package sessionresolver

import (
	"bitbucket.org/HeilaSystems/session"
	"context"
)

type sessionWrapper struct {
	repo session.SessionRepo
}

func (s *sessionWrapper) GetOrderIdFromToken(c context.Context,token string) (string,error) {
	type OrderId struct {
		OrderId string `json:"orderId"`
	}
	var order OrderId
	if err := s.repo.GetUserSessionByTokenToStruct(c,token,&order);err != nil {
		return "", err
	} else {
		return order.OrderId , nil
	}
}


