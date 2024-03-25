# premetheus-alert-webhook
Send prometheus alert using webhook. Customize the alert in webhook server and then send a mail.

-----------------------------
# Run webhook using helm chart
`git clone git@github.com:searchlight/prometheusWebHook.git` <br>
`cd prometheusWebHook` <br>
`helm install webhook webhookHelm/`

------------------------------
# AlertManager
Download and run alertManager manually. <br>
https://prometheus.io/download/ <br>
https://prometheus.io/docs/alerting/latest/alertmanager/

Or you can use <i> docker compose </i> with:
``
docker compose up
``

```
version: '3'

services:
  alertmanager:
    image: prom/alertmanager:v0.23.0
    restart: unless-stopped
    volumes:
      - "./webHookConfig.yml:/webHookConfig.yml"
    network_mode: "host"
    command: --config.file=/webHookConfig.yml --log.level=debug
```

To run(manually) and send alert in a webhook server : <br>
`./alertmanager --config.file=webHookConfig.yml`
------------------------------
# webhook.yml
````
global:
  resolve_timeout: 5m
route:
  receiver: webhook_receiver
  group_wait: 0s
  group_interval: 10s
  repeat_interval: 4h
receivers:
    - name: webhook_receiver
      webhook_configs:
        - url: http://172.24.0.2:30000/webhook
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
  2. Jmap does not show any error msg if failed to send mail 
----------------------

# Bookstore api
https://github.com/shn27/BookStoreApi-Go <br>
https://github.com/samiulsami/GolangBookstoreAPI

Run it using helm. If you want to create an original alert. Then run prometheus and add prometheus rule.
https://prometheus.io/docs/prometheus/latest/installation/