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

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	pkgcontroller "sigs.k8s.io/controller-runtime/pkg/controller"

	controllerapi "sandtech.io/sand-ops/api/v1"
	"sandtech.io/sand-ops/internal/utils"
)

// SandOpsIngressReconciler reconciles a SandOpsIngress object
type SandOpsIngressReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
	KubeClients
}

// +kubebuilder:rbac:groups=aasdev.sandtech.io,resources=sandopsingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=aasdev.sandtech.io,resources=sandopsingresses/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=aasdev.sandtech.io,resources=sandopsingresses/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SandOpsIngress object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.4/pkg/reconcile
func (r *SandOpsIngressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := r.Log.WithValues("reconciling:", req.NamespacedName)

	ingressResource := &controllerapi.SandOpsIngress{}
	err := r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, ingressResource)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		l.Error(err, fmt.Sprintf("failed to get: %s/%s", req.Name, req.Namespace))
		return ctrl.Result{}, err
	}

	ingressNamespaceResource, err := r.reconcileNamespace(ctx, ingressResource, l)
	if err != nil {
		if err.Error() != utils.FOUND {
			l.Error(err, fmt.Sprintf("failed to create namespace for ingress controller: %s/%s", ingressNamespaceResource.Name, ingressNamespaceResource.Namespace))
			return ctrl.Result{}, nil
		}
	} else {
		l.Info(fmt.Sprintf("successfully created namespace for ingress controller: %s/%s", ingressNamespaceResource.Name, ingressNamespaceResource.Namespace))
	}

	serviceAccountIngress, err := r.reconcileServiceAccountIngress(ctx, ingressResource, l)
	if err != nil {
		if err.Error() != utils.FOUND {
			l.Error(err, fmt.Sprintf("failed to create ingress service account: %s/%s", serviceAccountIngress.Name, serviceAccountIngress.Namespace))
			return ctrl.Result{}, nil
		}
	}

	serviceAccountAdmission, err := r.reconcileServiceAccountAdmission(ctx, ingressResource, l)
	if err != nil {
		if err.Error() != utils.FOUND {
			l.Error(err, fmt.Sprintf("failed to create ingress admission service account: %s/%s", serviceAccountAdmission.Name, serviceAccountAdmission.Namespace))
			return ctrl.Result{}, nil
		}
	}

	ingressRole, err := r.reconcileIngressRole(ctx, ingressResource, l)
	if err != nil {
		if err.Error() != utils.FOUND {
			l.Error(err, "failed to create ingress role")
			return ctrl.Result{}, nil
		}
	} else {
		l.Info(fmt.Sprintf("successfully created ingress role: %s/%s", ingressRole.Name, ingressRole.Namespace))
	}

	ingressAdmissionRole, err := r.reconcileIngressAdmissionRole(ctx, ingressResource, l)
	if err != nil {
		if err.Error() != utils.FOUND {
			l.Error(err, "failed to create ingress admission role")
			return ctrl.Result{}, nil
		}
	} else {
		l.Info(fmt.Sprintf("successfully created ingress admission role: %s/%s", ingressAdmissionRole.Name, ingressAdmissionRole.Namespace))
	}

	ingressClusterRole, err := r.reconcileClusterRole(ctx, ingressResource, l)
	if err != nil {
		if err.Error() != utils.FOUND {
			l.Error(err, "failed to create cluster role")
			return ctrl.Result{}, nil
		}
	} else {
		l.Info(fmt.Sprintf("successfully created ingress cluster role: %s/%s", ingressClusterRole.Name, ingressClusterRole.Namespace))
	}

	admissionClusterRole, err := r.reconcileAdmissionClusterRole(ctx, ingressResource, l)
	if err != nil {
		if err.Error() != utils.FOUND {
			l.Error(err, "failed to create admission cluster role")
			return ctrl.Result{}, nil
		}
	} else {
		l.Info(fmt.Sprintf("successfully created admission cluster role: %s/%s", admissionClusterRole.Name, admissionClusterRole.Namespace))
	}

	ingressRoleBinding, err := r.reconcileIngressRoleBinding(ctx, ingressResource, l)
	if err != nil {
		if err.Error() != utils.FOUND {
			l.Error(err, "failed to create ingress role binding")
			return ctrl.Result{}, nil
		}
	} else {
		l.Info(fmt.Sprintf("successfully created ingress role binding: %s/%s", ingressRoleBinding.Name, ingressRoleBinding.Namespace))
	}

	ingressAdmissionRoleBinding, err := r.reconcileIngressAdmissionRoleBinding(ctx, ingressResource, l)
	if err != nil {
		if err.Error() != utils.FOUND {
			l.Error(err, "failed to create ingress role binding")
			return ctrl.Result{}, nil
		}
	} else {
		l.Info(fmt.Sprintf("successfully created ingress role binding: %s/%s", ingressAdmissionRoleBinding.Name, ingressAdmissionRoleBinding.Namespace))
	}

	ingressClusterRoleBinding, err := r.reconcileIngressClusterRoleBinding(ctx, ingressResource, l)
	if err != nil {
		if err.Error() != utils.FOUND {
			l.Error(err, "failed to create ingress cluster role binding")
			return ctrl.Result{}, nil
		}
	} else {
		l.Info(fmt.Sprintf("successfully created ingress cluster role binding: %s/%s", ingressClusterRoleBinding.Name, ingressClusterRoleBinding.Namespace))
	}

	ingressAdmissionClusterRoleBinding, err := r.reconcileIngressAdmissionClusterRoleBinding(ctx, ingressResource, l)
	if err != nil {
		if err.Error() != utils.FOUND {
			l.Error(err, "failed to create ingress admission cluster role binding")
			return ctrl.Result{}, nil
		}
	} else {
		l.Info(fmt.Sprintf("successfully created ingress admission cluster role binding: %s/%s", ingressAdmissionClusterRoleBinding.Name, ingressAdmissionClusterRoleBinding.Namespace))
	}

	ingressConfigMap, err := r.reconcileConfigMap(ctx, ingressResource, l)
	if err != nil {
		if err.Error() != utils.FOUND {
			l.Error(err, "failed to create ingress configmap")
			return ctrl.Result{}, nil
		}
	} else {
		l.Info(fmt.Sprintf("successfully created ingress configmap: %s/%s", ingressConfigMap.Name, ingressConfigMap.Namespace))
	}

	ingressService, err := r.reconcileService(ctx, ingressResource, l)
	if err != nil {
		if err.Error() != utils.FOUND {
			l.Error(err, "failed to create ingress service")
			return ctrl.Result{}, nil
		}
	} else {
		l.Info(fmt.Sprintf("successfully created ingress service: %s/%s", ingressService.Name, ingressService.Namespace))
	}

	ingressAdmissionService, err := r.reconcileServiceAdmission(ctx, ingressResource, l)
	if err != nil {
		if err.Error() != utils.FOUND {
			l.Error(err, "failed to create ingress admission service")
			return ctrl.Result{}, nil
		}
	} else {
		l.Info(fmt.Sprintf("successfully created ingress admission service: %s/%s", ingressAdmissionService.Name, ingressAdmissionService.Namespace))
	}

	ingressDeployment, err := r.reconcileIngressControllerDeployment(ctx, ingressResource, l)
	if err != nil {
		if err.Error() != utils.FOUND {
			l.Error(err, "failed to create ingress deployment")
			return ctrl.Result{}, nil
		}
	} else {
		l.Info(fmt.Sprintf("successfully created ingress deployment: %s/%s", ingressDeployment.Name, ingressDeployment.Namespace))
	}

	jobAdmissionCreate, err := r.reconcileJobAdmissionCreate(ctx, ingressResource, l)
	if err != nil {
		if err.Error() != utils.FOUND {
			l.Error(err, "failed to create ingress job admission")
			return ctrl.Result{}, nil
		}
	} else {
		l.Info(fmt.Sprintf("successfully created ingress job admission: %s/%s", jobAdmissionCreate.Name, jobAdmissionCreate.Namespace))
	}

	jobAdmissionPatchCreate, err := r.reconcileJobPatchAdmissionCreate(ctx, ingressResource, l)
	if err != nil {
		if err.Error() != utils.FOUND {
			l.Error(err, "failed to create ingress job admission")
			return ctrl.Result{}, nil
		}
	} else {
		l.Info(fmt.Sprintf("successfully created ingress job admission: %s/%s", jobAdmissionPatchCreate.Name, jobAdmissionPatchCreate.Namespace))
	}

	ingressClass, err := r.reconcileIngressClass(ctx, ingressResource, l)
	if err != nil {
		if err.Error() != utils.FOUND {
			l.Error(err, "failed to create ingress class")
			return ctrl.Result{}, nil
		}
	} else {
		l.Info(fmt.Sprintf("successfully created ingress class: %s/%s", ingressClass.Name, ingressClass.Namespace))
	}

	ingressWebhook, err := r.reconcileIngressWebhook(ctx, ingressResource, l)
	if err != nil {
		if err.Error() != utils.FOUND {
			l.Error(err, "failed to create ingress webhook")
			return ctrl.Result{}, nil
		}
	} else {
		l.Info(fmt.Sprintf("successfully created ingress webhook: %s/%s", ingressWebhook.Name, ingressWebhook.Namespace))
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SandOpsIngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&controllerapi.SandOpsIngress{}).
		WithOptions(pkgcontroller.Options{MaxConcurrentReconciles: 2}).
		Owns(&corev1.Namespace{}).
		Owns(&corev1.ServiceAccount{}).
		Owns(&rbacv1.ClusterRole{}).
		Owns(&rbacv1.ClusterRoleBinding{}).
		Owns(&corev1.ServiceAccount{}).
		Owns(&rbacv1.Role{}).
		Owns(&rbacv1.RoleBinding{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Service{}).
		Owns(&appsv1.Deployment{}).
		Owns(&batchv1.Job{}).
		Owns(&corev1.Secret{}).
		Owns(&networkingv1.IngressClass{}).
		Owns(&admissionregistrationv1.ValidatingWebhookConfiguration{}).
		Owns(&networkingv1.NetworkPolicy{}).
		Owns(&networkingv1.Ingress{}).
		Complete(r)
}
