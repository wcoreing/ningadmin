#!/bin/bash

# 测试用户列表接口
echo "测试用户列表接口..."
curl -X GET "http://localhost:8080/api/system/user/list?pageNum=1&pageSize=10" \
  -H "Authorization: Bearer your_token_here" \
  | jq '.'

echo -e "\n\n测试按部门搜索..."
curl -X GET "http://localhost:8080/api/system/user/list?pageNum=1&pageSize=10&deptId=100" \
  -H "Authorization: Bearer your_token_here" \
  | jq '.' 