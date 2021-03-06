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
	"time"

	"github.com/thatInfrastructureGuy/VaultSync/pkg/common/data"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

// createSecretObject creates a Kubernetes Secret Object in memory.
func (k *Config) createSecretObject() {
	k.secretObject = &apiv1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      k.SecretName,
			Namespace: k.Namespace,
		},
		Type: "Opaque",
	}
}

// secretUpdater updates secrets into specified Kubernetes Secret
// If secret name not specified; secret with same name as vault is created.
// Errors out if namespace is not present.
func (k *Config) secretUpdater() error {
	namespace := k.secretObject.GetNamespace()
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Deployment before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		_, err := k.clientset.CoreV1().Secrets(namespace).Update(k.secretObject)
		return err
	})

	return retryErr
}

// secretCreator creates secrets into specified Kubernetes Secret
// If secret name not specified; secret with same name as vault is created.
// Errors out if namespace is not present.
func (k *Config) secretCreator() error {
	namespace := k.secretObject.GetNamespace()
	_, err := k.clientset.CoreV1().Secrets(namespace).Create(k.secretObject)
	if err != nil {
		log.Println("Error creating secret: ", err)
		return err
	}
	return nil
}

// secretsUpdater is an internal wrapper over kube secret operations.
func (k *Config) secretsUpdater(secretList map[string]data.SecretAttribute) error {
	// Instantiate secret data
	if len(k.secretObject.Data) == 0 {
		k.secretObject.Data = make(map[string][]byte)
	}
	for secretKey, secretAttributes := range secretList {
		k.secretObject.Data[secretKey] = []byte(secretAttributes.Value)
		//Delete the key if set to empty
		if secretAttributes.Value == "" || secretAttributes.MarkedForDeletion {
			delete(k.secretObject.Data, secretKey)
		}
	}
	// Set the date updated timestamp
	annotations := k.secretObject.GetAnnotations()
	if len(annotations) == 0 {
		annotations = make(map[string]string)
	}
	annotations["dateUpdated"] = time.Now().Format(time.RFC3339)
	k.secretObject.SetAnnotations(annotations)

	if !k.KubeSecretExists {
		return k.secretCreator()
	}
	return k.secretUpdater()
}

func (k *Config) getSecretObject() {
	var err error
	// Get the secret object
	k.secretObject, err = k.clientset.CoreV1().Secrets(k.Namespace).Get(k.SecretName, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		k.KubeSecretExists = false

		// Create kube secret empty object
		k.createSecretObject()
		return
	}
	k.KubeSecretExists = true
}

func (k *Config) GetLastUpdatedDate() (date time.Time, err error) {
	err = k.authenticate()
	if err != nil {
		return date, err
	}
	k.getSecretObject()
	if !k.KubeSecretExists {
		return date, nil
	}
	annotations := k.secretObject.GetAnnotations()
	value, ok := annotations["dateUpdated"]
	if !ok {
		return date, nil
	}
	date, err = time.Parse(time.RFC3339, value)
	if err != nil {
		return date, err
	}
	return date, nil
}

// PostSecrets is common interface function to post secrets to destination
func (k *Config) PostSecrets(secretList map[string]data.SecretAttribute) (err error) {
	err = k.secretsUpdater(secretList)
	if err != nil {
		return err
	}
	return nil
}
