#!/bin/bash

go install github.com/smart-echo/micro-toolkit/cmd/micro@latest

source <(micro completion bash)

micro new service github.com/smart-echo/micro-demo/hello && cd hello && make init proto update tidy && cd -

git add . && git commit -m "add hello service" && git push