package operator

import (
	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// param in crd
// AppGroup    = "apps.cwai.ctyun.cn"
// AppVersion  = "v1alpha1"
// AppResource = "cwaionlineapp"
// APPListKind = "CwaiOnlineAppList"
// AppKind     = "CwaiOnlineApp"
// AppPlural   = "cwaionlineapp"
// AppSingular = "cwaionlineapp"
var (
	AppGVR = schema.GroupVersionResource{
		Group:    "apps.cwai.ctyun.cn",
		Version:  "v1alpha1",
		Resource: "cwaionlineapp",
	}
	AppGVK = schema.GroupVersionKind{
		Group:   "apps.cwai.ctyun.cn",
		Version: "v1alpha1",
		Kind:    "CwaiOnlineApp",
	}
)

const (
	Kind       = "CwaiOnlineApp"
	APIVersion = "apps.cwai.ctyun.cn/v1alpha1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CwaiOnlineAppSpec defines the desired state of CwaiOnlineApp
type CwaiOnlineAppSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Replicas is the desired number of replicas of the given template.
	// If unspecified, defaults to 1.
	// +kubebuilder:validation:Minimum=0
	Replicas *int32 `json:"replicas,omitempty"`
	// Template is the object that describes the pod that
	// will be created for this replica.
	Template v1.PodTemplateSpec `json:"template,omitempty"`

	Ingress IngressSpec `json:"ingressSpec,omitempty"`

	// suspend specifies whether the  controller should create Pods or not.
	// If an App is created with suspend set to true, no Pods are created by
	// the  controller. If an App is suspended after creation (i.e. the
	// flag goes from false to true), the  controller will delete all
	// active Pods and PodGroups associated with this App.
	// Users must design their workload to gracefully handle this.
	// Suspending an App will reset the StartTime field of the App.
	//
	// Defaults to false.
	// +kubebuilder:default:=false
	// +optional
	Suspend *bool `json:"suspend,omitempty"`

	// Type is the type of the App.
	// +kubebuilder:validation:Enum=VsCode;Tensorboard;Inference
	Type AppType `json:"type,omitempty"`

	// ConfigMapKey is the key of the configmap. Now is only used for Inference App.
	ConfigMapKey string `json:"configMapKey,omitempty"`

	// SchedulingPolicy defines the policy related to scheduling, e.g. gang-scheduling
	// +optional
	SchedulingPolicy *SchedulingPolicy `json:"schedulingPolicy,omitempty"`

	Strategy appsV1.DeploymentStrategy `json:"strategy,omitempty"`

	Version string `json:"strategy,omitempty"`

	QueueID string `json:"queueID"`
}

// SchedulingPolicy encapsulates various scheduling policies of the distributed training
// job, for example `minAvailable` for gang-scheduling.
type SchedulingPolicy struct {
	MinAvailable           *int32                                 `json:"minAvailable,omitempty"`
	Queue                  string                                 `json:"queue,omitempty"`
	MinResources           *map[v1.ResourceName]resource.Quantity `json:"minResources,omitempty"`
	PriorityClass          string                                 `json:"priorityClass,omitempty"`
	ScheduleTimeoutSeconds *int32                                 `json:"scheduleTimeoutSeconds,omitempty"`
}

type AppType string

const (
	// VsCode is the type of the App.
	VsCode AppType = "VsCode"
	// Tensorboard is the type of the App.
	Tensorboard AppType = "Tensorboard"
	// Inference is the type of the App.
	Inference AppType = "Inference"
)

type IngressSpec struct {
	Annotations map[string]string `json:"annotations,omitempty" protobuf:"bytes,12,rep,name=annotations"`

	// Hostname is the host name of the ingress.
	// +required
	// +kubebuilder:validation:MaxLength=128
	// +kubebuilder:validation:Required
	Hostname string `json:"hostname,omitempty"`
	// Path is the path of the ingress.
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxItems=10
	Path []string `json:"path,omitempty"`

	// IngressClassName is the name of the IngressClass cluster resource.
	// The associated IngressClass defines which controller will implement the
	// resource.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=nginx
	IngressClassName string `json:"ingressClassName,omitempty"`
}

// AppConditionType defines all kinds of types of AppPhaseType.
type AppConditionType string

const (
	// AppCreated means the App has been accepted by the system,
	// but one or more of the pods/services has not been started.
	// This includes time before pods being scheduled and launched.
	AppCreated AppConditionType = "Created"

	// AppRunning means all sub-resources (e.g. services/pods) of this App
	// have been successfully scheduled and launched.
	// The training is running without error.
	AppRunning AppConditionType = "Running"

	// AppRestarting means one or more sub-resources (e.g. services/pods) of this App
	// reached phase failed but maybe restarted according to it's restart policy
	// which specified by user in v1.PodTemplateSpec.
	// The training is freezing/pending.
	AppRestarting AppConditionType = "Restarting"

	// AppSucceeded means all sub-resources (e.g. services/pods) of this App
	// reached phase have terminated in success.
	// The training is complete without error.
	AppSucceeded AppConditionType = "Succeeded"

	// AppSuspended means the App has been suspended.
	AppSuspended AppConditionType = "Suspended"

	// AppFailed means all replicas failed (RestartPolicy=Never)
	AppFailed AppConditionType = "Failed"

	// Unknown
	AppUnknown AppConditionType = "Unknown"
)

// AppCondition describes the state of the App at a certain point.
type AppCondition struct {
	// Type of App condition.
	Type AppConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status v1.ConditionStatus `json:"status"`
	// The reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`
	// A human readable message indicating details about the transition.
	Message string `json:"message,omitempty"`
	// The last time this condition was updated.
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
	// Last time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
}

type AppPhaseType string

const (
	// AppPhaseAvailable means App has at least one replica running
	AppPhaseAvailable AppPhaseType = "Available"
	// AppPhaseCreated means App in reconcile progress
	AppPhaseCreated AppPhaseType = "Created"
	// AppPhaseSuspended  means the App has been suspended.
	AppPhaseSuspended AppPhaseType = "Suspended"
	// AppPhaseFailed means app failed; according to AppFailed ConditionType = true
	AppPhaseFailed AppPhaseType = "Failed"
	// AppPhaseSucceeded means app succeeded; according to AppSucceeded ConditionType = true
	AppPhaseSucceeded AppPhaseType = "Succeeded"
	// AppPhaseUnknown means app Unknown;
	AppPhaseUnknown AppPhaseType = "Unknown"
)

// CwaiOnlineAppStatus defines the observed state of CwaiOnlineApp
type CwaiOnlineAppStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// AppPhaseType means app status
	Phase AppPhaseType `json:"phase,omitempty"`

	// Conditions is an array of current observed App conditions.
	Conditions []AppCondition `json:"conditions,omitempty"`

	// The number of actively running pods.
	Active int32 `json:"active,omitempty"`

	// The number of pods which reached phase Succeeded.
	Succeeded int32 `json:"succeeded,omitempty"`

	// The number of pods which reached phase Failed.
	Failed int32 `json:"failed,omitempty"`

	// The expect number of pods.
	Replicas int32 `json:"replicas,omitempty"`

	// Total number of non-terminated pods targeted by this deployment that have the desired template spec.
	// +optional
	UpdatedReplicas int32 `json:"updatedReplicas,omitempty"`

	// readyReplicas is the number of pods targeted by this Deployment with a Ready Condition.
	// +optional
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`

	// Total number of available pods (ready for at least minReadySeconds) targeted by this deployment.
	// +optional
	AvailableReplicas int32 `json:"availableReplicas,omitempty"`

	// A Selector is a label query over a set of resources. The result of matchLabels and
	// matchExpressions are ANDed. An empty Selector matches all objects. A null
	// Selector matches no objects.
	Selector string `json:"selector,omitempty"`

	// Represents time when the App was acknowledged by the App controller.
	// It is not guaranteed to be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	StartTime *metav1.Time `json:"startTime,omitempty"`

	// Represents time when the App was completed. It is not guaranteed to
	// be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`

	// Represents last time when the App was reconciled. It is not guaranteed to
	// be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	LastReconcileTime *metav1.Time `json:"lastReconcileTime,omitempty"`

	PodStatusList []PodStatus `json:"podStatusList,omitempty"`

	CurrentVersion string `json:"currentVersion,omitempty"`

	QueueID string `json:"queueID"`
}

type PodStatus struct {
	InferenceUUID   string `json:"inferenceID" `
	Version         string `json:"version" `
	PodName         string `json:"podName" `
	PodPort         string `json:"podPort"`
	RestartNum      int32  `json:"restartNum"`
	PodIP           string `json:"podIP"`
	HostIP          string `json:"hostIP"`
	Status          string `json:"status"`
	ContainerName   string `json:"containerName"`
	Namespace       string `json:"namespace"`
	PodCreateAt     string `json:"podCreateAt"`
	RunningTime     string `json:"runningTime"`
	ResourceGroupID string `json:"resourceGroupID"`
	GpuTypeName     string `json:"gpuTypeName"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:path=cwaionlineapp,scope=Namespaced,shortName=app,singular=cwaionlineapp
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="State",type=string,JSONPath=`.status.conditions[-1:].type`
//+kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
// +kubebuilder:subresource:scale:specpath=.spec.replicas,statuspath=.status.active,selectorpath=.status.selector

// CwaiOnlineApp is the Schema for the cwaionlineapps API
type CwaiOnlineApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CwaiOnlineAppSpec   `json:"spec,omitempty"`
	Status CwaiOnlineAppStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=cwaionlineapps
//+kubebuilder:object:root=true

// CwaiOnlineAppList contains a list of CwaiOnlineApp
type CwaiOnlineAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CwaiOnlineApp `json:"items"`
}
