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
	char *pid = getenv("TARGET_PID");
	if(pid==NULL){
		printf("[INFO] No PID mentioned running in current ns\n");
		return;
	}
	sprintf(ns, "/proc/%s/ns/mnt", pid);
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
	port := os.Getenv("PORT")
	if duration != "" {
		seconds, err := strconv.Atoi(duration)
		if err != nil {
			log.WithError(err).Fatal("Invalid chaos duration")
		}
		chaosDuration = time.Duration(seconds) * time.Second
	}
	if port == "" {
		port = server.DefaultDNSPort
	} else if _, err := strconv.Atoi(port); err != nil {
		log.WithError(err).Fatal("Invalid port")
	}
	log.WithField("port", port).Info("DNS Interceptor Port")
	dnsInterceptor, err := server.NewDNSInterceptor(resolvConfPath, os.Getenv("UPSTREAM_SERVER"))
	if err != nil {
		log.WithError(err).Fatal("Failed to create Interceptor")
	}

	dnsInterceptor.Serve(".", port)

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
