// +build appengine

package ipjogli

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

func mcGet(r *http.Request, ip string) (ipInfo, error) {
	ctx := appengine.NewContext(r)
	var info *ipInfo
	i, err := memcache.JSON.Get(ctx, ip, &info)
	if err != nil {
		return ipInfo{}, err
	}
	if err == nil {
		log.Infof(ctx, "%s", i.Value)
	}
	log.Debugf(ctx, "%+v", err)
	return *info, err
}

func mcSet(r *http.Request, ip ipInfo) error {
	ctx := appengine.NewContext(r)
	i := &memcache.Item{
		Key:    ip.IP,
		Object: ip,
	}
	return memcache.JSON.Set(ctx, i)
}