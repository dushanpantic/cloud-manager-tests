package scope

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	cloudresourcesv1beta1 "github.com/kyma-project/cloud-resources-control-plane/api/cloud-resources/v1beta1"
	"github.com/kyma-project/cloud-resources-control-plane/pkg/common/composed"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createScopeGcp(ctx context.Context, st composed.State) (error, context.Context) {
	logger := composed.LoggerFromCtx(ctx)
	state := st.(*State)

	js, ok := state.CredentialData["serviceaccount.json"]
	if !ok {
		err := errors.New("gardener credential for gcp missing serviceaccount.json key")
		logger.Error(err, "error defining GCP scope")
		return composed.StopAndForget, nil // no requeue
	}

	var data map[string]string
	err := json.Unmarshal([]byte(js), &data)
	if err != nil {
		err := fmt.Errorf("error decoding serviceaccount.json: %w", err)
		logger.Error(err, "error defining GCP scope")
		return composed.StopAndForget, nil // no requeue
	}

	project, ok := data["project_id"]
	if !ok {
		err := errors.New("gardener gcp credentials missing project_id")
		logger.Error(err, "error defining GCP scope")
		return composed.StopAndForget, nil // no requeue
	}

	state.Scope = &cloudresourcesv1beta1.Scope{
		ObjectMeta: metav1.ObjectMeta{
			Name:      state.Obj().GetName(),
			Namespace: state.Obj().GetNamespace(),
		},
		Spec: cloudresourcesv1beta1.ScopeSpec{
			Kyma:      "",
			ShootName: "",
			Scope: cloudresourcesv1beta1.ScopeInfo{
				Gcp: &cloudresourcesv1beta1.GcpScope{
					Project:    project,
					VpcNetwork: fmt.Sprintf("shoot--%s--%s", state.ShootNamespace, state.ShootName),
				},
			},
		},
	}

	return nil, nil
}
