# DNS Interceptor
DNS interceptor creates a mock dns server that intercepts dns requests and injects chaos on the provided settings

## Usage
```shell
TARGET_PID=39590 TARGET_HOSTNAMES='["google","fb.com"]' CHAOS_TYPE=error CHAOS_DURATION=5 MATCH_SCHEME=substring ./dns_interceptor
#or
TARGET_PID=39590 SPOOF_MAP='{"google.com":"fakegoogle.com"}' CHAOS_TYPE=spoof CHAOS_DURATION=5 ./dns_interceptor
```
### Env Vars
1. TARGET_PID : dns_interceptor can be paired with nsutil to run interceptor in the a different ns, this env var points to the target pid whose ns will be used.
2. CHAOS_TYPE: specifies the type of chaos to run, can be `error` or `spoof`
3. TARGET_HOSTNAMES: list of target host names to intercept (applicable for `error` chaos)
4. MATCH_SCHEME: there are 2 types of match schemes : `exact` and `substring`, this var determines whether the dns query has to match exactly with one of the targets or can have any of the targets as substring (applicable for `error` chaos)
5. SPOOF_MAP: map of host names, where the key specifies the target host name while the value is the host names to which the query is to be spoofed (applicable for `spoof` chaos)  
6. CHAOS_DURATION: time period in seconds during which the interceptor will run
