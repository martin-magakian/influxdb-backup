package common

type Output interface {
	SaveSeriesList([]string) error
	SaveFields(string) error
}
