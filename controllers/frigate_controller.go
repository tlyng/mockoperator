/*
Copyright 2020 Morkel.

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

package controllers

import (
	"context"
	"fmt"
	"log"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	shipv1beta1 "github.com/tlyng/mockoperator/api/v1beta1"
	"github.com/tlyng/mockoperator/component"
)

// FrigateReconciler reconciles a Frigate object
type FrigateReconciler struct {
	client.Client
	Log         logr.Logger
	Scheme      *runtime.Scheme
	Manipulator component.Manipulator
}

// Reconcile ...
// +kubebuilder:rbac:groups=ship.example.com,resources=frigates,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ship.example.com,resources=frigates/status,verbs=get;update;patch
func (r *FrigateReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	logger := r.Log.WithValues("frigate", req.NamespacedName)

	log.Println("Reconciling", req.NamespacedName)
	logger.Info("Reconciling")
	frigate := &shipv1beta1.Frigate{}
	if err := r.Client.Get(ctx, req.NamespacedName, frigate); err != nil {
		return ctrl.Result{}, fmt.Errorf("could not find resource")
	}
	frigate.Status.Foo = r.Manipulator.Manipulate(frigate.Spec.Foo)
	if err := r.Status().Update(context.Background(), frigate); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager ...
func (r *FrigateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&shipv1beta1.Frigate{}).
		Complete(r)
}
