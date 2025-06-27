#!/bin/bash

echo "=== 测试分销功能修复 ==="

echo "1. 检查数据库表..."
mysql -h 127.0.0.1 -uroot -p123456 -e "USE wnsys; SHOW TABLES LIKE 'share_%'; SHOW TABLES LIKE 'user_relation';" 2>/dev/null

echo -e "\n2. 测试分销收益接口..."
curl -X GET "http://localhost:8080/api/distribution/income/list" \
  -H "Content-Type: application/json" 2>/dev/null | jq . || echo "接口返回: $(curl -X GET "http://localhost:8080/api/distribution/income/list" -H "Content-Type: application/json" 2>/dev/null)"

echo -e "\n3. 测试分销申请接口..."
curl -X GET "http://localhost:8080/api/distribution/apply/list" \
  -H "Content-Type: application/json" 2>/dev/null | jq . || echo "接口返回: $(curl -X GET "http://localhost:8080/api/distribution/apply/list" -H "Content-Type: application/json" 2>/dev/null)"

echo -e "\n4. 测试分销提现接口..."
curl -X GET "http://localhost:8080/api/distribution/withdraw/list" \
  -H "Content-Type: application/json" 2>/dev/null | jq . || echo "接口返回: $(curl -X GET "http://localhost:8080/api/distribution/withdraw/list" -H "Content-Type: application/json" 2>/dev/null)"

echo -e "\n5. 检查后端服务状态..."
if curl -s "http://localhost:8080/" > /dev/null; then
    echo "✅ 后端服务正常运行"
else
    echo "❌ 后端服务未运行"
fi

echo -e "\n测试完成！" 