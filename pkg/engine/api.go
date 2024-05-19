package engine

// Interface defines the API for a runtime config updater.
type Interface interface {
	DefaultRuntime() string
	AddRuntime(string, string, bool) error
	RemoveRuntime(string) error
	Save(string) (int64, error)
}
