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
	"log"
	"os"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Config contains the kubernetes clientset and configs.
type Config struct {
	clientset        *kubernetes.Clientset
	secretObject     *apiv1.Secret
	DeploymentList   []string
	StatefulsetList  []string
	SecretName       string
	Namespace        string
	KubeSecretExists bool
}

// authenticate authenticates and authorizes the client.
func (k *Config) authenticate() (err error) {
	var config *rest.Config
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if len(kubeconfigPath) > 0 {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		log.Printf("Cannot connect to Kubernetes: %v\n", err)
		return err
	}
	k.clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
