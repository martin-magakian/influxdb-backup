package sqlite

import (
	"net/url"
)

func quoteName(in string) (string) {
	return url.QueryEscape(in)

}

func unqouteName(in string) (string,error) {
	return url.QueryUnescape(in)
}
