package controller

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	controllerapi "sandtech.io/sand-ops/api/v1"
	"sandtech.io/sand-ops/internal/utils"
)

func (r FrontendDeployReconciler) reconcileFrontend(ctx context.Context, frontendPod *controllerapi.FrontendDeploy, l logr.Logger) (appsv1.Deployment, error) {
	l.Info("reconcilling frontend deployment")

	frontendDeployment := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: frontendDeployment.Name, Namespace: frontendDeployment.Namespace}, frontendDeployment)

	if err == nil {
		return *frontendDeployment, fmt.Errorf(utils.FOUND)
	}

	if !errors.IsNotFound(err) {
		return *frontendDeployment, err
	}

	envVars := []corev1.EnvVar{}
	if frontendPod.Spec.EnvironmentVarialbles != nil {
		for _, envVar := range frontendPod.Spec.EnvironmentVarialbles {
			envVars = append(envVars, corev1.EnvVar{
				Name:  envVar.Name,
				Value: envVar.Value,
			})
		}
	}

	frontendDeployment = &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      frontendPod.Name,
			Namespace: frontendPod.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: frontendPod.APIVersion,
					Kind:       frontendPod.Kind,
					Name:       frontendPod.Name,
					UID:        frontendPod.UID,
					Controller: utils.BoolPointer(true),
				},
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": frontendPod.Name,
				},
			},
			Replicas: utils.ReplicasOrDefaultReplicas(frontendPod.Spec.Replicas, 1),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": frontendPod.Name,
					},
				},
				Spec: corev1.PodSpec{
					NodeSelector: utils.NodeSelectorLabel(frontendPod.Spec.NodeName),
					Containers: []corev1.Container{
						{
							Name:  "container-1",
							Image: frontendPod.Spec.ImageName,
							Env:   envVars,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: frontendPod.Spec.Port,
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU: resource.MustParse("500m"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU: resource.MustParse("2"),
								},
							},
						},
					},
				},
			},
		},
	}

	return *frontendDeployment, r.Create(ctx, frontendDeployment)
}
