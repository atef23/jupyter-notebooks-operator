package jupyternotebooks

import (
	"context"
	"reflect"
	cachev1alpha1 "github.com/atef23/jupyter-notebooks-operator/pkg/apis/cache/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	// route imports
	routev1 "github.com/openshift/api/route/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"github.com/atef23/jupyter-notebooks-operator/pkg/controller/ocp"
	"fmt"
)

var log = logf.Log.WithName("controller_jupyternotebooks")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new JupyterNotebooks Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {

	if err := routev1.AddToScheme(mgr.GetScheme()); err != nil {
		return err
	}

	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileJupyterNotebooks{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("jupyternotebooks-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource JupyterNotebooks
	err = c.Watch(&source.Kind{Type: &cachev1alpha1.JupyterNotebooks{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner JupyterNotebooks
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &cachev1alpha1.JupyterNotebooks{},
	})
	if err != nil {
		return err
	}

	// watch for Route only on OpenShift
	if err = c.Watch(&source.Kind{Type: &routev1.Route{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &cachev1alpha1.JupyterNotebooks{},
	}); err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileJupyterNotebooks implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileJupyterNotebooks{}

// ReconcileJupyterNotebooks reconciles a JupyterNotebooks object
type ReconcileJupyterNotebooks struct {
	// TODO: Clarify the split client
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a JupyterNotebooks object and makes changes based on the state read
// and what is in the JupyterNotebooks.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a JupyterNotebooks Deployment for each JupyterNotebooks CR
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileJupyterNotebooks) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling JupyterNotebooks")

	// Fetch the JupyterNotebooks instance
	jupyterNotebooks := &cachev1alpha1.JupyterNotebooks{}
	err := r.client.Get(context.TODO(), request.NamespacedName, jupyterNotebooks)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			reqLogger.Info("JupyterNotebooks resource not found. Ignoring since object must be deleted")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		reqLogger.Error(err, "Failed to get JupyterNotebooks")
		return reconcile.Result{}, err
	}

	// Check if the deployment already exists, if not create a new one
	found := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: jupyterNotebooks.Name, Namespace: jupyterNotebooks.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {








		// Define a new service
		service := r.serviceForJupyterNotebooks(jupyterNotebooks)
		reqLogger.Info("Creating a new Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
		err = r.client.Create(context.TODO(), service)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
			return reconcile.Result{}, err
		}









		route := &routev1.Route{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "route.openshift.io/v1",
				Kind:       "Route",
			},	
			ObjectMeta: metav1.ObjectMeta{
				Name:      jupyterNotebooks.Name,
				Namespace: jupyterNotebooks.Namespace,
			},
			Spec: routev1.RouteSpec{
				To: routev1.RouteTargetReference{
					Kind: "Service",
					Name: jupyterNotebooks.Name,
				},
				Port: &routev1.RoutePort{
					TargetPort: intstr.FromInt(8888),
				},
			},
		}

		reqLogger.Info("Route defined", route)
		reqLogger.Info("Creating a new Route", "Route.Namespace", route.Namespace, "Route.Name", route.Name)

	
		controllerutil.SetControllerReference(jupyterNotebooks, route, r.scheme)
		if err := r.client.Create(context.TODO(), route); err != nil {
			reqLogger.Error(err, "Failed to create new Route", "Route.Namespace", route.Namespace, "Route.Name", route.Name)
			return reconcile.Result{}, err
		}
	
		/*
		err = r.client.Create(context.TODO(), route)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Route")
			return reconcile.Result{}, err
		}
		*/






		// Define a new deployment
		dep := r.deploymentForJupyterNotebooks(jupyterNotebooks)
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.client.Create(context.TODO(), dep)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return reconcile.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Deployment")
		return reconcile.Result{}, err
	}




		// define a new route
		/*
		route := &routev1.Route{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "route.openshift.io/v1",
				Kind:       "Route",
			},	
			ObjectMeta: metav1.ObjectMeta{
				Name:      jupyterNotebooks.Name,
				Namespace: jupyterNotebooks.Namespace,
			},
			Spec: routev1.RouteSpec{
				TLS: &routev1.TLSConfig{
					InsecureEdgeTerminationPolicy: routev1.InsecureEdgeTerminationPolicyRedirect,
					Termination:                   routev1.TLSTerminationEdge,
				},
				To: routev1.RouteTargetReference{
					Kind: "Service",
					Name: jupyterNotebooks.Name,
				},
				Port: &routev1.RoutePort{
					TargetPort: intstr.FromInt(8888),
				},
			},
		}
		*/



	// Ensure the deployment size is the same as the spec
	size := jupyterNotebooks.Spec.Size
	if *found.Spec.Replicas != size {
		found.Spec.Replicas = &size
		err = r.client.Update(context.TODO(), found)
		if err != nil {
			reqLogger.Error(err, "Failed to update Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return reconcile.Result{}, err
		}
		// Spec updated - return and requeue
		return reconcile.Result{Requeue: true}, nil
	}

	// Update the JupyterNotebooks status with the pod names
	// List the pods for this jupyterNotebooks's deployment
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(jupyterNotebooks.Namespace),
		client.MatchingLabels(labelsForJupyterNotebooks(jupyterNotebooks.Name)),
	}
	if err = r.client.List(context.TODO(), podList, listOpts...); err != nil {
		reqLogger.Error(err, "Failed to list pods", "JupyterNotebooks.Namespace", jupyterNotebooks.Namespace, "JupyterNotebooks.Name", jupyterNotebooks.Name)
		return reconcile.Result{}, err
	}
	podNames := getPodNames(podList.Items)

	// Update status.Nodes if needed
	if !reflect.DeepEqual(podNames, jupyterNotebooks.Status.Nodes) {
		jupyterNotebooks.Status.Nodes = podNames
		err := r.client.Status().Update(context.TODO(), jupyterNotebooks)
		if err != nil {
			reqLogger.Error(err, "Failed to update JupyterNotebooks status")
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

// deploymentForJupyterNotebooks returns a jupyterNotebooks Deployment object
func (r *ReconcileJupyterNotebooks) deploymentForJupyterNotebooks(m *cachev1alpha1.JupyterNotebooks) *appsv1.Deployment {
	ls := labelsForJupyterNotebooks(m.Name)
	replicas := m.Spec.Size

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "jupyterlab",
							Image:   "quay.io/aaziz/jupyterlab:latest",
							Command: []string{"sleep", "3600"},
						},
					},
				},
			},
		},
	}
	// Set JupyterNotebooks instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}

// labelsForJupyterNotebooks returns the labels for selecting the resources
// belonging to the given jupyterNotebooks CR name.
func labelsForJupyterNotebooks(name string) map[string]string {
	return map[string]string{"app": "jupyterNotebooks", "jupyterNotebooks_cr": name}
}

func selectorsForService(name string) map[string]string {
	return map[string]string{
		"app": "jupyterNotebooks",
	}
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}










// serviceForJupyterNotebooks returns a jupyterNotebooks Service object
func (r *ReconcileJupyterNotebooks) serviceForJupyterNotebooks(m *cachev1alpha1.JupyterNotebooks) *corev1.Service {
	selectors := selectorsForService(m.Name)

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
	}

	service.Spec = corev1.ServiceSpec{
		Ports:     jupyterNoteBooksPort.asServicePorts(),
		Selector:  selectors,
	}

	// Set JupyterNotebooks instance as the owner and controller
	controllerutil.SetControllerReference(m, service, r.scheme)
	return service
}









func (r *ReconcileJupyterNotebooks) CreateRoute(m *cachev1alpha1.JupyterNotebooks) *routev1.Route {

	route := ocp.NewRoute(m.Name, m.Namespace, fmt.Sprintf("%s-server", m.Name), 8888)

	// Set JupyterNotebooks instance as the owner and controller
	controllerutil.SetControllerReference(m, route, r.scheme)
	return route
}







/*
func (r *ReconcileJupyterNotebooks) NewRoute(m *cachev1alpha1.JupyterNotebooks, port int) *routev1.Route {
	return &routev1.Route{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: routev1.RouteSpec{
			TLS: &routev1.TLSConfig{
				InsecureEdgeTerminationPolicy: routev1.InsecureEdgeTerminationPolicyRedirect,
				Termination:                   routev1.TLSTerminationEdge,
			},
			To: routev1.RouteTargetReference{
				Kind: "Service",
				Name: m.Name,
			},
			Port: &routev1.RoutePort{
				TargetPort: intstr.FromInt(port),
			},
		},
	}
}
*/