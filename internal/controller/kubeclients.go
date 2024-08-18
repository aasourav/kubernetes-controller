package controller

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type KubeClients struct {
	CRDClientSet        *clientset.Clientset
	KubernetesClientSet *kubernetes.Clientset
	DynamicClientSet    *dynamic.DynamicClient
}

func GetCRDClientSet() (*clientset.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	client, err := clientset.NewForConfig(config)
	if err != nil {
		return client, err
	}
	return client, nil
}

func GetKubernetesClientSet() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return client, err
	}
	return client, nil
}

func GetDynamicClient() (*dynamic.DynamicClient, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return dynamicClient, err
	}
	return dynamicClient, nil
}

func LoadCRDs(l logr.Logger) {
	crdClientSet, err := GetCRDClientSet()
	if crdClientSet == nil {
		l.Error(err, "Failed to get CRD ClientSet")
	}
	crdList, err := crdClientSet.ApiextensionsV1().CustomResourceDefinitions().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		l.Error(err, "CRD List not found")
	}

	// 3. get new empty schema holder
	scheme := runtime.NewScheme()

	// 4. loop over all the crd and add to the schema
	for _, crd := range crdList.Items {
		for _, v := range crd.Spec.Versions {
			//fmt.Printf("GROUP = %s        VERSION = %s     KIND = %s\n", crd.Spec.Group, v.Name, crd.Spec.Names.Kind)
			scheme.AddKnownTypeWithName(
				schema.GroupVersionKind{
					Group:   crd.Spec.Group,
					Version: v.Name,
					Kind:    crd.Spec.Names.Kind,
				},
				&unstructured.Unstructured{},
			)
		}
	}
}

func GetAllClients() (KubeClients, error) {
	CRDClient, err := GetCRDClientSet()
	if err != nil {
		return KubeClients{}, err
	}
	KubernetesClient, err := GetKubernetesClientSet()
	if err != nil {
		return KubeClients{}, err
	}
	DynamicClient, err := GetDynamicClient()
	if err != nil {
		return KubeClients{}, err
	}
	return KubeClients{
		CRDClientSet:        CRDClient,
		KubernetesClientSet: KubernetesClient,
		DynamicClientSet:    DynamicClient,
	}, nil

}
