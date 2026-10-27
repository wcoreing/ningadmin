package entry

import (
	"wnadmin/app/common/response"

	"github.com/gin-gonic/gin"
)

func Handle(c *gin.Context) {
	response.OkWithDetailed(
		"success", "success", c)

}
