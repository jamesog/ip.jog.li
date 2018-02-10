// +build !appengine

package ipjogli

import (
	"errors"
	"net/http"
)

func mcGet(r *http.Request, ip string) (ipInfo, error) {
	return ipInfo{}, errors.New("not implemented")
}

func mcSet(r *http.Request, ip ipInfo) error {
	return errors.New("not implemented")
}
