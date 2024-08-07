/*
Copyright 2024.

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

package controller

import (
	"context"
	"github.com/xiaobai0310/service-bind-ingress-operator/utils"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	servicebindingressv1alpha1 "github.com/xiaobai0310/service-bind-ingress-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ServiceBindingIngressReconciler reconciles a ServiceBindingIngress object
type ServiceBindingIngressReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=service-bind-ingress.bailu.io,resources=servicebindingingresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=service-bind-ingress.bailu.io,resources=servicebindingingresses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=service-bind-ingress.bailu.io,resources=servicebindingingresses/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:default:enable_ingres=false
// 简单的检验，可以使用CRD的scheme校验

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ServiceBindingIngress object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *ServiceBindingIngressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// TODO(user): your logic here

	// Fetch the ServiceBindingIngress instance
	instance := &servicebindingressv1alpha1.ServiceBindingIngress{}
	// Get the ServiceBindingIngress instance from cache
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Define a new Deployment object
	deployment := utils.NewDeployment(instance)
	if err := controllerutil.SetControllerReference(instance, deployment, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	d := &appv1.Deployment{}
	if err := r.Get(ctx, req.NamespacedName, d); err != nil {
		if errors.IsNotFound(err) {
			if err := r.Create(ctx, deployment); err != nil {
				logger.Error(err, "failed to create deployment.")
				return ctrl.Result{}, err
			}
		}
	} else {
		if err := r.Update(ctx, deployment); err != nil {
			return ctrl.Result{}, err
		}
	}

	// Define a new Service object
	service := utils.NewService(instance)
	if err := controllerutil.SetControllerReference(instance, service, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	s := &corev1.Service{}
	if err := r.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, s); err != nil {
		if errors.IsNotFound(err) && instance.Spec.EnableService {
			if err := r.Create(ctx, service); err != nil {
				logger.Error(err, "failed to create service.")
				return ctrl.Result{}, err
			}
		}
		if !errors.IsNotFound(err) && instance.Spec.EnableService {
			return ctrl.Result{}, err
		}
	} else {
		if instance.Spec.EnableService {
			logger.Info("skip update service.")
		} else {
			if err := r.Delete(ctx, service); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	// Define a new Ingress object
	// TODO  进行校验，如果启动了Ingress，则必须启动Service；设置默认值，默认值为false
	if instance.Spec.EnableIngress {
		return ctrl.Result{}, nil
	}
	ingress := utils.NewIngress(instance)
	if err := controllerutil.SetControllerReference(instance, ingress, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	i := &networkv1.Ingress{}
	if err := r.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, i); err != nil {
		if errors.IsNotFound(err) && instance.Spec.EnableIngress {
			if err := r.Create(ctx, ingress); err != nil {
				logger.Error(err, "failed to create ingress.")
				return ctrl.Result{}, err
			}
		}
		if !errors.IsNotFound(err) && instance.Spec.EnableIngress {
			return ctrl.Result{}, err
		}
	} else {
		if instance.Spec.EnableIngress {
			logger.Info("skip update ingress.")
		} else {
			if err := r.Delete(ctx, ingress); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ServiceBindingIngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&servicebindingressv1alpha1.ServiceBindingIngress{}).
		Owns(&appv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&networkv1.Ingress{}).
		Complete(r)
}
