package infrastructure

import "fmt"

// FileSystemError represents an error that occurs during file I/O operations
type FileSystemError struct {
	Operation string // The operation being performed (e.g., "read", "write", "create directory")
	Path      string // The file path involved
	Err       error  // The underlying error
}

func (e FileSystemError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("filesystem error during %s on '%s': %v", e.Operation, e.Path, e.Err)
	}
	return fmt.Sprintf("filesystem error during %s: %v", e.Operation, e.Err)
}

// Unwrap returns the underlying error for error unwrapping
func (e FileSystemError) Unwrap() error {
	return e.Err
}

// NewFileSystemError creates a new FileSystemError
func NewFileSystemError(operation, path string, err error) FileSystemError {
	return FileSystemError{
		Operation: operation,
		Path:      path,
		Err:       err,
	}
}

// WrapFileSystemError wraps an OS error with context about the file operation
func WrapFileSystemError(operation, path string, err error) error {
	if err == nil {
		return nil
	}
	return NewFileSystemError(operation, path, err)
}
