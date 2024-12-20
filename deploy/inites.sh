#!/bin/sh
set -e
echo "============init template==============" 
HTTP_CODE=$(curl  -s -o ./resp.md -w "%{http_code}" -L -XPUT -k --user elastic:$1 -H "Content-Type: application/json"  https://$2:9200/_template/events_template -d '
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
				"format": "strict_date_optional_time",
				"index": false
			},
			"data": {
				"properties": {
					"task_id": {
						"type": "keyword"
					},
					"task_record_id": {
						"type": "keyword",
						"index": false
					},
					"task_name": {
						"type": "keyword",
						"index": false
					},
					"account_id": {
						"type": "keyword",
						"index": false
					},
					"user_id": {
						"type": "keyword"
					},
					"compute_type": {
						"type": "keyword",
						"index": false
					},
					"node_ip": {
						"type": "ip",
						"index": false
					},
					"node_name": {
						"type": "keyword"
					},
					"node_uuid": {
						"type": "keyword",
						"index": false
					},
					"pod_namespace": {
						"type": "keyword",
						"index": false
					},
					"pod_ip": {
						"type": "ip",
						"index": false
					},
					"pod_name": {
						"type": "keyword",
						"index": false
					},
					"region_id": {
						"type": "keyword"
					},
					"resource_group_id": {
						"type": "keyword"
					},
					"resource_group_name": {
						"type": "keyword",
						"index": false
					},
					"level": {
						"type": "keyword"
					},
					"status": {
						"type": "keyword",
						"index": false
					},
					"event_message": {
            			"type": "text",
            			"analyzer": "ik_smart"
					},
					"localguid": {
						"type": "keyword",
						"index": false
					},
					"errcode": {
						"type": "keyword",
						"index": false
					},
					"workspace_name": {
						"type": "keyword",
						"index": false
					},
					"workspace_id": {
						"type": "keyword"
					},
					"event_time": {
						"type": "date",
						"format": "epoch_second"
					}
				}
			}
		}
	}
}')

if [ "$HTTP_CODE" -eq 200 ]; then
  cat resp.md && rm -rf resp.md
  echo "============init template succeed" 
else 
  echo "============init template failed"
fi

echo "=============init ilm================="
HTTP_CODE=$(curl -s -o ./resp.md -w "%{http_code}" -L -XPUT -k --user elastic:$1 -H "Content-Type: application/json" https://$2:9200/_ilm/policy/hot_delete  -d '
{
  "policy": {
    "phases": {
      "hot": {
            "min_age": "0ms",
            "actions": {
              "set_priority": {"priority": 100},
              "rollover": {
                "max_age": "30d",
                "max_size": "50gb"
              }
            }
          },
      "warm": {
        "min_age": "30d",
        "actions": {
          "forcemerge": {"max_num_segments": 1},
          "set_priority": {"priority": 50}
        }
      },
      "delete": {
        "actions": {
          "delete": {"delete_searchable_snapshot": true}
        }
      }
    }
  }
}')
if [ "$HTTP_CODE" -eq 200 ]; then
  cat resp.md && rm -rf resp.md
  echo "============init ilm succeed" 
else 
  echo "============init ilm failed"
fi

echo "============create index and alias=============="
HTTP_CODE=$(curl  -s -o ./resp.md -w "%{http_code}" -L -XPUT -k --user elastic:$1 -H "Content-Type: application/json"  https://$2:9200/events-00001 -d '{"aliases":{"yunxiao-events":{"is_write_index":true}}}')

if [ "$HTTP_CODE" -eq 200 ]; then
  cat resp.md && rm -rf resp.md
  echo "============create index and alias succeed"
else
  echo "============create index and alias failed"
fi