package lookup

//go:generate moq -stub -out locator_mock.go . Locator

// Locator defines the interface for locating files on a system.
type Locator interface {
	Locate(string) ([]string, error)
}
