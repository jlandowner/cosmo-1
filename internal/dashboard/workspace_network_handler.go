package dashboard

import (
	"context"

	connect_go "github.com/bufbuild/connect-go"

	cosmov1alpha1 "github.com/cosmo-workspace/cosmo/api/v1alpha1"
	"github.com/cosmo-workspace/cosmo/pkg/apiconv"
	"github.com/cosmo-workspace/cosmo/pkg/clog"
	dashv1alpha1 "github.com/cosmo-workspace/cosmo/proto/gen/dashboard/v1alpha1"
)

func (s *Server) UpsertNetworkRule(ctx context.Context, req *connect_go.Request[dashv1alpha1.UpsertNetworkRuleRequest]) (*connect_go.Response[dashv1alpha1.UpsertNetworkRuleResponse], error) {
	log := clog.FromContext(ctx).WithCaller()
	log.Debug().Info("request", "req", req)

	if err := userAuthentication(ctx, req.Msg.UserName); err != nil {
		return nil, ErrResponse(log, err)
	}

	m := req.Msg

	r := cosmov1alpha1.NetworkRule{
		PortNumber:       m.NetworkRule.PortNumber,
		CustomHostPrefix: m.NetworkRule.CustomHostPrefix,
		HTTPPath:         m.NetworkRule.HttpPath,
		Public:           m.NetworkRule.Public,
	}

	netRule, err := s.Klient.AddNetworkRule(ctx, m.WsName, m.UserName, r, int(m.Index))
	if err != nil {
		return nil, ErrResponse(log, err)
	}

	rule := apiconv.C2D_NetworkRule(*netRule)
	res := &dashv1alpha1.UpsertNetworkRuleResponse{
		Message:     "Successfully upserted network rule",
		NetworkRule: &rule,
	}
	return connect_go.NewResponse(res), nil
}

func (s *Server) DeleteNetworkRule(ctx context.Context, req *connect_go.Request[dashv1alpha1.DeleteNetworkRuleRequest]) (*connect_go.Response[dashv1alpha1.DeleteNetworkRuleResponse], error) {
	log := clog.FromContext(ctx).WithCaller()
	log.Debug().Info("request", "req", req)

	if err := userAuthentication(ctx, req.Msg.UserName); err != nil {
		return nil, ErrResponse(log, err)
	}

	m := req.Msg

	delRule, err := s.Klient.DeleteNetworkRule(ctx, m.WsName, m.UserName, int(m.Index))
	if err != nil {
		return nil, ErrResponse(log, err)
	}

	rule := apiconv.C2D_NetworkRule(*delRule)
	res := &dashv1alpha1.DeleteNetworkRuleResponse{
		Message:     "Successfully removed network rule",
		NetworkRule: &rule,
	}
	return connect_go.NewResponse(res), nil
}
