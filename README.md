[![Build Status](https://travis-ci.org/zoer/optima.svg)](https://travis-ci.org/zoer/optima)

# Installation
```bash
make cold_start
```

# Usage
```bash
$ curl -X POST http://$(docker-machine ip default 2> /dev/null || echo '127.0.0.1'):8089/configs -d '{"Type": "Test.vpn", "Data": "Rabbit.log"}'
{"host": "10.0.5.42", "port": "5671", "virtualhost": "/", "user": "guest", "password": "guest"}
```
