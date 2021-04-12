# DNS Interceptor
DNS interceptor creates a mock dns server that intercepts dns requests and injects chaos on the provided settings

## Usage
```shell
TARGET_PID=39590 TARGET_HOSTNAMES='["google","fb.com"]' CHAOS_DURATION=5 MATCH_SCHEME=substring ./dns_interceptor
```
### Env Vars
1. TARGET_PID : dns_interceptor can be paired with nsutil to run interceptor in the a different ns, this env var points to the target pid whose ns will be used.
2. TARGET_HOSTNAMES: list of target host names to intercept
3. MATCH_SCHEME: there are 2 types of match schemes : exact and substring, this var determines whether the dns query has to match exactly with one of the targets or can have any of the targets as substring
4. CHAOS_DURATION: time period in seconds during which the interceptor will run
