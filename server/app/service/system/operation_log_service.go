package system

import (
	"encoding/json"
	"fmt"
	"time"
	"wnadmin/app/model/do/SystemDB"
	"wnadmin/app/model/dto"
)

type OperationLogService struct{}

// RecordLog 记录操作日志
func (s *OperationLogService) RecordLog(logData *dto.OperationLogData) error {
	// 将请求参数和响应数据转换为JSON字符串
	requestParamsStr := ""
	responseDataStr := ""

	if logData.RequestParams != nil {
		if bytes, err := json.Marshal(logData.RequestParams); err == nil {
			requestParamsStr = string(bytes)
		}
	}

	if logData.ResponseData != nil {
		if bytes, err := json.Marshal(logData.ResponseData); err == nil {
			responseDataStr = string(bytes)
		}
	}

	// 创建操作日志记录
	operationLog := &SystemDB.OperationLogDO{
		TraceID:       logData.TraceID,
		UserID:        logData.UserID,
		Username:      logData.Username,
		UserType:      logData.UserType,
		OperationType: logData.OperationType,
		Module:        logData.Module,
		OperationName: logData.OperationName,
		Method:        logData.Method,
		RequestURL:    logData.RequestURL,
		RequestParams: requestParamsStr,
		ResponseData:  responseDataStr,
		IPAddress:     logData.IPAddress,
		UserAgent:     logData.UserAgent,
		Location:      logData.Location,
		Status:        logData.Status,
		ErrorMsg:      logData.ErrorMsg,
		CostTime:      logData.CostTime,
	}

	return operationLog.Create()
}

// GetLogList 获取操作日志列表
func (s *OperationLogService) GetLogList(req *dto.OperationLogListReq) (*dto.PageResult, error) {
	logDO := &SystemDB.OperationLogDO{}

	// 构建查询条件
	conditions := make(map[string]interface{})
	if req.Username != "" {
		conditions["username"] = req.Username
	}
	if req.Module != "" {
		conditions["module"] = req.Module
	}
	if req.OperationType != "" {
		conditions["operation_type"] = req.OperationType
	}
	if req.StartTime != "" {
		conditions["start_time"] = req.StartTime
	}
	if req.EndTime != "" {
		conditions["end_time"] = req.EndTime
	}

	// 查询数据
	logs, total, err := logDO.GetLogList(req.Page, req.PageSize, conditions)
	if err != nil {
		return nil, err
	}

	return &dto.PageResult{
		List:     logs,
		Total:    total,
		PageNum:  req.Page,
		PageSize: req.PageSize,
	}, nil
}

// GetOperationMapping 获取操作映射配置
func (s *OperationLogService) GetOperationMapping() map[string]dto.OperationInfo {
	return map[string]dto.OperationInfo{
		// 用户管理
		"GET:/system/user/list":      {Module: "用户管理", Name: "查询用户列表", Type: "SELECT"},
		"GET:/system/user/info":      {Module: "用户管理", Name: "查看用户详情", Type: "SELECT"},
		"POST:/system/user":          {Module: "用户管理", Name: "新增用户", Type: "CREATE"},
		"PUT:/system/user":           {Module: "用户管理", Name: "修改用户", Type: "UPDATE"},
		"DELETE:/system/user/delete": {Module: "用户管理", Name: "删除用户", Type: "DELETE"},
		"PUT:/system/user/status":    {Module: "用户管理", Name: "修改用户状态", Type: "UPDATE"},
		"PUT:/system/user/password":  {Module: "用户管理", Name: "重置密码", Type: "UPDATE"},

		// 角色管理
		"GET:/system/role/list":   {Module: "角色管理", Name: "查询角色列表", Type: "SELECT"},
		"GET:/system/role":        {Module: "角色管理", Name: "查看角色详情", Type: "SELECT"},
		"POST:/system/role":       {Module: "角色管理", Name: "新增角色", Type: "CREATE"},
		"PUT:/system/role":        {Module: "角色管理", Name: "修改角色", Type: "UPDATE"},
		"DELETE:/system/role":     {Module: "角色管理", Name: "删除角色", Type: "DELETE"},
		"PUT:/system/role/status": {Module: "角色管理", Name: "修改角色状态", Type: "UPDATE"},

		// 菜单管理
		"GET:/system/menu/list":      {Module: "菜单管理", Name: "查询菜单列表", Type: "SELECT"},
		"GET:/system/menu/tree":      {Module: "菜单管理", Name: "查询菜单树", Type: "SELECT"},
		"GET:/system/menu/info":      {Module: "菜单管理", Name: "查看菜单详情", Type: "SELECT"},
		"POST:/system/menu/add":      {Module: "菜单管理", Name: "新增菜单", Type: "CREATE"},
		"PUT:/system/menu/update":    {Module: "菜单管理", Name: "修改菜单", Type: "UPDATE"},
		"DELETE:/system/menu/delete": {Module: "菜单管理", Name: "删除菜单", Type: "DELETE"},
		"PUT:/system/menu/status":    {Module: "菜单管理", Name: "修改菜单状态", Type: "UPDATE"},
		"PUT:/system/menu/sort":      {Module: "菜单管理", Name: "修改菜单排序", Type: "UPDATE"},

		// 部门管理
		"GET:/system/dept/list":   {Module: "部门管理", Name: "查询部门列表", Type: "SELECT"},
		"GET:/system/dept/tree":   {Module: "部门管理", Name: "查询部门树", Type: "SELECT"},
		"GET:/system/dept":        {Module: "部门管理", Name: "查看部门详情", Type: "SELECT"},
		"POST:/system/dept":       {Module: "部门管理", Name: "新增部门", Type: "CREATE"},
		"PUT:/system/dept":        {Module: "部门管理", Name: "修改部门", Type: "UPDATE"},
		"DELETE:/system/dept":     {Module: "部门管理", Name: "删除部门", Type: "DELETE"},
		"PUT:/system/dept/status": {Module: "部门管理", Name: "修改部门状态", Type: "UPDATE"},

		// 系统配置
		"GET:/system/config/list":      {Module: "系统配置", Name: "查询配置列表", Type: "SELECT"},
		"GET:/system/config/tree":      {Module: "系统配置", Name: "查询配置树", Type: "SELECT"},
		"GET:/system/config/info":      {Module: "系统配置", Name: "查看配置详情", Type: "SELECT"},
		"POST:/system/config/add":      {Module: "系统配置", Name: "新增配置", Type: "CREATE"},
		"PUT:/system/config/update":    {Module: "系统配置", Name: "修改配置", Type: "UPDATE"},
		"DELETE:/system/config/delete": {Module: "系统配置", Name: "删除配置", Type: "DELETE"},

		// 登录相关
		"POST:/auth/login":  {Module: "系统登录", Name: "用户登录", Type: "LOGIN"},
		"POST:/auth/logout": {Module: "系统登录", Name: "用户登出", Type: "LOGOUT"},
	}
}

// GetStatistics 获取操作日志统计
func (s *OperationLogService) GetStatistics(days int) (map[string]interface{}, error) {
	logDO := &SystemDB.OperationLogDO{}

	// 统计总数
	var totalCount int64
	if err := logDO.Query().Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// 统计指定天数内的数据
	var recentCount int64
	startTime := fmt.Sprintf("%s 00:00:00", time.Now().AddDate(0, 0, -days).Format("2006-01-02"))
	if err := logDO.Query().Where("create_time >= ?", startTime).Count(&recentCount).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"totalCount":  totalCount,
		"recentCount": recentCount,
	}, nil
}

// DeleteLog 删除操作日志
func (s *OperationLogService) DeleteLog(id int64) error {
	logDO := &SystemDB.OperationLogDO{}
	return logDO.Query().Where("id = ?", id).Delete(logDO).Error
}

// BatchDeleteLogs 批量删除操作日志
func (s *OperationLogService) BatchDeleteLogs(ids []int64) error {
	logDO := &SystemDB.OperationLogDO{}
	return logDO.Query().Where("id IN ?", ids).Delete(logDO).Error
}

// ClearLogs 清空指定天数之前的日志
func (s *OperationLogService) ClearLogs(days int) (int64, error) {
	logDO := &SystemDB.OperationLogDO{}
	cutoffTime := time.Now().AddDate(0, 0, -days).Format("2006-01-02 15:04:05")

	result := logDO.Query().Where("create_time < ?", cutoffTime).Delete(logDO)
	return result.RowsAffected, result.Error
}
