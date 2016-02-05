package common

type Output interface {
	Run(in []chan *Field) error
	Close() error
	SaveSeriesList([]string) error
	SaveFields(string) error
}
