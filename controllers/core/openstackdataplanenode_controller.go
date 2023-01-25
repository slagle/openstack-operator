/*
Copyright 2022.

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

package core

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	corev1beta1 "github.com/openstack-k8s-operators/openstack-operator/apis/core/v1beta1"
)

// OpenStackDataPlaneNodeReconciler reconciles a OpenStackDataPlaneNode object
type OpenStackDataPlaneNodeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

//+kubebuilder:rbac:groups=core.openstack.org,resources=openstackdataplanenodes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core.openstack.org,resources=openstackdataplanenodes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core.openstack.org,resources=openstackdataplanenodes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the OpenStackDataPlaneNode object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *OpenStackDataPlaneNodeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Fetch the OpenStackDataPlaneNode instance
	instance := &corev1beta1.OpenStackDataPlaneNode{}
	err := r.Client.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if k8s_errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected.
			// For additional cleanup logic use finalizers. Return and don't requeue.
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	if instance.Spec.Node.Managed {
		err = r.Provision(ctx, instance)
		if err != nil {
			r.Log.Error(err, fmt.Sprintf("Unable to OpenStackDataPlaneNode %s", instance.Name))
			return ctrl.Result{}, err
		}
	}

	err = r.GenerateInventory(ctx, instance)
	if err != nil {
		r.Log.Error(err, fmt.Sprintf("Unable to generate inventory for %s", instance.Name))
		return ctrl.Result{}, err
	}

	r.ConfigureNetwork(ctx, instance)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OpenStackDataPlaneNodeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1beta1.OpenStackDataPlaneNode{}).
		Complete(r)
}

func (r *OpenStackDataPlaneNodeReconciler) Provision(ctx context.Context, instance *corev1beta1.OpenStackDataPlaneNode) error {
	return nil
}

type Inventory struct {
	all struct {
		hosts struct {
			host struct {
				host_var struct {
					host_var_value string
				}
			}
		}
	}
}

func (r *OpenStackDataPlaneNodeReconciler) GenerateInventory(ctx context.Context, instance *corev1beta1.OpenStackDataPlaneNode) error {
	var err error

	inventory := make(map[string]map[string]map[string]map[string]string)
	all := make(map[string]map[string]map[string]string)
	host := make(map[string]map[string]string)
	host_vars := make(map[string]string)
	host_vars["ansible_host"] = instance.Spec.Node.HostName
	host[instance.Name] = host_vars
	all["hosts"] = host
	inventory["all"] = all

	configMapName := fmt.Sprintf("dataplanenode-%s-inventory", instance.Name)
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: instance.Namespace,
		},
	}

	_, err = controllerutil.CreateOrUpdate(ctx, r.Client, cm, func() error {
		cm.TypeMeta = metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		}
		cm.ObjectMeta = metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: instance.Namespace,
		}
		invData, err := yaml.Marshal(inventory)
		if err != nil {
			return err
		}
		cm.Data = map[string]string{
			"inventory": string(invData),
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *OpenStackDataPlaneNodeReconciler) ConfigureNetwork(ctx context.Context, instance *corev1beta1.OpenStackDataPlaneNode) error {

	return nil
}
