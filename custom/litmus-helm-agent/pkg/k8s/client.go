package k8s

import (
	"context"
	"flag"
	"fmt"
	"os"

	corev1r "k8s.io/api/core/v1"

	metav1r "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func ConnectKubeApi() *kubernetes.Clientset {
	config, err := getKubeConfig()
	if err != nil {
		fmt.Printf("❌ Cannot create config: " + err.Error() + "\n")
		os.Exit(1)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("❌ Cannot create clientset: " + err.Error() + "\n")
		os.Exit(1)
	}
	return clientset
}

func CreateConfigMap(configmapName string, configMapData map[string]string, NAMESPACE string, clientset *kubernetes.Clientset) {
	configMap := corev1r.ConfigMap{
		TypeMeta: metav1r.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1r.ObjectMeta{
			Name:      configmapName,
			Namespace: NAMESPACE,
		},
		Data: configMapData,
	}

	_, err := clientset.CoreV1().ConfigMaps(NAMESPACE).Update(context.TODO(), &configMap, metav1r.UpdateOptions{})
	if err != nil {
		fmt.Printf("❌ Cannot update configmap " + configmapName + " : " + err.Error() + "\n")
		os.Exit(1)
	}
}

func CreateSecret(secretName string, secretData map[string][]byte, NAMESPACE string, clientset *kubernetes.Clientset) {
	secret := corev1r.Secret{
		TypeMeta: metav1r.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1r.ObjectMeta{
			Name:      secretName,
			Namespace: NAMESPACE,
		},
		Data: secretData,
	}

	var sec *corev1r.Secret

	sec, err := clientset.CoreV1().Secrets(NAMESPACE).Update(context.TODO(), &secret, metav1r.UpdateOptions{})
	if err != nil {
		fmt.Printf("❌ Cannot update secret " + secretName + " : " + err.Error() + "\n")
		os.Exit(1)
	}
	_ = sec
}

func getKubeConfig() (*rest.Config, error) {
	kubeconfig := flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	flag.Parse()
	// Use in-cluster config if kubeconfig path is not specified
	if *kubeconfig == "" {
		return rest.InClusterConfig()
	}

	return clientcmd.BuildConfigFromFlags("", *kubeconfig)
}
