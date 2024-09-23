package controller

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	controllerapi "sandtech.io/sand-ops/api/v1"
	"sandtech.io/sand-ops/internal/utils"
)

func (r FrontendDeployReconciler) reconcileFrontendIngress(ctx context.Context, frontendPod *controllerapi.FrontendDeploy, l logr.Logger) (networkingv1.Ingress, error) {
	l.Info("reconcilling frontend ingress")
	ingressResource, err := utils.GetIngress(frontendPod.Namespace, ctx, r.Client)
	ingress := &networkingv1.Ingress{}
	deploymentPath := "/"

	if err == nil {
		ingressError := r.Get(ctx, types.NamespacedName{Name: frontendPod.Namespace + "-ingress-service", Namespace: frontendPod.Namespace}, ingress)
		if ingressError == nil {

			if frontendPod.Spec.IsHost {
				deploymentPath = "/?(.*)"

			} else {
				deploymentPath = "/" + frontendPod.Name + "/?(.*)"
			}
			deploymentSvc := "/" + frontendPod.Name + "-frontend-svc"
			if deploymentPathExist, path, index := utils.IngressPathExists(ingress.Spec.Rules[0].HTTP.Paths, deploymentPath); deploymentPathExist {
				if path.Backend.Service.Name != deploymentSvc {
					ingress.Spec.Rules[0].HTTP.Paths[index].Backend.Service.Name = deploymentSvc
					ingress.Spec.Rules[0].HTTP.Paths[index].Backend.Service.Port = networkingv1.ServiceBackendPort{
						Number: frontendPod.Spec.Port,
					}
					return *ingress, r.Update(ctx, ingress)
				}
				return *ingress, fmt.Errorf(utils.FOUND)
			}

			// new path add
			newPathType := networkingv1.PathTypeImplementationSpecific

			newIngressPath := networkingv1.HTTPIngressPath{
				Path:     deploymentPath,
				PathType: &newPathType,
				Backend: networkingv1.IngressBackend{
					Service: &networkingv1.IngressServiceBackend{
						Name: deploymentSvc,
						Port: networkingv1.ServiceBackendPort{
							Number: frontendPod.Spec.Port,
						},
					},
				},
			}

			ingress.Spec.Rules[0].HTTP.Paths = append(ingress.Spec.Rules[0].HTTP.Paths, newIngressPath)

			return *ingress, r.Update(ctx, ingress)
		}

		if !errors.IsNotFound(ingressError) {
			return *ingress, ingressError
		}

		if frontendPod.Spec.IsHost {
			deploymentPath = "/?(.*)"
		} else {
			deploymentPath = "/" + frontendPod.Name + "/?(.*)"
		}

		pathType := networkingv1.PathTypeImplementationSpecific
		ingress = &networkingv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      frontendPod.Namespace + "-ingress-service",
				Namespace: frontendPod.Namespace,
				OwnerReferences: []metav1.OwnerReference{
					{
						APIVersion: ingressResource.APIVersion,
						Kind:       ingressResource.Kind,
						Name:       ingressResource.Name,
						UID:        ingressResource.UID,
						Controller: utils.DataTypePointerRef(true),
					},
				},
				Annotations: map[string]string{
					"nginx.ingress.kubernetes.io/use-regex":       "true",
					"nginx.ingress.kubernetes.io/rewrite-target":  "/$1",
					"nginx.ingress.kubernetes.io/proxy-body-size": "8m",
				},
			},
			Spec: networkingv1.IngressSpec{
				IngressClassName: utils.DataTypePointerRef("nginx-" + frontendPod.Namespace),
				Rules: []networkingv1.IngressRule{
					{
						IngressRuleValue: networkingv1.IngressRuleValue{
							HTTP: &networkingv1.HTTPIngressRuleValue{
								Paths: []networkingv1.HTTPIngressPath{
									{
										Path:     deploymentPath,
										PathType: &pathType,
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: frontendPod.Name + "-frontend-svc",
												Port: networkingv1.ServiceBackendPort{
													Number: frontendPod.Spec.Port,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		return *ingress, r.Create(ctx, ingress)
	}
	return *ingress, err
}

func (r FrontendDeployReconciler) reconcileFrontend(ctx context.Context, frontendPod *controllerapi.FrontendDeploy, l logr.Logger) (appsv1.Deployment, error) {
	l.Info("reconcilling frontend deployment")

	frontendDeployment := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: frontendPod.Name, Namespace: frontendPod.Namespace}, frontendDeployment)

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
					Controller: utils.DataTypePointerRef(true),
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
