#!/bin/bash

go install github.com/smart-echo/micro-toolkit/cmd/micro@latest

source <(micro completion bash)

rm -rf hello && hello-client

micro new service --trace github.com/smart-echo/micro-demo/hello && cd hello && make init proto update tidy && cd -

git add . && git commit -m "upgrade hello service, support trace" && git push

micro new client github.com/smart-echo/micro-demo/hello && cd hello-client && make update tidy && cd -

git add . && git commit -m "upgrade hello client" && git push