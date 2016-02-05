package common

type Output interface {
	Run(in []chan *Field) error
	// Gracefully shutdown writer
	Shutdown() error
	SaveSeriesList([]string) error
	SaveFields(string) error
}
