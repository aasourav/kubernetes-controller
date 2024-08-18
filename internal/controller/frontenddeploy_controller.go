/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	controllerapi "sandtech.io/sand-ops/api/v1"
	"sandtech.io/sand-ops/internal/utils"
)

// FrontendDeployReconciler reconciles a FrontendDeploy object
type FrontendDeployReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
	KubeClients
}

// +kubebuilder:rbac:groups=frontends.sandtech.io,resources=frontenddeploys,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=frontends.sandtech.io,resources=frontenddeploys/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=frontends.sandtech.io,resources=frontenddeploys/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the FrontendDeploy object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.4/pkg/reconcile
func (r *FrontendDeployReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := r.Log.WithValues("reconciling ", req.NamespacedName)

	frontendDeploy := &controllerapi.FrontendDeploy{}

	err := r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, frontendDeploy)
	if err != nil {
		if errors.IsNotFound(err) {
			l.Info(fmt.Sprintf("could not find: %s/%s", req.Name, req.Namespace))
			return ctrl.Result{}, nil
		}
		l.Error(err, fmt.Sprintf("failed to get: %s/%s", req.Name, req.Namespace))
		return ctrl.Result{}, err
	}

	frontendSvc, err := r.reconcileFrontendService(ctx, frontendDeploy, l)

	if err != nil {
		if err.Error() != utils.FOUND {
			l.Error(err, fmt.Sprintf("failed to create frontend service: %s/%s", frontendSvc.Name, frontendSvc.Namespace))
			return ctrl.Result{}, nil
		}
	} else {
		l.Info(fmt.Sprintf("successfully created frontend service: %s/%s", frontendSvc.Name, frontendSvc.Namespace))
	}

	frontendPod, err := r.reconcileFrontend(ctx, frontendDeploy, l)

	if err != nil {
		if err.Error() != utils.FOUND {
			l.Error(err, fmt.Sprintf("failed to create frontend deployment: %s/%s", frontendPod.Name, frontendPod.Namespace))
			return ctrl.Result{}, nil
		}
	} else {
		l.Info(fmt.Sprintf("successfully created frontend deployment: %s/%s", frontendPod.Name, frontendPod.Namespace))
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *FrontendDeployReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&controllerapi.FrontendDeploy{}).
		Complete(r)
}
