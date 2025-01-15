package StudentActiveDB

import (
	"gorm.io/gorm"
	"time"
	"wnadmin/app/provider/db"
)

// SysUser 对应 sys_user 表的结构体
type SysUserDO struct {
	UserID        uint64    `gorm:"primaryKey;autoIncrement;column:user_id;comment:用户ID"`
	DeptID        uint64    `gorm:"column:dept_id;comment:部门ID"`
	Code          string    `gorm:"column:code;type:varchar(50);comment:erp工号"`
	UserName      string    `gorm:"column:user_name;type:varchar(30);not null;comment:用户账号"`
	NickName      string    `gorm:"column:nick_name;type:varchar(30);not null;comment:用户昵称"`
	UserType      string    `gorm:"column:user_type;type:varchar(2);default:'00';comment:用户类型（00系统用户）"`
	Email         string    `gorm:"column:email;type:varchar(50);default:'';comment:用户邮箱"`
	Phonenumber   string    `gorm:"column:phonenumber;type:varchar(11);default:'';comment:手机号码"`
	Sex           string    `gorm:"column:sex;type:char(1);default:'0';comment:用户性别（0男 1女 2未知）"`
	Avatar        string    `gorm:"column:avatar;type:varchar(100);default:'';comment:头像地址"`
	Password      string    `gorm:"column:password;type:varchar(100);default:'';comment:密码"`
	Status        string    `gorm:"column:status;type:char(1);default:'0';comment:帐号状态（0正常 1停用）"`
	DelFlag       string    `gorm:"column:del_flag;type:char(1);default:'0';comment:删除标志（0代表存在 2代表删除）"`
	LoginIP       string    `gorm:"column:login_ip;type:varchar(128);default:'';comment:最后登录IP"`
	LoginDate     time.Time `gorm:"column:login_date;comment:最后登录时间"`
	CreateBy      string    `gorm:"column:create_by;type:varchar(64);default:'';comment:创建者"`
	CreateTime    time.Time `gorm:"column:create_time;comment:创建时间"`
	UpdateBy      string    `gorm:"column:update_by;type:varchar(64);default:'';comment:更新者"`
	UpdateTime    time.Time `gorm:"column:update_time;comment:更新时间"`
	Remark        string    `gorm:"column:remark;type:varchar(500);comment:备注"`
	PwdUpdateDate time.Time `gorm:"column:pwd_update_date;comment:密码最后更新时间"`
}

func (t *SysUserDO) Query() *gorm.DB {
	return db.DB.WnAdmin.Table("sys_user")
}
