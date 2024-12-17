#!/bin/sh
set -e

echo "============send event==============\n"

# 定义基础 JSON 模板
json_template='{
  "specversion": "2.0",
  "id": "__UNIQUE_ID__",
  "source": "ctyun.yunxiao_taskinfo",
  "ctyunregion": "cn-beijing",
  "type": "__TYPE__",
  "datacontenttype": "application/json",
  "subject": "ctyun.yunxiao_task.cn-beijing.1234567890.task_record_id:__TASK_RECORD_ID__;task_pod:pod123;",
  "time": "2024-12-10T15:09:00Z",
  "data": {
    "task_id": "__TASK_ID__",
    "task_record_id": "__TASK_RECORD_ID__",
    "task_name": "Test Task",
    "task_detail": "This is a test task.",
    "account_id": "1234567890",
    "user_id": "user1234",
    "region_id": "cn-beijing",
    "resource_group_id": "rg1234",
    "resource_group_name": "Resource Group 1",
    "workspace_name": "Workspace 1",
    "workspace_id": "ws1234",
    "level": "Critical",
    "event_time": "2024-12-10T15:09:00Z",
    "status": "failed",
    "status_message": "The task has failed due to an error.",
    "event_message": "The task has failed due to an error."
  }
}'

# 检查 uuidgen 是否存在
if ! command -v uuidgen >/dev/null; then
  echo "Error: UUID generation tool (uuidgen) is not installed."
  exit 1
fi

# 捕获终止信号
trap 'echo "Script terminated."; exit' SIGTERM SIGINT

# 初始化 task_id 和 task_record_id 的起始值
counter=1

# 定义可用的 type 值
type_values=("Critical" "Warning" "Info")

# 无限循环，每 2 秒执行一次
while true; do
  unique_id=$(uuidgen)
  current_time=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

  # 随机选择一个 type 值
  type_index=$(( RANDOM % 3 )) # 0 到 2 的随机数
  selected_type=${type_values[$type_index]}

  # 动态替换占位符
  json_body=$(echo "$json_template" | sed "s/__UNIQUE_ID__/$unique_id/" \
                                      | sed "s/__TASK_ID__/$counter/" \
                                      | sed "s/__TASK_RECORD_ID__/$counter/" \
                                      | sed "s/__TYPE__/$selected_type/" \
                                      | sed "s/2024-12-10T15:09:00Z/$current_time/")

  # 使用 curl 执行请求
  HTTP_CODE=$(curl -s -o ./resp.md -w "%{http_code}" -L -XPOST -k \
    -H "Content-Type: application/json" \
    -H "EventToken: xwViJm3InjSCXAhNxyTNsdI" \
    -d "$json_body" \
    http://10.233.87.172:9200/apis/v1/event-agent/events)

  if [ "$HTTP_CODE" -eq 200 ]; then
    echo "\n============send event succeeded\n"
    cat resp.md
    rm -rf resp.md
  else
    echo "\n============send event failed, HTTP_CODE: $HTTP_CODE\n"
    [ -f resp.md ] && cat resp.md
  fi

  # 每次循环后递增 counter
  counter=$((counter + 1))

  # 等待 2 秒
  sleep 2
done