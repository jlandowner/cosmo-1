package dashboard

import (
	"context"
	"fmt"
	"net/http"

	dashv1alpha1 "github.com/cosmo-workspace/cosmo/api/openapi/dashboard/v1alpha1"
	wsv1alpha1 "github.com/cosmo-workspace/cosmo/api/workspace/v1alpha1"
	"github.com/cosmo-workspace/cosmo/pkg/clog"
	"github.com/cosmo-workspace/cosmo/pkg/wsnet"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/utils/pointer"
)

func (s *Server) PutNetworkRule(ctx context.Context, userId string, workspaceName string, networkRuleName string, req dashv1alpha1.UpsertNetworkRuleRequest) (dashv1alpha1.ImplResponse, error) {
	log := clog.FromContext(ctx).WithCaller()
	log.Debug().Info("request", "userId", userId, "workspaceName", workspaceName, "networkRuleName", networkRuleName, "req", req)

	user := userFromContext(ctx)
	if user == nil {
		log.Info("user not found in context")
		return ErrorResponse(http.StatusInternalServerError, "")
	}

	ws := workspaceFromContext(ctx)
	if ws == nil {
		log.Info("workspace not found in context")
		return ErrorResponse(http.StatusInternalServerError, "")
	}
	before := ws.DeepCopy()

	netRule := wsv1alpha1.NetworkRule{
		PortName:   networkRuleName,
		PortNumber: int(req.PortNumber),
		Group:      pointer.String(req.Group),
		HTTPPath:   req.HttpPath,
	}
	log.Debug().Info("upserting network rule", "ws", ws.Name, "namespace", ws.Namespace, "netRule", netRule)

	var err error
	ws.Spec.Network, err = wsnet.UpsertNetRule(ws.Spec.Network, netRule)
	if err != nil {
		message := err.Error()
		log.Error(err, message, "userid", user.Name, "workspace", ws.Name, "netRuleName", networkRuleName)
		return ErrorResponse(http.StatusBadRequest, message)
	}

	log.DebugAll().Info("NetworkRule upserted", "ws", ws, "namespace", ws.Namespace, "netRule", netRule)
	log.DebugAll().PrintObjectDiff(before, ws)

	if equality.Semantic.DeepEqual(before, ws) {
		log.Info("no change", "userid", user.Name, "workspace", ws.Name, "netRuleName", networkRuleName)
		return ErrorResponse(http.StatusBadRequest, "")
	}

	if err := s.Klient.Update(ctx, ws); err != nil {
		message := "Failed to upsert network rule"
		log.Error(err, message, "userid", user.Name, "workspace", ws.Name, "netRuleName", networkRuleName)
		return ErrorResponse(http.StatusInternalServerError, message)
	}

	res := &dashv1alpha1.UpsertNetworkRuleResponse{}
	res.Message = "Successfully upserted network rule"
	res.NetworkRule = convertNetRuleTodashv1alpha1NetRule(netRule)
	return dashv1alpha1.Response(http.StatusOK, res), nil
}

func (s *Server) DeleteNetworkRule(ctx context.Context, userId string, workspaceName string, networkRuleName string) (dashv1alpha1.ImplResponse, error) {
	log := clog.FromContext(ctx).WithCaller()
	log.Debug().Info("request", "userId", userId, "workspaceName", workspaceName, "networkRuleName", networkRuleName)

	user := userFromContext(ctx)
	if user == nil {
		log.Info("user not found in context")
		return ErrorResponse(http.StatusInternalServerError, "")
	}

	ws := workspaceFromContext(ctx)
	if ws == nil {
		log.Info("workspace not found in context")
		return ErrorResponse(http.StatusInternalServerError, "")
	}

	before := ws.DeepCopy()

	var delRule *wsv1alpha1.NetworkRule
	for _, v := range ws.Spec.Network {
		if v.PortName == networkRuleName {
			delRule = v.DeepCopy()
		}
	}
	if delRule == nil {
		message := fmt.Sprintf("port name %s is not found", networkRuleName)
		log.Info(message, "userid", user.Name, "workspace", ws.Name, "netRuleName", networkRuleName)
		return ErrorResponse(http.StatusBadRequest, message)
	}

	ws.Spec.Network = wsnet.RemoveNetworkOverrideByName(ws.Spec.Network, *delRule)
	log.DebugAll().Info("NetworkRule removed", "ws", ws, "userid", user.Name, "netRuleName", networkRuleName)

	log.DebugAll().PrintObjectDiff(before, ws)
	if equality.Semantic.DeepEqual(before, ws) {
		log.Info("no change", "userid", user.Name, "workspace", ws.Name, "netRuleName", networkRuleName)
		return ErrorResponse(http.StatusBadRequest, "")
	}

	if err := s.Klient.Update(ctx, ws); err != nil {
		message := "Failed to remove network rule"
		log.Error(err, message, "userid", user.Name, "workspace", ws.Name, "netRuleName", networkRuleName)
		return ErrorResponse(http.StatusInternalServerError, message)
	}

	res := &dashv1alpha1.RemoveNetworkRuleResponse{}
	res.Message = "Successfully removed network rule"
	res.NetworkRule = convertNetRuleTodashv1alpha1NetRule(*delRule)
	return NormalResponse(http.StatusOK, res)
}

func convertNetRulesTodashv1alpha1NetRules(netRules []wsv1alpha1.NetworkRule, urlMap map[string]string, serviceMainPortName string) []dashv1alpha1.NetworkRule {
	apirules := make([]dashv1alpha1.NetworkRule, 0, len(netRules))
	for _, v := range netRules {
		if v.PortName == serviceMainPortName {
			continue
		}

		r := convertNetRuleTodashv1alpha1NetRule(v)
		r.Url = urlMap[v.PortName]

		apirules = append(apirules, r)
	}
	return apirules
}

func convertNetRuleTodashv1alpha1NetRule(v wsv1alpha1.NetworkRule) dashv1alpha1.NetworkRule {
	return dashv1alpha1.NetworkRule{
		PortName:   v.PortName,
		PortNumber: int32(v.PortNumber),
		Group:      *v.Group,
		HttpPath:   v.HTTPPath,
	}
}
