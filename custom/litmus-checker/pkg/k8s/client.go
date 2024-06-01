package k8s

import (
	"context"
	"errors"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/discovery"
	memory "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

type ResourceDef struct {
	Name      string
	Group     string
	Version   string
	Kind      string
	Namespace string
	Selectors string
}

var decUnstructured = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)

func CreateChaosDeployment(manifest []byte, dc discovery.DiscoveryInterface, dyn dynamic.Interface) (*unstructured.Unstructured, error) {

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))
	// Decode YAML manifest into unstructured.Unstructured
	obj := &unstructured.Unstructured{}
	_, gvk, err := decUnstructured.Decode(manifest, nil, obj)
	if err != nil {
		return nil, err
	}

	if obj.GetKind() != "ChaosEngine" {
		return nil, errors.New("not a ChaosEngine Spec")
	}
	// Find GVR
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return nil, err
	}

	// Obtain REST interface for the GVR
	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// Namespaced resources should specify the namespace
		dr = dyn.Resource(mapping.Resource).Namespace(obj.GetNamespace())
	} else {
		// for cluster-wide resources
		dr = dyn.Resource(mapping.Resource)
	}

	return dr.Create(context.TODO(), obj, metav1.CreateOptions{})
}

func GetResourceDetails(dc discovery.DiscoveryInterface, dyn dynamic.Interface, res ResourceDef) (*unstructured.Unstructured, error) {
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))
	// Find GVR
	mapping, err := mapper.RESTMapping(schema.GroupKind{Group: res.Group, Kind: res.Kind}, res.Version)
	if err != nil {
		return nil, err
	}
	// Obtain REST interface for the GVR
	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources should specify the namespace
		dr = dyn.Resource(mapping.Resource).Namespace(res.Namespace)
	} else {
		// for cluster-wide resources
		dr = dyn.Resource(mapping.Resource)
	}
	return dr.Get(context.TODO(), res.Name, metav1.GetOptions{})
}

func GetKubeConfig(kubeconfig *string) (*rest.Config, error) {
	// Use in-cluster config if kubeconfig path is specified
	if *kubeconfig == "" {
		return rest.InClusterConfig()
	}
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return config, err
	}
	return config, err
}

func GetDynamicClient(kubeconfig *string) (discovery.DiscoveryInterface, dynamic.Interface, error) {
	cfg, err := GetKubeConfig(kubeconfig)
	if err != nil {
		return nil, nil, err
	}

	dc, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return nil, nil, err
	}
	// Prepare the dynamic client
	dyn, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, nil, err
	}
	return dc, dyn, err
}
