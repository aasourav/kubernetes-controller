package utils

func BoolPointer(booleanValue bool) *bool {
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
