package entry

import (
	"github.com/gin-gonic/gin"
	"wnadmin/app/common/response"
	"wnadmin/app/model/do/StudentActiveDB"
	"wnadmin/app/provider/mylog"
)

func Handle(c *gin.Context) {
	sysUser := StudentActiveDB.SysUserDO{}
	sysUser.Query().First(&sysUser)
	mylog.Log.Info("response", sysUser)
	response.OkWithDetailed(
		sysUser, "success", c)

}
