package v1alpha1

const (
	// ReplicaIndexLabel represents the label key for the replica-index, e.g. 0, 1, 2.. etc
	ReplicaIndexLabel = "cwai.ctyun.cn/replica-index"

	// ReplicaTypeLabel represents the label key for the replica-type, e.g. ps, worker etc.
	ReplicaTypeLabel = "cwai.ctyun.cn/replica-type"

	// OperatorNameLabel represents the label key for the operator name, e.g. tf-operator, mpi-operator, etc.
	OperatorNameLabel = "cwai.ctyun.cn/operator-name"

	// AppNameLabel represents the label key for the app name, the value is the app name.
	AppNameLabel = "cwai.ctyun.cn/app-name"

	ServerPortKey = "server"

	MetricsPortKey = "metrics"

	DefaultServerContainerName = "default"

	DefaultIngressClass = "nginx"

	DefaultIngressPort = 8080
)
