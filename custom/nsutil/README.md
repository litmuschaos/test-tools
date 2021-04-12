# NSUtil

NSUtil is a cli tool that can be used to execute commands in target namespaces, very similar to nsenter. This tool also forwards any kill signals to the executed command
and also pipes the standard input and output from the target command. Currently, this does not support mount and user namespaces.

## Usage

```sh
./nsutil -p -n -t 39590 -- <target-command>
```
### Flags
1. t (target) : `-t <pid>` points to the process id whose namespaces will be used
2. p (pid) : flag to enter pid namespace for the target process
3. n (net) : flag to enter net namespace for the target process
4. c (cgroup): flag to enter cgroup namespace for the target process
5. u (uts) : flag to enter uts namespace for the target process
6. i (ipc) : flag to enter ipc namespace for the target process