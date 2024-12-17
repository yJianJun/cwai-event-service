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
				"type": "keyword",
				"index": false
			},
			"id": {
				"type": "keyword"
			},
			"source": {
				"type": "keyword",
				"index": false
			},
			"ctyunregion": {
				"type": "keyword",
				"index": false
			},
			"type": {
				"type": "keyword"
			},
			"datacontenttype": {
				"type": "keyword",
				"index": false
			},
			"subject": {
				"type": "keyword",
				"index": false
			},
			"time": {
				"type": "date",
				"format": "date_optional_time||strict_date_optional_time",
				"index": false
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
						"format": "epoch_second||epoch_millis"
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