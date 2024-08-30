package controller

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	controllerapi "sandtech.io/sand-ops/api/v1"
	utils "sandtech.io/sand-ops/internal/utils"
)

func (r *FrontendDeployReconciler) reconcileFrontendService(ctx context.Context, frontendDeploy *controllerapi.FrontendDeploy, l logr.Logger) (corev1.Service, error) {
	l.Info("reconcilling frontend svc")

	frontendSvc := &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{Name: utils.FrontendSVCSuffixedString(frontendDeploy.Name), Namespace: frontendDeploy.Namespace}, frontendSvc)

	if err == nil {
		return *frontendSvc, fmt.Errorf(utils.FOUND)
	}

	if !errors.IsNotFound(err) {
		return *frontendSvc, err
	}
	fmt.Println("here is name::::: ", frontendDeploy.Name)
	frontendSvc = &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.FrontendSVCSuffixedString(frontendDeploy.Name),
			Namespace: frontendDeploy.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: frontendDeploy.APIVersion,
					Kind:       frontendDeploy.Kind,
					UID:        frontendDeploy.UID,
					Name:       frontendDeploy.Name,
					Controller: utils.DataTypePointerRef(true),
				},
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": frontendDeploy.Name,
			},
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       frontendDeploy.Spec.Port,
					TargetPort: intstr.FromInt(int(frontendDeploy.Spec.Port)),
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	return *frontendSvc, r.Create(ctx, frontendSvc)
}
