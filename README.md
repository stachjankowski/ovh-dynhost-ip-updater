ovh-dynhost-ip-updater
=====

Lightweight Go tool to update IP address in [OVH DynHost](https://docs.ovh.com/gb/en/domains/hosting_dynhost/).

## Configuration
1. Configure your domain in [OVH DynHost](https://docs.ovh.com/gb/en/domains/hosting_dynhost/)
2. Create your app keys https://docs.ovh.com/gb/en/api/first-steps-with-ovh-api/#create-your-app-keys
3. Save app keys. Because we using [go-ovh](https://github.com/ovh/go-ovh) lib you have two options file config or evoiremenets. All information you have [here](https://github.com/ovh/go-ovh#configuration).


## Installation and run
### Docker

```
docker run -e ENDPOINT=ovh-eu -e APPLICATION_KEY=... -e APPLICATION_SECRET=... -e CUSTOMER_KEY=... stachjankowski/ovh-dynhost-ip-updater:latest --zone example.com --subdomain abc
```

### Standalone
```
make
./bin/ovh-dynhost-ip-updater --zone example.com --subdomain abc
```

The program can be run once or run in a loop (every 60s). To run in a loop, use the switch  ```--loop```.

You have two options with IP address transfer.
* Manual: e.g.: ```--ip 127.0.0.3```
* Detect public IP using external service, e.g.:
  * ```--ipurl https://ip.42.pl/raw```
  * ```--ipurl "http://ip-api.com/json/" --jsonpath "\$.query"```
  
### Help:
```$ ./bin/ovh-dynhost-ip-updater -h
Usage: ovh-dynhost-ip-updater --zone ZONE --subdomain SUBDOMAIN [--ipurl IPURL] [--jsonpath JSONPATH] [--ip IP] [--loop] [--verbose]

Options:
  --zone ZONE            domain
  --subdomain SUBDOMAIN
                         subdomain
  --ipurl IPURL          address to the service that returns your public ip address [default: https://api.ipify.org]
  --jsonpath JSONPATH
  --ip IP                IPv4 address
  --loop, -l             work in an infinite loop
  --verbose, -v          verbosity level
  --help, -h             display this help and exit```