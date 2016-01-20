package output

type Output interface {
	SaveSeriesList([]string) error
}
