package utils

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
