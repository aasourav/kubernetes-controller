package utils

import (
	"context"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/types"
	controllerapi "sandtech.io/sand-ops/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func DataTypePointerRef[T bool | string | int | int32 | int64](booleanValue T) *T {
	return &booleanValue
}

func FrontendSVCSuffixedString(name string) string {
	return name + "-frontend-svc"
}

func ReplicasOrDefaultReplicas(numberOfReplicas int32, defaultReplica int32) *int32 {
	if numberOfReplicas < defaultReplica {
		return &defaultReplica
	} else {
		return &numberOfReplicas
	}
}

func NodeSelectorLabel(nodeNames string) map[string]string {
	if nodeNames == "" {
		return nil
	} else {
		return map[string]string{
			"kubernetes.io/hostname": nodeNames,
		}
	}
}

func IngressLabel(labelType string) map[string]string {
	return map[string]string{
		"app.kubernetes.io/component": labelType,
		"app.kubernetes.io/instance":  "ingress-nginx",
		"app.kubernetes.io/name":      "ingress-nginx",
		"app.kubernetes.io/part-of":   "ingress-nginx",
		"app.kubernetes.io/version":   "1.11.2",
	}
}

func NSSuffixedNamespace(name string) string {
	return name + "-ns"
}

func GetIngress(namespace string, ctx context.Context, client client.Client) (*controllerapi.SandOpsIngress, error) {
	ingressName := namespace[:len(namespace)-3]
	ingress := &controllerapi.SandOpsIngress{}
	err := client.Get(ctx, types.NamespacedName{Name: ingressName, Namespace: namespace}, ingress)
	if err == nil {
		return ingress, nil
	}
	return nil, err
}

func IngressPathExists(paths []networkingv1.HTTPIngressPath, targetPath string) (bool, *networkingv1.HTTPIngressPath, int) {
	for i, p := range paths {
		if p.Path == targetPath {

			return true, &p, i
		}
	}
	return false, nil, 0
}
