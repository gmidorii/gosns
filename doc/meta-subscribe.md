# API doc

This is API documentation for Subscribe API. This is generated by `httpdoc`. Don't edit by hand.

## Table of contents

- [[200] POST /meta/subscribe](#200-post-metasubscribe)


## [200] POST /meta/subscribe

Register topic subscribed

### Request



Headers

| Name  | Value  | Description |
| ----- | :----- | :--------- |
| Content-Type | application/json |  |





Request example

```

{
  "channel": "/meta/subscribe",
  "client_id": "hogehoge",
  "subscription" : [
  	"/golang"
  ],
  "method" : {
    "format": "slack",
    "webhook_url": "https://hooks.slack.com/services/XXXXX"
  }
}

```


### Response

Headers

| Name  | Value  | Description |
| ----- | :----- | :--------- |
| Content-Type | application/json |  |





Response example

```
{"channel":"/meta/subscribe","successful":true,"clientId":"hogehoge","subscription":["/golang"]}
```

