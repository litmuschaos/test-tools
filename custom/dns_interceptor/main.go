package main

// cgo code to enter mnt ns of target process

/*
#cgo CFLAGS: -Werror -I/usr/include
#cgo LDFLAGS: -L/usr/lib64/ -L/usr/lib/x86_64-linux-gnu/ -ldl
#define _GNU_SOURCE
#include <sched.h>
#include <stdio.h>
#include <fcntl.h>
#include <unistd.h>
#include <stdlib.h>
__attribute__((constructor)) void enter_namespace(void) {
	char *ns = (char*)malloc(30 * sizeof(char));
	sprintf(ns, "/proc/%s/ns/mnt", getenv("TARGET_PID"));
	int fd = open(ns, O_RDONLY);
	int res = setns(fd, 0);
	if(res!=0){
		printf("[ERROR] Could not enter mnt ns of PID : %s\n",getenv("TARGET_PID"));
		exit(1);
	}
	close(fd);
}
*/
import "C"
import (
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/litmuschaos/dns_interceptor/server"
	log "github.com/sirupsen/logrus"
)

const resolvConfPath = "/etc/resolv.conf"

func main() {
	chaosDuration := 2 * time.Minute
	duration := os.Getenv("CHAOS_DURATION")
	if duration != "" {
		seconds, err := strconv.Atoi(duration)
		if err != nil {
			log.WithError(err).Fatal("Invalid chaos duration")
		}
		chaosDuration = time.Duration(seconds) * time.Second
	}
	dnsInterceptor, err := server.NewDNSInterceptor(resolvConfPath)
	if err != nil {
		log.WithError(err).Fatal("Failed to create Interceptor")
	}

	dnsInterceptor.Serve(".")

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	timeChan := time.Tick(chaosDuration)

	select {
	case <-timeChan:
		log.Info("Chaos duration complete, stopping")
	case s := <-sig:
		log.Infof("Signal (%v) received, stopping", s)
	}
	dnsInterceptor.Shutdown()
}
