# premetheus-alert-webhook
Send prometheus alert using webhook. Customize the alert in webhook server and then send a mail.

# send alerts manually using the alertmanager-api
````
curl -X POST -H "Content-Type: application/json" -d\
'[{
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
Prometheus send alert to webhook. Webhook print last 10 lines log for current running pod. 


--------------------

# TODO
 1. Proper log parsing
 2. Creating and Sending mail(alert)

----------------------
# Issue
  1.Prometheus is vulnerable with multiple replicas
