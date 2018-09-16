package tecutils

import (
	"fmt"
	"net/url"
)

//Devuelve la url base de una direccion http o https
func ParseBaseUrl(urlAddress string) (parsed string, err error) {
	res, err := url.Parse(urlAddress)
	if err != nil {
		return
	}
	parsed = fmt.Sprintf("%s://%s", res.Scheme, res.Host)
	return
}
