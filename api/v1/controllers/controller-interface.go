package controllers

// ControllerDescriber -
type ControllerDescriber interface {
	Connect(connectionString string) (setupError error)
	Close() error
}
