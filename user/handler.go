package user

import (
	"github.com/gin-gonic/gin"
	"github.com/nexptr/omnigram-server/log"
	"github.com/nexptr/omnigram-server/middleware"
	"github.com/nexptr/omnigram-server/user/schema"
	"github.com/nexptr/omnigram-server/utils"
)

func loginHandle(c *gin.Context) {
}

func logoutHandle(c *gin.Context) {

}

func createAPIKeyHandle(c *gin.Context) {

	userID := c.GetInt64(middleware.XUserIDTag)

	token := schema.NewAPIToken(userID)

	if err := token.Save(orm.DB); err != nil {
		log.E(`创建APIKey失败：`, err.Error())
		c.JSON(500, utils.ErrSaveToken)
		return
	}

	c.JSON(200, utils.SUCCESS.WithData(token.APIKey))

}

// DELETE /user/apikeys/:id
func deleteAPIKeyHandle(c *gin.Context) {

	id := c.Param(`id`)

	if err := schema.DeleteAPIKey(orm.DB, id); err != nil {
		log.E(`删除APIKey失败：`, err.Error())
		c.JSON(500, utils.ErrDeleteToken)
		return
	}

}

func getAPIKeysHandle(c *gin.Context) {

	userID := c.GetInt64(middleware.XUserIDTag)

	keys, err := schema.GetAPIKeysByUserID(orm.DB, userID)

	if err != nil {
		log.E(`获取APIKey失败：`, err.Error())
		c.JSON(500, utils.ErrGetTokens)
		return
	}

	c.JSON(200, utils.SUCCESS.WithData(keys))

}
