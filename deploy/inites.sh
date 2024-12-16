#!/bin/sh
set -e
echo "============init template==============\n" 
HTTP_CODE=$(curl  -s -o ./resp.md -w "%{http_code}" -L -XPUT -k --user elastic:$1 -H "Content-Type: application/json"  http://$2:9200/_template/events_template -d '
{
	"index_patterns": "events-*",
	"order": 1,
	"settings": {
		"number_of_shards": 3,
		"number_of_replicas": 1,
		"index.lifecycle.name": "hot_delete",
		"index.lifecycle.rollover_alias": "yunxiao-events"
	},
	"mappings": {
		"properties": {
			"specversion": {
				"type": "keyword"
			},
			"id": {
				"type": "keyword"
			},
			"source": {
				"type": "keyword"
			},
			"ctyunregion": {
				"type": "keyword"
			},
			"type": {
				"type": "keyword"
			},
			"datacontenttype": {
				"type": "keyword"
			},
			"subject": {
				"type": "keyword"
			},
			"time": {
				"type": "date",
				"format": "strict_date_optional_time||yyyy-MM-dd HH:mm:ss||yyyy-MM-dd"
			},
			"data": {
				"properties": {
					"task_id": {
						"type": "keyword"
					},
					"task_record_id": {
						"type": "keyword"
					},
					"task_name": {
						"type": "text"
					},
					"account_id": {
						"type": "keyword"
					},
					"user_id": {
						"type": "keyword"
					},
					"compute_type": {
						"type": "keyword"
					},
					"node_ip": {
						"type": "ip"
					},
					"node_name": {
						"type": "keyword"
					},
					"pod_namespace": {
						"type": "keyword"
					},
					"pod_ip": {
						"type": "ip"
					},
					"pod_name": {
						"type": "keyword"
					},
					"region_id": {
						"type": "keyword"
					},
					"resource_group_id": {
						"type": "keyword"
					},
					"resource_group_name": {
            "type": "text",
            "analyzer": "ik_smart"
					},
					"level": {
						"type": "keyword"
					},
					"status": {
						"type": "keyword"
					},
					"event_message": {
            "type": "text",
            "analyzer": "ik_smart"
					},
					"localguid": {
						"type": "keyword"
					},
					"errcode": {
						"type": "keyword"
					},
					"workspace_name": {
            "type": "text",
            "analyzer": "ik_smart"
					},
					"workspace_id": {
						"type": "keyword"
					},
					"event_time": {
						"type": "date",
						"format": "strict_date_optional_time||epoch_second||yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"
					}
				}
			}
		}
	}
}')

if [ "$HTTP_CODE" -eq 200 ]; then
  cat resp.md && rm -rf resp.md
  echo "\n ============init template succeed\n" 
else 
  echo "\n ============init template failed\n"
fi
echo "=============init ilm================= \n"
HTTP_CODE=$(curl -s -o ./resp.md -w "%{http_code}" -L -XPUT -k --user elastic:$1 -H "Content-Type: application/json" http://$2:9200/_ilm/policy/hot_delete  -d '
{
  "policy": {
    "phases": {
      "hot": {
            "min_age": "0ms",
            "actions": {
              "set_priority": {"priority": 100},
              "rollover": {"max_age": "3m"}
            }
          },
      "warm": {
        "min_age": "3m",
        "actions": {
          "forcemerge": {"max_num_segments": 1},
          "set_priority": {"priority": 50}
        }
      },
      "delete": {
        "min_age": "4m",
        "actions": {
          "delete": {"delete_searchable_snapshot": true}
        }
      }
    }
  }
}')
if [ "$HTTP_CODE" -eq 200 ]; then
  cat resp.md && rm -rf resp.md
  echo "============init ilm succeed\n" 
else 
  echo "============init ilm failed\n"
fi
set -e
echo "============create index and alias==============\n"
HTTP_CODE=$(curl  -s -o ./resp.md -w "%{http_code}" -L -XPUT -k --user elastic:$1 -H "Content-Type: application/json"  https://$2:9200/events-00001 -d '{"aliases":{"yunxiao-events":{"is_write_index":true}}}')

if [ "$HTTP_CODE" -eq 200 ]; then
  cat resp.md && rm -rf resp.md
  echo "\n ============create index and alias succeed\n"
else
  echo "\n ============create index and alias failed\n"
fi
set -e

echo "============send event==============\n"

# 定义基础 JSON 模板
json_template='{
  "specversion": "2.0",
  "id": "__UNIQUE_ID__",
  "source": "ctyun.yunxiao_taskinfo",
  "ctyunregion": "cn-beijing",
  "type": "task_failed",
  "datacontenttype": "application/json",
  "subject": "ctyun.yunxiao_task.cn-beijing.1234567890.task_record_id:1234;task_pod:pod123;",
  "time": "2024-12-10T15:09:00Z",
  "data": {
    "task_id": "1234",
    "task_record_id": "1234",
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

# 无限循环，每 2 秒执行一次
while true; do
  unique_id=$(uuidgen)
  current_time=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

  # 动态替换占位符
  json_body=$(echo "$json_template" | sed "s/__UNIQUE_ID__/$unique_id/" | sed "s/2024-12-10T15:09:00Z/$current_time/")

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

  # 等待 2 秒
  sleep 2
done