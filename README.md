# premetheus-alert-webhook
Send prometheus alert using webhook. Customize the alert in webhook server and then send a mail.

-----------------------------
# AlertManager
Download and run alertManager manually. <br>
https://prometheus.io/download/ <br>
https://prometheus.io/docs/alerting/latest/alertmanager/

<strike>
Or you can use <i> docker compose </i>

```
version: '3'

services:
  alertmanager:
    image: prom/alertmanager:v0.23.0
    restart: unless-stopped
    ports:
      - "9093:9093"
    volumes:
      - "./alertmanager:/webhookconfig.yml"
    #command: --config.file=/config/alertmanager.yml --log.level=debug
```
</strike>

To run(manually) and send alert in a webhook server : <br>
`./alertmanager --config.file=webhook.yml`
------------------------------
# webhook.yml
````
global:
  resolve_timeout: 10s
route:
  receiver: webhook_receiver
  group_wait: 0s
  group_interval: 10s
  repeat_interval: 4h
receivers:
    - name: webhook_receiver
      webhook_configs:
        - url: http://localhost:8080/webhook
          send_resolved: false
````

---------------------------------

# send alerts manually using the alertmanager-api
````
curl -X POST -H "Content-Type: application/json" -d\
'[{
 "status":"firing",
 "annotations": {
     "property1": "string1",
     "property2": "string2"
  },
  "labels": {
     "label1": "string3",
     "label2": "string4"
   },
   "generatorURL": "http://example.com"
}]' http://localhost:9093/api/v2/alerts
````

-------------------- 

# Up to now
Prometheus send alert to webhook. Webhook print last 10 lines log for current running pod and then send an email 
parsing the logs.


--------------------

# TODO
 1. Proper log parsing

----------------------
# Issue
  1. Prometheus is vulnerable with multiple replicas
----------------------

# Bookstore api
https://github.com/shn27/BookStoreApi-Go <br>
https://github.com/samiulsami/GolangBookstoreAPI

Run it using helm. If you want to create an original alert. Then run prometheus and add prometheus rule.
https://prometheus.io/docs/prometheus/latest/installation/