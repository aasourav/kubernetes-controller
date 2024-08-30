package controller

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	controllerapi "sandtech.io/sand-ops/api/v1"
	"sandtech.io/sand-ops/internal/utils"
)

func (r *SandOpsIngressReconciler) reconcileServiceAccountAdmission(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress, l logr.Logger) (corev1.ServiceAccount, error) {
	l.Info("reconcilling ingress service account")
	serviceAccount := &corev1.ServiceAccount{}
	err := r.Get(ctx, types.NamespacedName{Name: utils.INGRESS_NGINX_ADMISSION, Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name)}, serviceAccount)
	if err == nil {
		return *serviceAccount, fmt.Errorf(utils.FOUND)
	}

	if !errors.IsNotFound(err) {
		return *serviceAccount, err
	}

	serviceAccount = &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.INGRESS_NGINX_ADMISSION,
			Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name),
			Labels:    utils.IngressLabel(utils.ADMISSION_WEBHOOK),
			OwnerReferences: []metav1.OwnerReference{
				{
					Name:       ingressDeployment.Name,
					APIVersion: ingressDeployment.APIVersion,
					Kind:       ingressDeployment.Kind,
					UID:        ingressDeployment.UID,
					Controller: utils.DataTypePointerRef(true),
				},
			},
		},
		AutomountServiceAccountToken: utils.DataTypePointerRef(true),
	}

	return *serviceAccount, r.Create(ctx, serviceAccount)
}

func (r *SandOpsIngressReconciler) reconcileServiceAccountIngress(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress, l logr.Logger) (corev1.ServiceAccount, error) {
	l.Info("reconcilling ingress service account")
	serviceAccount := &corev1.ServiceAccount{}
	err := r.Get(ctx, types.NamespacedName{Name: utils.INGRESS_NGINX, Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name)}, serviceAccount)
	if err == nil {
		return *serviceAccount, fmt.Errorf(utils.FOUND)
	}

	if !errors.IsNotFound(err) {
		return *serviceAccount, err
	}

	serviceAccount = &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.INGRESS_NGINX,
			Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name),
			Labels:    utils.IngressLabel(utils.CONTROLLER),
			OwnerReferences: []metav1.OwnerReference{
				{
					Name:       ingressDeployment.Name,
					APIVersion: ingressDeployment.APIVersion,
					Kind:       ingressDeployment.Kind,
					UID:        ingressDeployment.UID,
					Controller: utils.DataTypePointerRef(true),
				},
			},
		},
		AutomountServiceAccountToken: utils.DataTypePointerRef(true),
	}

	return *serviceAccount, r.Create(ctx, serviceAccount)
}

func (r *SandOpsIngressReconciler) reconcileNamespace(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress, l logr.Logger) (corev1.Namespace, error) {
	l.Info("reconcilling ingress controller namespace")
	namespace := &corev1.Namespace{}
	err := r.Get(ctx, types.NamespacedName{Name: utils.NSSuffixedNamespace(ingressDeployment.Name), Namespace: ingressDeployment.Namespace}, namespace)
	if err == nil {
		return *namespace, fmt.Errorf(utils.FOUND)
	}

	if !errors.IsNotFound(err) {
		return *namespace, err
	}

	namespace = &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app.kubernetes.io/instance": "ingress-nginx",
				"app.kubernetes.io/name":     "ingress-nginx",
				"namespace":                  utils.NSSuffixedNamespace(ingressDeployment.Name),
			},
			Name: utils.NSSuffixedNamespace(ingressDeployment.Name),
		},
	}

	return *namespace, r.Create(ctx, namespace)
}

func (r *SandOpsIngressReconciler) reconcileIngressRole(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress, l logr.Logger) (rbacv1.Role, error) {
	l.Info("reoncilling ingress role")

	ingressRole := &rbacv1.Role{}
	err := r.Get(ctx, types.NamespacedName{Name: utils.INGRESS_NGINX, Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name)}, ingressRole)
	if err == nil {
		return *ingressRole, fmt.Errorf(utils.FOUND)
	}
	if !errors.IsNotFound(err) {
		return *ingressRole, err
	}

	ingressRole = &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.INGRESS_NGINX,
			Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name),
			Labels:    utils.IngressLabel(utils.CONTROLLER),
			OwnerReferences: []metav1.OwnerReference{
				{
					Name:       ingressDeployment.Name,
					UID:        ingressDeployment.UID,
					APIVersion: ingressDeployment.APIVersion,
					Kind:       ingressDeployment.Kind,
					Controller: utils.DataTypePointerRef(true),
				},
			},
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"namespaces"},
				Verbs:     []string{"get"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{
					"configmaps",
					"pods",
					"secrets",
					"endpoints",
				},
				Verbs: []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"services"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{"networking.k8s.io"},
				Resources: []string{"ingresses"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{"networking.k8s.io"},
				Resources: []string{"ingresses/status"},
				Verbs:     []string{"update"},
			},
			{
				APIGroups: []string{"networking.k8s.io"},
				Resources: []string{"ingressclasses"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups:     []string{"coordination.k8s.io"},
				ResourceNames: []string{"ingress-nginx-leader"},
				Resources:     []string{"leases"},
				Verbs:         []string{"get", "update"},
			},
			{
				APIGroups: []string{"coordination.k8s.io"},
				Resources: []string{"leases"},
				Verbs:     []string{"create"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"events"},
				Verbs:     []string{"create", "patch"},
			},
			{
				APIGroups: []string{"discovery.k8s.io"},
				Resources: []string{"endpointslices"},
				Verbs:     []string{"list", "watch", "get"},
			},
		},
	}

	return *ingressRole, r.Create(ctx, ingressRole)
}

func (r *SandOpsIngressReconciler) reconcileIngressRoleBinding(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress, l logr.Logger) (rbacv1.RoleBinding, error) {
	l.Info("reconciling ingress role binding")

	roleBinding := &rbacv1.RoleBinding{}
	err := r.Get(ctx, types.NamespacedName{Name: utils.INGRESS_NGINX, Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name)}, roleBinding)
	if err == nil {
		return *roleBinding, fmt.Errorf(utils.FOUND)
	}
	if !errors.IsNotFound(err) {
		return *roleBinding, err
	}

	roleBinding = &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.INGRESS_NGINX,
			Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name),
			Labels:    utils.IngressLabel(utils.CONTROLLER),
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: ingressDeployment.APIVersion,
					Kind:       ingressDeployment.Kind,
					Name:       ingressDeployment.Name,
					UID:        ingressDeployment.UID,
					Controller: utils.DataTypePointerRef(true),
				},
			},
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      utils.INGRESS_NGINX,
				Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name),
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     utils.INGRESS_NGINX,
		},
	}

	return *roleBinding, r.Create(ctx, roleBinding)
}

func (r *SandOpsIngressReconciler) reconcileIngressAdmissionRoleBinding(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress, l logr.Logger) (rbacv1.RoleBinding, error) {
	l.Info("reconciling ingress admission role binding")

	roleBinding := &rbacv1.RoleBinding{}
	err := r.Get(ctx, types.NamespacedName{Name: utils.INGRESS_NGINX_ADMISSION, Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name)}, roleBinding)
	if err == nil {
		return *roleBinding, fmt.Errorf(utils.FOUND)
	}
	if !errors.IsNotFound(err) {
		return *roleBinding, err
	}

	roleBinding = &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.INGRESS_NGINX_ADMISSION,
			Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name),
			Labels:    utils.IngressLabel(utils.ADMISSION_WEBHOOK),
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: ingressDeployment.APIVersion,
					Kind:       ingressDeployment.Kind,
					Name:       ingressDeployment.Name,
					UID:        ingressDeployment.UID,
					Controller: utils.DataTypePointerRef(true),
				},
			},
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      utils.INGRESS_NGINX_ADMISSION,
				Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name),
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     utils.INGRESS_NGINX_ADMISSION,
		},
	}

	return *roleBinding, r.Create(ctx, roleBinding)
}

func (r *SandOpsIngressReconciler) reconcileIngressAdmissionRole(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress, l logr.Logger) (rbacv1.Role, error) {
	l.Info("reoncilling ingress role")

	ingressAdmissionRole := &rbacv1.Role{}
	err := r.Get(ctx, types.NamespacedName{Name: utils.INGRESS_NGINX_ADMISSION, Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name)}, ingressAdmissionRole)
	if err == nil {
		return *ingressAdmissionRole, fmt.Errorf(utils.FOUND)
	}
	if !errors.IsNotFound(err) {
		return *ingressAdmissionRole, err
	}

	ingressAdmissionRole = &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.INGRESS_NGINX_ADMISSION,
			Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name),
			Labels:    utils.IngressLabel(utils.CONTROLLER),
			OwnerReferences: []metav1.OwnerReference{
				{
					Name:       ingressDeployment.Name,
					UID:        ingressDeployment.UID,
					APIVersion: ingressDeployment.APIVersion,
					Kind:       ingressDeployment.Kind,
					Controller: utils.DataTypePointerRef(true),
				},
			},
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"secrets"},
				Verbs:     []string{"get", "create"},
			},
		},
	}

	return *ingressAdmissionRole, r.Create(ctx, ingressAdmissionRole)
}

func (r *SandOpsIngressReconciler) reconcileClusterRole(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress, l logr.Logger) (rbacv1.ClusterRole, error) {
	l.Info("reconcilling sandopsingress cluster role")

	clusterRole := &rbacv1.ClusterRole{}
	err := r.Get(ctx, types.NamespacedName{Name: utils.INGRESS_NGINX, Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name)}, clusterRole)
	if err == nil {
		return *clusterRole, fmt.Errorf(utils.FOUND)
	}

	if !errors.IsNotFound(err) {
		return *clusterRole, err
	}

	clusterRole = &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.INGRESS_NGINX,
			Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name),
			Labels:    utils.IngressLabel(utils.CONTROLLER),
			OwnerReferences: []metav1.OwnerReference{
				{
					Name:       ingressDeployment.Name,
					UID:        ingressDeployment.UID,
					APIVersion: ingressDeployment.APIVersion,
					Kind:       ingressDeployment.Kind,
					Controller: utils.DataTypePointerRef(true),
				},
			},
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{
					"configmaps",
					"endpoints",
					"nodes",
					"pods",
					"secrets",
					"namespaces",
				},
				Verbs: []string{"list", "watch"},
			},
			{
				APIGroups: []string{"coordination.k8s.io"},
				Resources: []string{"leases"},
				Verbs:     []string{"list", "watch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"nodes"},
				Verbs:     []string{"get"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"services"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{"networking.k8s.io"},
				Resources: []string{"ingresses"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"events"},
				Verbs:     []string{"create", "patch"},
			},
			{
				APIGroups: []string{"networking.k8s.io"},
				Resources: []string{"ingresses/status"},
				Verbs:     []string{"update"},
			},
			{
				APIGroups: []string{"networking.k8s.io"},
				Resources: []string{"ingressclasses"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{"discovery.k8s.io"},
				Resources: []string{"endpointslices"},
				Verbs:     []string{"list", "watch", "get"},
			},
		},
	}
	return *clusterRole, r.Create(ctx, clusterRole)
}

func (r *SandOpsIngressReconciler) reconcileIngressClusterRoleBinding(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress, l logr.Logger) (rbacv1.ClusterRoleBinding, error) {
	l.Info("reconciling ingress admission role binding")

	clusterRoleBinding := &rbacv1.ClusterRoleBinding{}
	err := r.Get(ctx, types.NamespacedName{Name: "ingress-nginx-" + utils.NSSuffixedNamespace(ingressDeployment.Name), Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name)}, clusterRoleBinding)
	if err == nil {
		return *clusterRoleBinding, fmt.Errorf(utils.FOUND)
	}
	if !errors.IsNotFound(err) {
		return *clusterRoleBinding, err
	}

	clusterRoleBinding = &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ingress-nginx-" + utils.NSSuffixedNamespace(ingressDeployment.Name),
			Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name),
			Labels:    utils.IngressLabel(utils.CONTROLLER),
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: ingressDeployment.APIVersion,
					Kind:       ingressDeployment.Kind,
					Name:       ingressDeployment.Name,
					UID:        ingressDeployment.UID,
					Controller: utils.DataTypePointerRef(true),
				},
			},
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      utils.INGRESS_NGINX,
				Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name),
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     utils.INGRESS_NGINX,
		},
	}

	return *clusterRoleBinding, r.Create(ctx, clusterRoleBinding)
}

func (r *SandOpsIngressReconciler) reconcileAdmissionClusterRole(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress, l logr.Logger) (rbacv1.ClusterRole, error) {
	l.Info("reconciling admission cluster role")

	admissionClusterRole := &rbacv1.ClusterRole{}
	err := r.Get(ctx, types.NamespacedName{Name: utils.INGRESS_NGINX_ADMISSION, Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name)}, admissionClusterRole)

	if err == nil {
		return *admissionClusterRole, fmt.Errorf(utils.FOUND)
	}

	if !errors.IsNotFound(err) {
		return *admissionClusterRole, err
	}

	admissionClusterRole = &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.INGRESS_NGINX_ADMISSION,
			Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name),
			Labels:    utils.IngressLabel(utils.ADMISSION_WEBHOOK),
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"admissionregistration.k8s.io"},
				Resources: []string{"validatingwebhookconfigurations"},
				Verbs:     []string{"get", "update"},
			},
		},
	}

	return *admissionClusterRole, r.Create(ctx, admissionClusterRole)
}

func (r *SandOpsIngressReconciler) reconcileIngressAdmissionClusterRoleBinding(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress, l logr.Logger) (rbacv1.ClusterRoleBinding, error) {
	l.Info("reconciling ingress admission role binding")

	clusterRoleBinding := &rbacv1.ClusterRoleBinding{}
	err := r.Get(ctx, types.NamespacedName{Name: "ingress-nginx-admission-" + utils.NSSuffixedNamespace(ingressDeployment.Name), Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name)}, clusterRoleBinding)
	if err == nil {
		return *clusterRoleBinding, fmt.Errorf(utils.FOUND)
	}
	if !errors.IsNotFound(err) {
		return *clusterRoleBinding, err
	}

	clusterRoleBinding = &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ingress-nginx-admission-" + utils.NSSuffixedNamespace(ingressDeployment.Name),
			Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name),
			Labels:    utils.IngressLabel(utils.ADMISSION_WEBHOOK),
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: ingressDeployment.APIVersion,
					Kind:       ingressDeployment.Kind,
					Name:       ingressDeployment.Name,
					UID:        ingressDeployment.UID,
					Controller: utils.DataTypePointerRef(true),
				},
			},
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      utils.INGRESS_NGINX_ADMISSION,
				Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name),
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     utils.INGRESS_NGINX_ADMISSION,
		},
	}

	return *clusterRoleBinding, r.Create(ctx, clusterRoleBinding)
}

func (r *SandOpsIngressReconciler) reconcileConfigMap(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress, l logr.Logger) (corev1.ConfigMap, error) {
	l.Info("reconciling ingress configmap")

	configMap := &corev1.ConfigMap{}
	err := r.Get(ctx, types.NamespacedName{Name: utils.INGRESS_NGINX_CONTROLLER, Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name)}, configMap)
	if err == nil {
		return *configMap, fmt.Errorf(utils.FOUND)
	}
	if !errors.IsNotFound(err) {
		return *configMap, err
	}

	configMap = &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.INGRESS_NGINX_CONTROLLER,
			Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name),
			Labels:    utils.IngressLabel(utils.CONTROLLER),
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: ingressDeployment.APIVersion,
					Kind:       ingressDeployment.Kind,
					Name:       ingressDeployment.Name,
					UID:        ingressDeployment.UID,
					Controller: utils.DataTypePointerRef(true),
				},
			},
		},

		Data: map[string]string{
			"allow-snippet-annotations": "true",
		},
	}

	return *configMap, r.Create(ctx, configMap)
}

func (r *SandOpsIngressReconciler) reconcileService(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress, l logr.Logger) (corev1.Service, error) {
	l.Info("reconciling ingress service")

	service := &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{Name: utils.INGRESS_NGINX_CONTROLLER, Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name)}, service)
	if err == nil {
		return *service, fmt.Errorf(utils.FOUND)
	}
	if !errors.IsNotFound(err) {
		return *service, err
	}

	singleStack := corev1.IPFamilyPolicySingleStack
	service = &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.INGRESS_NGINX_CONTROLLER,
			Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name),
			Labels:    utils.IngressLabel(utils.CONTROLLER),
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: ingressDeployment.APIVersion,
					Kind:       ingressDeployment.Kind,
					Name:       ingressDeployment.Name,
					UID:        ingressDeployment.UID,
					Controller: utils.DataTypePointerRef(true),
				},
			},
		},
		Spec: corev1.ServiceSpec{
			IPFamilies:     []corev1.IPFamily{corev1.IPv4Protocol},
			IPFamilyPolicy: &singleStack,
			Ports: []corev1.ServicePort{
				{
					AppProtocol: utils.DataTypePointerRef("http"),
					Name:        "http",
					Port:        80,
					Protocol:    corev1.ProtocolTCP,
					TargetPort:  intstr.FromString("http"),
				},
				{
					AppProtocol: utils.DataTypePointerRef("https"),
					Name:        "https",
					Port:        443,
					Protocol:    corev1.ProtocolTCP,
					TargetPort:  intstr.FromString("https"),
				},
			},
			Selector: map[string]string{
				"app.kubernetes.io/component": utils.CONTROLLER,
				"app.kubernetes.io/instance":  utils.INGRESS_NGINX,
				"app.kubernetes.io/name":      utils.INGRESS_NGINX,
			},
			Type: corev1.ServiceTypeLoadBalancer,
		},
	}

	return *service, r.Create(ctx, service)
}

func (r *SandOpsIngressReconciler) reconcileServiceAdmission(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress, l logr.Logger) (corev1.Service, error) {
	l.Info("reconciling ingress service admission")

	service := &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{Name: utils.INGRESS_NGINX_CONTROLLER_ADMISSION, Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name)}, service)
	if err == nil {
		return *service, fmt.Errorf(utils.FOUND)
	}
	if !errors.IsNotFound(err) {
		return *service, err
	}

	singleStack := corev1.IPFamilyPolicySingleStack
	service = &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.INGRESS_NGINX_CONTROLLER_ADMISSION,
			Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name),
			Labels:    utils.IngressLabel(utils.CONTROLLER),
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: ingressDeployment.APIVersion,
					Kind:       ingressDeployment.Kind,
					Name:       ingressDeployment.Name,
					UID:        ingressDeployment.UID,
					Controller: utils.DataTypePointerRef(true),
				},
			},
		},
		Spec: corev1.ServiceSpec{
			IPFamilies:     []corev1.IPFamily{corev1.IPv4Protocol},
			IPFamilyPolicy: &singleStack,
			Ports: []corev1.ServicePort{
				{
					AppProtocol: utils.DataTypePointerRef("https"),
					Name:        "https-webhook",
					Port:        443,
					TargetPort:  intstr.FromString("webhook"),
				},
			},
			Selector: map[string]string{
				"app.kubernetes.io/component": utils.CONTROLLER,
				"app.kubernetes.io/instance":  utils.INGRESS_NGINX,
				"app.kubernetes.io/name":      utils.INGRESS_NGINX,
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	return *service, r.Create(ctx, service)
}

func (r *SandOpsIngressReconciler) reconcileIngressClass(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress, l logr.Logger) (networkingv1.IngressClass, error) {
	l.Info("reconciling ingress class")
	ingressClass := &networkingv1.IngressClass{}
	err := r.Get(ctx, types.NamespacedName{Name: "nginx-" + utils.NSSuffixedNamespace(ingressDeployment.Name)}, ingressClass)
	if err == nil {
		return *ingressClass, fmt.Errorf(utils.FOUND)
	}

	if !errors.IsNotFound(err) {
		return *ingressClass, err
	}

	ingressClass = &networkingv1.IngressClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx-" + utils.NSSuffixedNamespace(ingressDeployment.Name),
			Labels: map[string]string{
				"app.kubernetes.io/component": "controller",
				"app.kubernetes.io/instance":  "ingress-nginx",
				"app.kubernetes.io/name":      "ingress-nginx",
				"app.kubernetes.io/part-of":   "ingress-nginx",
				"app.kubernetes.io/version":   "1.11.2",
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: ingressDeployment.APIVersion,
					Kind:       ingressDeployment.Kind,
					Name:       ingressDeployment.Name,
					UID:        ingressDeployment.UID,
					Controller: utils.DataTypePointerRef(true),
				},
			},
		},
		Spec: networkingv1.IngressClassSpec{
			Controller: "k8s.io/ingress-nginx-" + utils.NSSuffixedNamespace(ingressDeployment.Name),
		},
	}

	return *ingressClass, r.Create(ctx, ingressClass)
}

func (r *SandOpsIngressReconciler) reconcileIngressWebhook(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress, l logr.Logger) (admissionregistrationv1.ValidatingWebhookConfiguration, error) {
	l.Info("reconciling ingress webhook")

	webhookConfig := &admissionregistrationv1.ValidatingWebhookConfiguration{}
	err := r.Get(ctx, types.NamespacedName{Name: "ingress-nginx-admission-" + utils.NSSuffixedNamespace(ingressDeployment.Name)}, webhookConfig)
	if err == nil {
		return *webhookConfig, fmt.Errorf(utils.FOUND)
	}
	if !errors.IsNotFound(err) {

		return *webhookConfig, err
	}

	webhookConfig = &admissionregistrationv1.ValidatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: "ingress-nginx-admission-" + utils.NSSuffixedNamespace(ingressDeployment.Name),
			Labels: map[string]string{
				"app.kubernetes.io/component": "admission-webhook",
				"app.kubernetes.io/instance":  "ingress-nginx",
				"app.kubernetes.io/name":      "ingress-nginx",
				"app.kubernetes.io/part-of":   "ingress-nginx",
				"app.kubernetes.io/version":   "1.11.2",
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: ingressDeployment.APIVersion,
					Kind:       ingressDeployment.Kind,
					Name:       ingressDeployment.Name,
					UID:        ingressDeployment.UID,
					Controller: utils.DataTypePointerRef(true),
				},
			},
		},
		Webhooks: []admissionregistrationv1.ValidatingWebhook{
			{
				Name: "validate.nginx.ingress.kubernetes.io",
				ClientConfig: admissionregistrationv1.WebhookClientConfig{
					Service: &admissionregistrationv1.ServiceReference{
						Name:      "ingress-nginx-controller-admission",
						Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name),
						Path:      utils.DataTypePointerRef("/networking/v1/ingresses"),
					},
				},
				FailurePolicy: (*admissionregistrationv1.FailurePolicyType)(utils.DataTypePointerRef("Fail")),
				MatchPolicy:   (*admissionregistrationv1.MatchPolicyType)(utils.DataTypePointerRef("Equivalent")),
				Rules: []admissionregistrationv1.RuleWithOperations{
					{
						Operations: []admissionregistrationv1.OperationType{"CREATE", "UPDATE"},
						Rule: admissionregistrationv1.Rule{
							APIGroups:   []string{"networking.k8s.io"},
							APIVersions: []string{"v1"},
							Resources:   []string{"ingresses"},
						},
					},
				},
				SideEffects: func(s admissionregistrationv1.SideEffectClass) *admissionregistrationv1.SideEffectClass {
					return &s
				}(admissionregistrationv1.SideEffectClassNone),
				AdmissionReviewVersions: []string{"v1"},
			},
		},
	}

	return *webhookConfig, r.Create(ctx, webhookConfig)
}

func (r *SandOpsIngressReconciler) reconcileJobPatchAdmissionCreate(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress, l logr.Logger) (batchv1.Job, error) {
	l.Info("reconciling ingress job admission patch create")

	job := &batchv1.Job{}
	err := r.Get(ctx, types.NamespacedName{Name: "ingress-nginx-admission-patch", Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name)}, job)
	if err == nil {
		return *job, fmt.Errorf(utils.FOUND)
	}
	if !errors.IsNotFound(err) {
		return *job, err
	}

	job = &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ingress-nginx-admission-patch",
			Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name),
			Labels: map[string]string{
				"app.kubernetes.io/component": "admission-webhook",
				"app.kubernetes.io/instance":  "ingress-nginx",
				"app.kubernetes.io/name":      "ingress-nginx",
				"app.kubernetes.io/part-of":   "ingress-nginx",
				"app.kubernetes.io/version":   "1.11.2",
			},
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: "ingress-nginx-admission-patch",
					Labels: map[string]string{
						"app.kubernetes.io/component": "admission-webhook",
						"app.kubernetes.io/instance":  "ingress-nginx",
						"app.kubernetes.io/name":      "ingress-nginx",
						"app.kubernetes.io/part-of":   "ingress-nginx",
						"app.kubernetes.io/version":   "1.11.2",
					},
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: ingressDeployment.APIVersion,
							Kind:       ingressDeployment.Kind,
							Name:       ingressDeployment.Name,
							UID:        ingressDeployment.UID,
							Controller: utils.DataTypePointerRef(true),
						},
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "patch",
							Image: "registry.k8s.io/ingress-nginx/kube-webhook-certgen:v1.4.3@sha256:a320a50cc91bd15fd2d6fa6de58bd98c1bd64b9a6f926ce23a600d87043455a3",
							Args: []string{
								"patch",
								"--webhook-name=ingress-nginx-admission",
								"--namespace=$(POD_NAMESPACE)",
								"--patch-mutating=false",
								"--secret-name=ingress-nginx-admission",
								"--patch-failure-policy=Fail",
							},
							Env: []corev1.EnvVar{
								{
									Name: "POD_NAMESPACE",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
							},
							SecurityContext: &corev1.SecurityContext{
								AllowPrivilegeEscalation: func(b bool) *bool { return &b }(false),
								Capabilities: &corev1.Capabilities{
									Drop: []corev1.Capability{"ALL"},
								},
								ReadOnlyRootFilesystem: utils.DataTypePointerRef(true),
								RunAsNonRoot:           utils.DataTypePointerRef(true),
								RunAsUser:              utils.DataTypePointerRef(int64(65532)),
								SeccompProfile: &corev1.SeccompProfile{
									Type: corev1.SeccompProfileTypeRuntimeDefault,
								},
							},
							ImagePullPolicy: corev1.PullIfNotPresent,
						},
					},
					NodeSelector: map[string]string{
						"kubernetes.io/os": "linux",
					},
					RestartPolicy:      corev1.RestartPolicyOnFailure,
					ServiceAccountName: "ingress-nginx-admission",
				},
			},
		},
	}

	return *job, r.Create(ctx, job)
}

func (r *SandOpsIngressReconciler) reconcileJobAdmissionCreate(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress, l logr.Logger) (batchv1.Job, error) {
	l.Info("reconciling ingress job admission create")

	job := &batchv1.Job{}
	err := r.Get(ctx, types.NamespacedName{Name: "ingress-nginx-admission-create", Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name)}, job)
	if err == nil {
		return *job, fmt.Errorf(utils.FOUND)
	}
	if !errors.IsNotFound(err) {
		return *job, err
	}

	job = &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ingress-nginx-admission-create",
			Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name),
			Labels: map[string]string{
				"app.kubernetes.io/component": "admission-webhook",
				"app.kubernetes.io/instance":  "ingress-nginx",
				"app.kubernetes.io/name":      "ingress-nginx",
				"app.kubernetes.io/part-of":   "ingress-nginx",
				"app.kubernetes.io/version":   "1.11.2",
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: ingressDeployment.APIVersion,
					Kind:       ingressDeployment.Kind,
					Name:       ingressDeployment.Name,
					UID:        ingressDeployment.UID,
					Controller: utils.DataTypePointerRef(true),
				},
			},
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: "ingress-nginx-admission-create",
					Labels: map[string]string{
						"app.kubernetes.io/component": "admission-webhook",
						"app.kubernetes.io/instance":  "ingress-nginx",
						"app.kubernetes.io/name":      "ingress-nginx",
						"app.kubernetes.io/part-of":   "ingress-nginx",
						"app.kubernetes.io/version":   "1.11.2",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "create",
							Image: "registry.k8s.io/ingress-nginx/kube-webhook-certgen:v1.4.3@sha256:a320a50cc91bd15fd2d6fa6de58bd98c1bd64b9a6f926ce23a600d87043455a3",
							Args: []string{
								"create",
								"--host=ingress-nginx-controller-admission,ingress-nginx-controller-admission.$(POD_NAMESPACE).svc",
								"--namespace=$(POD_NAMESPACE)",
								"--secret-name=ingress-nginx-admission",
							},
							Env: []corev1.EnvVar{
								{
									Name: "POD_NAMESPACE",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
							},
							SecurityContext: &corev1.SecurityContext{
								AllowPrivilegeEscalation: utils.DataTypePointerRef(false),
								Capabilities: &corev1.Capabilities{
									Drop: []corev1.Capability{"ALL"},
								},
								ReadOnlyRootFilesystem: utils.DataTypePointerRef(true),
								RunAsNonRoot:           utils.DataTypePointerRef(true),
								RunAsUser:              utils.DataTypePointerRef(int64(65532)),
								SeccompProfile: &corev1.SeccompProfile{
									Type: corev1.SeccompProfileTypeRuntimeDefault,
								},
							},
							ImagePullPolicy: corev1.PullIfNotPresent,
						},
					},
					NodeSelector: map[string]string{
						"kubernetes.io/os": "linux",
					},
					RestartPolicy:      corev1.RestartPolicyOnFailure,
					ServiceAccountName: "ingress-nginx-admission",
				},
			},
		},
	}

	return *job, r.Create(ctx, job)
}

func (r *SandOpsIngressReconciler) reconcileIngressControllerDeployment(ctx context.Context, ingressDeployment *controllerapi.SandOpsIngress, l logr.Logger) (appsv1.Deployment, error) {
	l.Info("reconciling ingress deployment")

	deployment := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: utils.INGRESS_NGINX_CONTROLLER, Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name)}, deployment)
	if err == nil {
		return *deployment, fmt.Errorf(utils.FOUND)
	}

	if !errors.IsNotFound(err) {
		return *deployment, err
	}

	deployment = &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.INGRESS_NGINX_CONTROLLER,
			Namespace: utils.NSSuffixedNamespace(ingressDeployment.Name),
			Labels: map[string]string{
				"app.kubernetes.io/component": utils.CONTROLLER,
				"app.kubernetes.io/instance":  utils.INGRESS_NGINX,
				"app.kubernetes.io/name":      utils.INGRESS_NGINX,
				"app.kubernetes.io/part-of":   utils.INGRESS_NGINX,
				"app.kubernetes.io/version":   "1.11.2",
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: ingressDeployment.APIVersion,
					Kind:       ingressDeployment.Kind,
					Name:       ingressDeployment.Name,
					UID:        ingressDeployment.UID,
					Controller: utils.DataTypePointerRef(true),
				},
			},
		},
		Spec: appsv1.DeploymentSpec{
			MinReadySeconds:      0,
			RevisionHistoryLimit: utils.DataTypePointerRef(int32(10)),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/component": utils.CONTROLLER,
					"app.kubernetes.io/instance":  utils.INGRESS_NGINX,
					"app.kubernetes.io/name":      utils.INGRESS_NGINX,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app.kubernetes.io/component": utils.CONTROLLER,
						"app.kubernetes.io/instance":  utils.INGRESS_NGINX,
						"app.kubernetes.io/name":      utils.INGRESS_NGINX,
						"app.kubernetes.io/part-of":   utils.INGRESS_NGINX,
						"app.kubernetes.io/version":   "1.11.2",
					},
				},
				Spec: corev1.PodSpec{
					AutomountServiceAccountToken: utils.DataTypePointerRef(true),
					Containers: []corev1.Container{
						{
							Name:  utils.CONTROLLER,
							Image: "registry.k8s.io/ingress-nginx/controller:v1.8.1@sha256:e5c4824e7375fcf2a393e1c03c293b69759af37a9ca6abdb91b13d78a93da8bd",
							Args: []string{
								"/nginx-ingress-controller",
								"--election-id=ingress-nginx-leader",
								"--controller-class=k8s.io/ingress-nginx-" + utils.NSSuffixedNamespace(ingressDeployment.Name),
								"--ingress-class=nginx-" + utils.NSSuffixedNamespace(ingressDeployment.Name),
								"--configmap=$(POD_NAMESPACE)/ingress-nginx-controller",
								"--validating-webhook=:8443",
								"--validating-webhook-certificate=/usr/local/certificates/cert",
								"--validating-webhook-key=/usr/local/certificates/key",
								"--tcp-services-configmap=" + ingressDeployment.Name + "-ns/" + ingressDeployment.Name + "-ns-tcp-service-cm",
							},
							Env: []corev1.EnvVar{
								{
									Name: "POD_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.name",
										},
									},
								},
								{
									Name: "POD_NAMESPACE",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
								{
									Name:  "LD_PRELOAD",
									Value: "/usr/local/lib/libmimalloc.so",
								},
							},
							ImagePullPolicy: corev1.PullIfNotPresent,
							Lifecycle: &corev1.Lifecycle{
								PreStop: &corev1.LifecycleHandler{
									Exec: &corev1.ExecAction{
										Command: []string{"/wait-shutdown"},
									},
								},
							},
							LivenessProbe: &corev1.Probe{
								FailureThreshold: 5,
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path:   "/healthz",
										Port:   intstr.FromInt(10254),
										Scheme: corev1.URISchemeHTTP,
									},
								},
								InitialDelaySeconds: 10,
								PeriodSeconds:       10,
								SuccessThreshold:    1,
								TimeoutSeconds:      1,
							},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
								},
								{
									ContainerPort: 443,
									Name:          "https",
									Protocol:      corev1.ProtocolTCP,
								},
								{
									ContainerPort: 8443,
									Name:          "webhook",
									Protocol:      corev1.ProtocolTCP,
								},
							},
							ReadinessProbe: &corev1.Probe{
								FailureThreshold: 3,
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path:   "/healthz",
										Port:   intstr.FromInt(10254),
										Scheme: corev1.URISchemeHTTP,
									},
								},
								InitialDelaySeconds: 10,
								PeriodSeconds:       10,
								SuccessThreshold:    1,
								TimeoutSeconds:      1,
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("100m"),
									corev1.ResourceMemory: resource.MustParse("90Mi"),
								},
							},
							SecurityContext: &corev1.SecurityContext{
								AllowPrivilegeEscalation: utils.DataTypePointerRef(true),
								Capabilities: &corev1.Capabilities{
									Add:  []corev1.Capability{"NET_BIND_SERVICE"},
									Drop: []corev1.Capability{"ALL"},
								},
								RunAsUser: utils.DataTypePointerRef(int64(101)),
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									MountPath: "/usr/local/certificates/",
									Name:      "webhook-cert",
									ReadOnly:  true,
								},
							},
						},
					},
					DNSPolicy: corev1.DNSClusterFirst,
					NodeSelector: map[string]string{
						"kubernetes.io/os": "linux",
					},
					ServiceAccountName:            utils.INGRESS_NGINX,
					TerminationGracePeriodSeconds: utils.DataTypePointerRef(int64(300)),
					Volumes: []corev1.Volume{
						{
							Name: "webhook-cert",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: utils.INGRESS_NGINX_ADMISSION,
								},
							},
						},
					},
				},
			},
		},
	}

	return *deployment, r.Create(ctx, deployment)
}
