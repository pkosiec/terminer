package shared

// Operation is type of installer recipe operation
type Operation string

const (
	// OperationInstall is a recipe install operation
	OperationInstall Operation = "installation"

	// OperationRollback is a recipe rollback operation
	OperationRollback Operation = "rollback"
)
