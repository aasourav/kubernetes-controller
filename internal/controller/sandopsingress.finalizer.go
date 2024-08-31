package controller

import (
	"context"

	"github.com/go-logr/logr"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	controllerapi "sandtech.io/sand-ops/api/v1"
	"sandtech.io/sand-ops/internal/utils"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *SandOpsIngressReconciler) ingressControllerFinalizer(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress, l logr.Logger) (bool, error) {
	if ingressDeployment.ObjectMeta.DeletionTimestamp.IsZero() {
		controllerutil.AddFinalizer(ingressDeployment, utils.INGRESS_FINALIZER)
		if err := r.Update(ctx, ingressDeployment); err != nil {
			return false, err
		}
		l.Info("Added ingress finalizer")
	} else {
		if controllerutil.ContainsFinalizer(ingressDeployment, utils.INGRESS_FINALIZER) {
			if err := r.deleteWebhook(ctx, ingressDeployment); err != nil {
				return false, err
			}
			if err := r.deleteIngressClass(ctx, ingressDeployment); err != nil {
				return false, err
			}
			if err := r.deletePatchJob(ctx, ingressDeployment); err != nil {
				return false, err
			}
			if err := r.deleteIngressJob(ctx, ingressDeployment); err != nil {
				return false, err
			}
			if err := r.deleteDeployment(ctx, ingressDeployment); err != nil {
				return false, err
			}
			if err := r.deleteServiceAdmission(ctx, ingressDeployment); err != nil {
				return false, err
			}
			if err := r.deleteService(ctx, ingressDeployment); err != nil {
				return false, err
			}
			if err := r.deleteConfigMap(ctx, ingressDeployment); err != nil {
				return false, err
			}
			if err := r.deleteAdmissionClusterRoleBinding(ctx, ingressDeployment); err != nil {
				return false, err
			}
			if err := r.deleteClusterRoleBinding(ctx, ingressDeployment); err != nil {
				return false, err
			}
			if err := r.deleteAdmissionRoleBinding(ctx, ingressDeployment); err != nil {
				return false, err
			}
			if err := r.deleteRoleBinding(ctx, ingressDeployment); err != nil {
				return false, err
			}
			if err := r.deleteAdmissionClusterRole(ctx, ingressDeployment); err != nil {
				return false, err
			}
			if err := r.deleteClusterRole(ctx, ingressDeployment); err != nil {
				return false, err
			}
			if err := r.deleteAdmissionRole(ctx, ingressDeployment); err != nil {
				return false, err
			}
			if err := r.deleteRole(ctx, ingressDeployment); err != nil {
				return false, err
			}
			if err := r.deleteAdmissionServiceAccount(ctx, ingressDeployment); err != nil {
				return false, err
			}
			if err := r.deleteServiceAccount(ctx, ingressDeployment); err != nil {
				return false, err
			}
			if err := r.deleteNamespace(ctx, ingressDeployment); err != nil {
				return false, err
			}
		}
		if ingressDeployment.ResourceVersion != "" {
			controllerutil.RemoveFinalizer(ingressDeployment, utils.INGRESS_FINALIZER)

			if err := r.Update(ctx, ingressDeployment); err != nil {
				return false, err
			}
			l.Info("ingress finalizer removed")
			return true, nil
		}
		return true, nil
	}

	return false, nil
}

func (r *SandOpsIngressReconciler) deleteWebhook(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress) error {

	webhookConfig := &admissionregistrationv1.ValidatingWebhookConfiguration{}

	err := r.Get(ctx, types.NamespacedName{Name: "ingress-nginx-admission-" + utils.NSSuffixedNamespace(ingressDeployment.Name)}, webhookConfig)
	if errors.IsNotFound(err) {
		return nil
	}

	err = r.Delete(ctx, webhookConfig)
	if err != nil {
		return err
	}

	return nil
}

func (r *SandOpsIngressReconciler) deleteIngressClass(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress) error {

	ingressClass := &networkingv1.IngressClass{}
	err := r.Get(ctx, types.NamespacedName{Name: "nginx-" + utils.NSSuffixedNamespace(ingressDeployment.Name)}, ingressClass)

	err = r.Delete(ctx, ingressClass)
	if err != nil {
		return err
	}

	return nil
}

func (r *SandOpsIngressReconciler) deletePatchJob(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress) error {

	job := &batchv1.Job{}
	err := r.Get(ctx, types.NamespacedName{Name: "ingress-nginx-admission-patch", Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name)}, job)

	err = r.Delete(ctx, job)
	if err != nil {
		return err
	}

	return nil
}

func (r *SandOpsIngressReconciler) deleteIngressJob(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress) error {

	job := &batchv1.Job{}
	err := r.Get(ctx, types.NamespacedName{Name: "ingress-nginx-admission-create", Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name)}, job)

	err = r.Delete(ctx, job)
	if err != nil {
		return err
	}

	return nil
}

func (r *SandOpsIngressReconciler) deleteDeployment(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress) error {

	deployment := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: utils.INGRESS_NGINX_CONTROLLER, Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name)}, deployment)

	err = r.Delete(ctx, deployment)
	if err != nil {
		return err
	}

	return nil
}

func (r *SandOpsIngressReconciler) deleteServiceAdmission(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress) error {

	service := &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{Name: utils.INGRESS_NGINX_CONTROLLER_ADMISSION, Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name)}, service)

	err = r.Delete(ctx, service)
	if err != nil {
		return err
	}

	return nil
}

func (r *SandOpsIngressReconciler) deleteService(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress) error {

	service := &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{Name: utils.INGRESS_NGINX_CONTROLLER, Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name)}, service)

	err = r.Delete(ctx, service)
	if err != nil {
		return err
	}

	return nil
}

func (r *SandOpsIngressReconciler) deleteConfigMap(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress) error {

	configMap := &corev1.ConfigMap{}
	err := r.Get(ctx, types.NamespacedName{Name: utils.INGRESS_NGINX_CONTROLLER, Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name)}, configMap)

	err = r.Delete(ctx, configMap)
	if err != nil {
		return err
	}

	return nil
}

func (r *SandOpsIngressReconciler) deleteAdmissionClusterRoleBinding(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress) error {

	clusterRoleBinding := &rbacv1.ClusterRoleBinding{}
	err := r.Get(ctx, types.NamespacedName{Name: "ingress-nginx-admission-" + utils.NSSuffixedNamespace(ingressDeployment.Name), Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name)}, clusterRoleBinding)

	err = r.Delete(ctx, clusterRoleBinding)
	if err != nil {
		return err
	}

	return nil
}

func (r *SandOpsIngressReconciler) deleteClusterRoleBinding(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress) error {

	clusterRoleBinding := &rbacv1.ClusterRoleBinding{}
	err := r.Get(ctx, types.NamespacedName{Name: "ingress-nginx-" + utils.NSSuffixedNamespace(ingressDeployment.Name), Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name)}, clusterRoleBinding)

	err = r.Delete(ctx, clusterRoleBinding)
	if err != nil {
		return err
	}

	return nil
}

func (r *SandOpsIngressReconciler) deleteAdmissionRoleBinding(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress) error {

	roleBinding := &rbacv1.RoleBinding{}
	err := r.Get(ctx, types.NamespacedName{Name: utils.INGRESS_NGINX_ADMISSION, Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name)}, roleBinding)

	err = r.Delete(ctx, roleBinding)
	if err != nil {
		return err
	}

	return nil
}

func (r *SandOpsIngressReconciler) deleteRoleBinding(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress) error {

	roleBinding := &rbacv1.RoleBinding{}
	err := r.Get(ctx, types.NamespacedName{Name: utils.INGRESS_NGINX, Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name)}, roleBinding)

	err = r.Delete(ctx, roleBinding)
	if err != nil {
		return err
	}

	return nil
}

func (r *SandOpsIngressReconciler) deleteAdmissionClusterRole(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress) error {

	admissionClusterRole := &rbacv1.ClusterRole{}
	err := r.Get(ctx, types.NamespacedName{Name: utils.INGRESS_NGINX_ADMISSION, Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name)}, admissionClusterRole)

	err = r.Delete(ctx, admissionClusterRole)
	if err != nil {
		return err
	}

	return nil
}

func (r *SandOpsIngressReconciler) deleteClusterRole(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress) error {

	clusterRole := &rbacv1.ClusterRole{}
	err := r.Get(ctx, types.NamespacedName{Name: utils.INGRESS_NGINX, Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name)}, clusterRole)

	err = r.Delete(ctx, clusterRole)
	if err != nil {
		return err
	}

	return nil
}

func (r *SandOpsIngressReconciler) deleteAdmissionRole(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress) error {

	ingressAdmissionRole := &rbacv1.Role{}
	err := r.Get(ctx, types.NamespacedName{Name: utils.INGRESS_NGINX_ADMISSION, Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name)}, ingressAdmissionRole)

	err = r.Delete(ctx, ingressAdmissionRole)
	if err != nil {
		return err
	}

	return nil
}

func (r *SandOpsIngressReconciler) deleteRole(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress) error {

	ingressRole := &rbacv1.Role{}
	err := r.Get(ctx, types.NamespacedName{Name: utils.INGRESS_NGINX, Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name)}, ingressRole)

	err = r.Delete(ctx, ingressRole)
	if err != nil {
		return err
	}

	return nil
}

func (r *SandOpsIngressReconciler) deleteAdmissionServiceAccount(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress) error {

	serviceAccount := &corev1.ServiceAccount{}
	err := r.Get(ctx, types.NamespacedName{Name: utils.INGRESS_NGINX_ADMISSION, Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name)}, serviceAccount)

	err = r.Delete(ctx, serviceAccount)
	if err != nil {
		return err
	}

	return nil
}

func (r *SandOpsIngressReconciler) deleteServiceAccount(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress) error {

	serviceAccount := &corev1.ServiceAccount{}
	err := r.Get(ctx, types.NamespacedName{Name: utils.INGRESS_NGINX, Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name)}, serviceAccount)

	err = r.Delete(ctx, serviceAccount)
	if err != nil {
		return err
	}

	return nil
}

func (r *SandOpsIngressReconciler) deleteNamespace(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress) error {

	namespace := &corev1.Namespace{}
	err := r.Get(ctx, types.NamespacedName{Name: utils.NSSuffixedNamespace(ingressDeployment.Name), Namespace: ingressDeployment.Namespace}, namespace)

	err = r.Delete(ctx, namespace)
	if err != nil {
		return err
	}

	return nil
}
