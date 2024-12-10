package model

type DatasetStatus int

const (
	DatasetCreating DatasetStatus = iota + 1
	DatasetCreated
	DatasetFailedCreate
	DatasetBound
	DatasetDeleting
	DatasetDeleted
	DatasetFailedDelete
)

func (f DatasetStatus) String() string {
	switch f {
	case DatasetCreating:
		return "creating"
	case DatasetCreated:
		return "created"
	case DatasetFailedCreate:
		return "failed create"
	case DatasetBound:
		return "bound"
	case DatasetDeleting:
		return "deleting"
	case DatasetDeleted:
		return "deleted"
	case DatasetFailedDelete:
		return "failed delete"
	default:
		return "invalid status"
	}
}
