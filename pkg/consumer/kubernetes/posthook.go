/*
 * Copyright 2020 Kulkarni, Ashish <thatInfrastructureGuy@gmail.com>
 * Author: Ashish Kulkarni
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package kubernetes

import (
	"strings"
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
		err = k.RedeployDeployment()
	}
	if k.StatefulsetList != nil {
		err = k.RedeployStatefulsets()
	}
	return err
}

func (k *Config) RedeployDeployment() (err error) {
	redeployDate := time.Now().Format("2006-01-02-15_04_05")

	for _, deployment := range k.DeploymentList {
		deployment = strings.TrimSpace(deployment)
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
		statefulset = strings.TrimSpace(statefulset)
		statefulsetObject, err := k.clientset.AppsV1().StatefulSets(k.Namespace).Get(statefulset, metav1.GetOptions{})
		if err != nil {
			return err
		}

		retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
			statefulsetObject.Spec.Template.ObjectMeta.Labels["redeployedByVaultSync"] = redeployDate
			statefulsetObject.Spec.UpdateStrategy.Type = "RollingUpdate"
			_, err = k.clientset.AppsV1().StatefulSets(k.Namespace).Update(statefulsetObject)
			return err
		})
		return retryErr
	}
	return nil
}
