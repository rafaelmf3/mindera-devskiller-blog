package bootstrap

import (
	"bitbucket.org/mindera/go-rest-blog/service"
)

func Init(port int) error {
	api := service.NewRestApiService()
	return api.ServeContent(port)
}
