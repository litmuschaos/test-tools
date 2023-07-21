module github.com/litmuschaos/dns_interceptor

go 1.20

require (
	github.com/miekg/dns v1.1.41
	github.com/sirupsen/logrus v1.9.3
)

require (
	golang.org/x/net v0.0.0-20220906165146-f3363e06e74c // indirect
	golang.org/x/sys v0.5.0 // indirect
)

replace golang.org/x/net => golang.org/x/net v0.7.0
