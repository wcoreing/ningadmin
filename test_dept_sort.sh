#!/bin/bash

echo "=== 测试部门排序更新功能 ==="

# 先获取部门列表
echo "1. 获取部门树..."
curl -X GET "http://localhost:8080/api/system/dept/tree" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq .

echo -e "\n2. 测试更新部门排序..."
# 更新部门排序（使用完整的部门信息）
curl -X PUT "http://localhost:8080/api/system/dept" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "deptId": 100,
    "parentId": 0,
    "deptName": "总公司",
    "orderNum": 999,
    "leader": "管理员",
    "phone": "",
    "email": "",
    "status": "0"
  }' | jq .

echo -e "\n3. 测试更新部门状态..."
# 测试状态更新接口
curl -X PUT "http://localhost:8080/api/system/dept/status" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "deptId": 100,
    "status": "1"
  }' | jq .

echo -e "\n测试完成！" 