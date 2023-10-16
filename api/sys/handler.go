package sys

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nexptr/omnigram-server/api/epub"
	"github.com/nexptr/omnigram-server/api/epub/selfhost"
	"github.com/nexptr/omnigram-server/conf"
	"github.com/nexptr/omnigram-server/utils"
)

type SysInfo struct {
	Version     string              `json:"version"`
	ChatEnabled bool                `json:"chat_enabled,omitempty"`
	M4tEnabled  bool                `json:"m4t_enabled,omitempty"`
	ScanStatus  selfhost.ScanStatus `json:"scan_stats"`
}

// getSysInfoHandle get User Authorization
/**
 * @api {get} /sys/info Get Current Server Info
 * @apiName getSysInfoHandle
 * @apiGroup sys
 * @apiDescription Get server configs. if chat server has configed, or the
 * m4t service is support.
 *
 * @apiHeader {String} Authorization Users unique auth key.
 *
 * @apiSuccess {Boolean} chatserver     Always set to Bearer.
 * @apiSuccess {Number} expires_in     Number of seconds that the included access token is valid for.
 * @apiSuccess {String} refresh_token  Issued if the original scope parameter included offline_access.
 * @apiSuccess {String} access_token   Issued for the scopes that were requested.
 */
func getSysInfoHandle(c *gin.Context) {

	mng := epub.GetManager()

	info := &SysInfo{
		Version:     conf.Version,
		ChatEnabled: true,
		M4tEnabled:  true,
		ScanStatus:  mng.Status(),
	}

	c.JSON(http.StatusOK, utils.SUCCESS.WithData(info))

}
