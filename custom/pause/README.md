# Pause

Pause is a CLI tool that can be used to begin a command in paused state. It usually helps in performing some operations before the command is actually executed. 
For example, actions like loading `cgroup`.

## Usage

```sh
mv pause /usr/local/bin
pause nsutil -p -n -t 39590 -- stress-ng -c 2 -t 60s
```
