package server

import (
	"github.com/gin-gonic/gin"
	"github.com/cloudposse/openvpn-api/src/config"
	log "github.com/Sirupsen/logrus"
	"github.com/cloudposse/openvpn-api/src/api"
)

// Run - start http server
func Run(cfg config.Config) {

	router := gin.Default()

	router.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")

		logger := log.WithFields(log.Fields{"class": "RootCmd", "method": "RunE"})

		err := api.EnsureUserCerts(name)
		if err != nil {
			logger.Errorf("Ensure user certs failed: %v", err.Error())
			c.String(404, "Can not get user certs")
		}

		clientConf, err := api.GetClientConfig(name)

		if err == nil {
			c.String(200, "%v", clientConf)
		} else {
			logger.Errorf("Get user client configuration failed: %v", err.Error())
			c.String(404, "")
		}
	})

	router.DELETE("/user/:name", func(c *gin.Context) {
		name := c.Param("name")


		err := api.RevokeUser(name)
		if err == nil {
			c.String(200, "Access for user %v revoked", name)
		} else {
			c.String(404, "Error revoking access for user %v", name)
		}
	})

	router.Run(cfg.Listen)
}
