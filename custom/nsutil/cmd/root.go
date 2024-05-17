package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"
)

// user and mnt ns not supported
var (
	ns         = []string{"net", "pid", "cgroup", "uts", "ipc", "mnt"}
	nsSelected = make([]bool, 7)
	t          int // target pid
	libPath    string
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "nsutil",
	Short: "A cli tool to execute commands in target namespace",
	Long: `A cli tool to execute commands in target namespace, very similar to nsenter. 
This tool also forwards any kill signals to the executed command 
and also pipes the standard input and output from the target command`,
	Run: nsutil,
}

const (
	nsFSMagic   = 0x6e736673
	procFSMagic = 0x9fa0
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&nsSelected[0], "net", "n", false, "network namespace to enter")
	rootCmd.PersistentFlags().BoolVarP(&nsSelected[1], "pid", "p", false, "pid namespace to enter")
	rootCmd.PersistentFlags().BoolVarP(&nsSelected[2], "cgroup", "c", false, "cgroup namespace to enter")
	rootCmd.PersistentFlags().BoolVarP(&nsSelected[3], "uts", "u", false, "uts namespace to enter")
	rootCmd.PersistentFlags().BoolVarP(&nsSelected[4], "ipc", "i", false, "ipc namespace to enter")
	rootCmd.PersistentFlags().BoolVarP(&nsSelected[5], "mnt", "m", false, "mnt namespace to enter")
	rootCmd.PersistentFlags().IntVarP(&t, "target", "t", 0, "target process id (required)")
	rootCmd.PersistentFlags().StringVar(&libPath, "lib-path", "/usr/local/lib/nsutil.so", "shared library path to be preloaded")
	err := rootCmd.MarkPersistentFlagRequired("target")
	if err != nil {
		log.WithError(err).Fatal("Failed to mark required flag")
	}
}

// nsutil handles the task of executing the required process in the given namespaces
func nsutil(cmd *cobra.Command, args []string) {
	nsMap := getNSFiles()

	// target command
	nCmd := exec.Command(args[0], args[1:]...)
	nCmd.Env = os.Environ()
	if nsSelected[5] {
		nCmd.Env = append(nCmd.Env, fmt.Sprintf("LD_PRELOAD=%s", libPath), fmt.Sprintf("MNT_PATH=%s", fmt.Sprintf("/proc/"+strconv.Itoa(t)+"/ns/mnt")))
	}
	nCmd.Stdin = os.Stdin
	nCmd.Stdout = os.Stdout
	nCmd.Stderr = os.Stderr

	sig := make(chan os.Signal)

	// go routine responsible to enter the appropriate ns and execute the required command
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		for n, f := range nsMap {
			if err := unix.Setns(int(f.Fd()), 0); err != nil {
				log.WithError(err).WithField("ns-type", n).Fatal("Failed to setns")
			}
		}
		err := nCmd.Run()
		if err != nil {
			log.WithError(err).WithField("cmd", nCmd.String()).Fatal("Failed to run command")
		}
		nCmd = nil
		// notify main thread that the target command has exited
		sig <- syscall.SIGKILL
	}()

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	s := <-sig
	// kill the target command process when kill signal received
	if nCmd != nil {
		log.WithField("signal", s).Info("Killing target command")
		err := nCmd.Process.Signal(syscall.SIGINT)
		if err != nil {
			log.WithError(err).Fatal("Failed to kill command")
		}
	}

	log.Infof("Signal (%v) received, stopping", s)
}

// getNSFiles generates a map of all the required ns files
func getNSFiles() map[string]*os.File {
	nsMap := map[string]*os.File{}
	for i, n := range ns {
		if !nsSelected[i] || n == "mnt" {
			continue
		}
		file, err := getFileFromNS("/proc/" + strconv.Itoa(t) + "/ns/" + n)
		if err != nil {
			log.WithError(err).Fatal("Failed to get ns file")
		}
		nsMap[n] = file
	}
	return nsMap
}

// getFileFromNS gets the os.File pointer for the required ns
func getFileFromNS(nsPath string) (*os.File, error) {
	stat := syscall.Statfs_t{}
	if err := syscall.Statfs(nsPath, &stat); err != nil {
		return nil, fmt.Errorf("failed to Statfs %q: %v", nsPath, err)
	}

	switch stat.Type {
	case nsFSMagic, procFSMagic:
		break
	default:
		return nil, fmt.Errorf("unknown FS magic on %q: %x", nsPath, stat.Type)
	}

	file, err := os.Open(nsPath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// Execute is the entrypoint for the cli tool
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
