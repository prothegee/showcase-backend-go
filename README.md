# showcase backend go

__*table of contents:*__

- [intro](#intro)
- [important](#important)

<br>

- [database](./docs/database/README.md)

<br>

[quic preview](https://youtu.be/TvGujQngAJ0)

[![](https://img.youtube.com/vi/TvGujQngAJ0/hqdefault.jpg)](https://www.youtube.com/watch?v=TvGujQngAJ0)

## intro

a showcase backend using go, postgresql, & redis

__*prequiste:*__

- go: 1.25.0 or higher

<br>

- [__required packages__](./go.mod)

---

## important

__*before anything (build, run, test):*__

1. check on the this root project for config.json:
    - if doesn't exists, copy paste from config.json.template to config.json

2. check:
    - [backend_api listener](./config.json.template:4)
    - [postgresql main db](./config.json.template:11)
    - [redis main db](./config.json.template:21)

3. scripts:
    - [to build](./dbuild.sh)
    - [to debug use dlv](./ddebug.sh)
    - [to run the development](./drun.sh)
    - [to test *required to run the service/s first](./dtest.sh)

<br>

__*to run the test:*__

1. after all those 3 already checked
2. open a terminal, then you can run the service by run [`./drun.sh`](./drun.sh)
3. open another terminal session then run [`./dtest.sh`](./dtest.sh)

---

###### end of readme

