package taxman

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
)

func (s *Service) Register(ctx context.Context) {
	service := api.AgentServiceRegistration{
		ID:   s.ID,
		Name: s.config.ServiceName,
		Checks: api.AgentServiceChecks{
			&api.AgentServiceCheck{
				Interval: "10s",
				HTTP:     fmt.Sprintf("http://%s/health", s.listener.Addr().String()),
				Status:   "passing",
				DeregisterCriticalServiceAfter: "20s",
			},
		},
	}
	err := s.discovery.Consul().Agent().ServiceRegister(&service)
	if err != nil {
		s.errors <- errors.Wrap(err, "discovery error")
	}
	for {
		select {
		case <-ctx.Done():
			err = s.discovery.Consul().Agent().ServiceDeregister(service.ID)
			if err != nil {
				s.errors <- errors.Wrap(err, "failed to deregister")
				return
			}
			s.log.Debug("service deregistered")
			return
		case <-time.NewTimer(60 * time.Second).C:
			err = s.discovery.Consul().Agent().ServiceRegister(&service)
			if err != nil {
				s.errors <- errors.Wrap(err, "discovery error")
			}
		}
	}
}

//func (s *Service) AcquireLeadership(ctx context.Context) {
//	return
//	// TODO make it similar to ServiceRegister
//	var session *discovery.Session
//	var err error
//	ticker := time.NewTicker(5 * time.Second)
//	for {
//		select {
//		default:
//			if session == nil {
//				session, err = discovery.NewSession(s.discovery)
//				if err != nil {
//					s.errors <- errors.Wrap(err, "failed to register session")
//					continue
//				}
//				session.EndlessRenew()
//			}
//
//			ok, err := s.discovery.TryAcquire(&discovery.KVPair{
//				Key:     s.config.LeadershipKey,
//				Value:   []byte(fmt.Sprintf("%d", s.LocalPayoutCounter)),
//				Session: session,
//			})
//
//			if err != nil {
//				s.errors <- err
//				s.IsLeader = false
//				continue
//			}
//
//			if ok {
//				s.IsLeader = true
//			} else {
//				// probably will never happen, but just in case
//				s.IsLeader = false
//			}
//			<-ticker.C
//		case <-ctx.Done():
//			// TODO release leadership
//			s.log.Debug("releasing leadership")
//			return
//		}
//	}
//}
