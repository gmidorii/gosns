# gosns

## Overview
gosns is messaging server like Amazon SNS.  
Using messaging model is Pub-Sub messaging model.

## Install
```sh
% go get -u "github.com/midorigreen/gosns"
```

## Build
```sh
# build
% go build

# run
% ./gosns
```

## API Doc
```
doc
├── meta-handshake.md
└── meta-subscribe.md
```

## TODO
- Increasing Subsriver accesepting method
  - [ ] Mail
  - [x] Slack
  - [ ] HTTP
  - etc
- [ ] Add PubSubModel
- [x] Add Topic API
