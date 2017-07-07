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

### Docker
```sh
# build
% docker build -t gosns:1.0 ./

# run (port 8080 -> 8888)
% docker run -it -p 8080:8888 gosns:1.0
```

## API Doc
```
doc
├── meta-handshake.md
└── meta-subscribe.md
```

