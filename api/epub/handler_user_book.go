package epub

import (
	"github.com/gin-gonic/gin"
	"github.com/nexptr/omnigram-server/api/epub/schema"
	"github.com/nexptr/omnigram-server/log"
	"github.com/nexptr/omnigram-server/middleware"
	"github.com/nexptr/omnigram-server/utils"
)

// PersonalBooksHandle 获取用户喜欢的书籍列表
/**
 * @api {get} /book/fav Get User Favorite Books
 * @apiName FavBookHandle
 * @apiGroup book
 * @apiDescription Get Personal liked , reading, and marked books.
 *
 * @apiHeader {String} Authorization Users unique auth key.
 *
 * @apiSuccess {Boolean} chatserver     Always set to Bearer.
 * @apiSuccess {Number} expires_in     Number of seconds that the included access token is valid for.
 * @apiSuccess {String} refresh_token  Issued if the original scope parameter included offline_access.
 * @apiSuccess {String} access_token   Issued for the scopes that were requested.
 */
func FavBookHandle(c *gin.Context) {

	userID := c.GetInt64(middleware.XUserIDTag)

	req := &struct {
		Limit  int `form:"limit"`
		Offset int `form:"offset"`
	}{10, 0}

	if err := c.ShouldBind(req); err != nil {
		log.I(`请求参数异常`, err)
		c.JSON(400, utils.ErrReqArgs.WithMessage(err.Error()))
		return
	}

	likes, err := schema.LikedBooks(orm, userID, req.Offset, req.Limit)

	if err != nil {
		log.I(`用户登录参数异常`, err)
		c.JSON(200, utils.ErrInnerServer.WithMessage(err.Error()))
		return
	}

	c.JSON(200, utils.SUCCESS.WithData(likes))
}

// PersonalBooksHandle 获取用户喜欢的书籍列表
/**
 * @api {get} /book/personal Get User Personal Books
 * @apiName PersonalBooksHandle
 * @apiGroup book
 * @apiDescription Get Personal liked , reading, and marked books.
 *
 * @apiHeader {String} Authorization Users unique auth key.
 *
 * @apiSuccess {Boolean} chatserver     Always set to Bearer.
 * @apiSuccess {Number} expires_in     Number of seconds that the included access token is valid for.
 * @apiSuccess {String} refresh_token  Issued if the original scope parameter included offline_access.
 * @apiSuccess {String} access_token   Issued for the scopes that were requested.
 */
func PersonalBooksHandle(c *gin.Context) {

	userID := c.GetInt64(middleware.XUserIDTag)

	readings, err := schema.ReadingBooks(orm, userID, 0, 20)

	if err != nil {
		log.I(`用户登录参数异常`, err)
		c.JSON(200, utils.ErrInnerServer.WithMessage(err.Error()))
		return
	}

	likes, err := schema.LikedBooks(orm, userID, 0, 20)

	if err != nil {
		log.I(`用户登录参数异常`, err)
		c.JSON(200, utils.ErrInnerServer.WithMessage(err.Error()))
		return
	}

	c.JSON(200, utils.SUCCESS.WithData(map[string][]schema.ProcessBook{
		"readings": readings,
		"likes":    likes,
	}))

}
