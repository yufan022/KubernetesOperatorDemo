/*

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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1 "yufan.info/m/v2/api/v1"
)

// ApplicationDemoReconciler reconciles a ApplicationDemo object
type ApplicationDemoReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps.yufan.com,resources=applicationdemoes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.yufan.com,resources=applicationdemoes/status,verbs=get;update;patch

func (r *ApplicationDemoReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("applicationdemo", req.NamespacedName)

	// your logic here
	// 1. Print Spec.Detail and Status.Created in log
	obj := &appsv1.ApplicationDemo{}
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		_ = fmt.Errorf("couldn't find object:%s", req.String())
	} else {
		//打印Detail和Created
		log.V(1).Info("Successfully get detail", "Detail", obj.Spec.Detail)
		log.V(1).Info("", "Created", obj.Status.Created)
	}
	// 2. Change Created
	if !obj.Status.Created {
		obj.Status.Created = true
		_ = r.Update(ctx, obj)
	}
	obj = &appsv1.ApplicationDemo{}
	obj.Status.Created = false
	obj.Spec.Detail = "detail"
	obj.Spec.Foo = "foo"
	obj.Namespace = "operator-test"
	if err := r.Create(ctx, obj); err != nil {
		log.Info(err.Error())
	} else {
		log.Info("create success")
	}

	return ctrl.Result{}, nil
}

func (r *ApplicationDemoReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.ApplicationDemo{}).
		Complete(r)
}
