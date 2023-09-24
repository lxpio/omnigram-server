package user

import (
	"strconv"

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

func getUserInfoHandle(c *gin.Context) {

	userID := c.GetInt64(middleware.XUserIDTag)

	key := userKeyPrefix + strconv.FormatInt(userID, 10)

	if user, ok := userInfoCache.Get(key); ok {
		c.JSON(200, utils.SUCCESS.WithData(user))
		return
	}

	user, err := schema.FirstUserByID(orm, userID)

	if err != nil {
		log.E(`获取用户信息失败：`, err.Error())
		c.JSON(404, utils.ErrGetUserInfo)
		return
	}

	userInfoCache.Add(key, user)

	c.JSON(200, utils.SUCCESS.WithData(user))
}

func createAPIKeyHandle(c *gin.Context) {

	userID := c.GetInt64(middleware.XUserIDTag)

	token := schema.NewAPIToken(userID)

	if err := token.Save(orm); err != nil {
		log.E(`创建APIKey失败：`, err.Error())
		c.JSON(500, utils.ErrSaveToken)
		return
	}

	c.JSON(200, utils.SUCCESS.WithData(token.APIKey))

}

// DELETE /user/apikeys/:id
func deleteAPIKeyHandle(c *gin.Context) {

	id := c.Param(`id`)

	if err := schema.DeleteAPIKey(orm, id); err != nil {
		log.E(`删除APIKey失败：`, err.Error())
		c.JSON(500, utils.ErrDeleteToken)
		return
	}

}

func getAPIKeysHandle(c *gin.Context) {

	userID := c.GetInt64(middleware.XUserIDTag)

	log.D(`userID`, userID)

	keys, err := schema.GetAPIKeysByUserID(orm, userID)

	if err != nil {
		log.E(`获取APIKey失败：`, err.Error())
		c.JSON(500, utils.ErrGetTokens)
		return
	}

	c.JSON(200, utils.SUCCESS.WithData(keys))

}

func createAccountHandle(c *gin.Context) {
	panic(`TODO`)
}
func listAccountHandle(c *gin.Context) {
	panic(`TODO`)
}
func getAccountHandle(c *gin.Context) {
	panic(`TODO`)
}
func deleteAccountHandle(c *gin.Context) {
	panic(`TODO`)
}
