package kubernetes

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

func (k *Config) PostExec() (err error) {
	err = k.authenticate()
	if err != nil {
		return err
	}
	if k.DeploymentList != nil {
		return k.RedeployDeployment()
	}
	if k.StatefulsetList != nil {
		return k.RedeployStatefulsets()
	}
	return nil
}

func (k *Config) RedeployDeployment() (err error) {
	redeployDate := time.Now().Format("2006-01-02-15_04_05")

	for _, deployment := range k.DeploymentList {
		// Retrieve the latest version of Deployment before attempting update
		deploymentObject, err := k.clientset.AppsV1().Deployments(k.Namespace).Get(deployment, metav1.GetOptions{})
		if err != nil {
			return err
		}

		retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
			deploymentObject.Spec.Template.ObjectMeta.Labels["redeployedByVaultSync"] = redeployDate
			_, err = k.clientset.AppsV1().Deployments(k.Namespace).Update(deploymentObject)
			return err
		})
		if err != nil {
			return retryErr
		}
	}
	return nil
}

func (k *Config) RedeployStatefulsets() (err error) {
	redeployDate := time.Now().Format("2006-01-02-15_04_05")

	for _, statefulset := range k.StatefulsetList {
		statefulsetObject, err := k.clientset.AppsV1().StatefulSets(k.Namespace).Get(statefulset, metav1.GetOptions{})
		if err != nil {
			return err
		}

		retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
			statefulsetObject.Spec.Template.ObjectMeta.Labels["redeployedByVaultSync"] = redeployDate
			_, err = k.clientset.AppsV1().StatefulSets(k.Namespace).Update(statefulsetObject)
			return err
		})
		return retryErr
	}
	return nil
}
