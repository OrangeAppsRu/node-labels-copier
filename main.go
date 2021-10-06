/*
Copyright 2016 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"fmt"
	"time"
	"strings"
	"encoding/json"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/apimachinery/pkg/types"
)

type patchStringValue struct {
    Op    string `json:"op"`
    Path  string `json:"path"`
    Value string `json:"value"`
}

const nameFieldManager = "node-labels-copier"

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		for _, n := range nodes.Items {
			nodeName := n.ObjectMeta.Name
			// Пропускаем ноду. если нет метки "role"
			role := ""
			if _, ok := n.ObjectMeta.Labels["role"]; !ok {
				continue
			}
			role, _ = n.ObjectMeta.Labels["role"]

			nodeRole := ""
			for lKey, _ := range n.ObjectMeta.Labels {
				if strings.Contains(lKey, "node-role.kubernetes.io") {
					split := strings.Split(lKey, "/")
					if len(split) == 2 {
						nodeRole = split[1]
					}
				}
			}

			if role != nodeRole {
				fmt.Printf("copy label for node: %s, role: %s => node-role.kubernetes.io/%s\n", nodeName, role, role)
				payload := []patchStringValue{{
					Op:    "replace",
					Path:  fmt.Sprintf("/metadata/labels/node-role.kubernetes.io~1%s", role),
					Value: "",
				}}
				payloadBytes, _ := json.Marshal(payload)
				
				_ , err = clientset.CoreV1().Nodes().Patch(context.TODO(), nodeName, types.JSONPatchType, payloadBytes, metav1.PatchOptions{FieldManager: nameFieldManager})
				if err != nil {
					fmt.Fprintf(os.Stderr, "%v\n", err)
				}
			}
		}
		time.Sleep(300 * time.Second)
	}
}