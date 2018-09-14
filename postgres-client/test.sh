#!/usr/bin/expect -f
 
set ns $::env(NAMESPACE)
set sv $::env(SERVICE_NAME)
set db $::env(DATABASE_NAME)
set password $::env(PASSWORD)
set user $::env(DATABASE_USER)
set port $::env(PORT)
set parallel $::env(PARALLEL_TRANSACTION)
set transaction $::env(TRANSACTIONS)
set timeout -1   
spawn pgbench -i -h $sv.$ns.svc.cluster.local -p $port  -U $user  -s 30 $db
# Look for passwod prompt
expect "Password:"
# Send password aka $password 
send "$password\r"
# send blank line (\r) to make sure we get back to gui
send -- "\r"
expect "#"

#set i [lindex $argv 0];
set i 1
while {$i > 0 } {
#puts "count : $i\n";   #this is for printing th value of variable
spawn pgbench -c 4 -h $sv.$ns.svc.cluster.local -p $port -U $user  -j $parallel -t $transaction $db
# Look for passwod prompt
expect "Password"
# Send password aka $password 
send "$password\r"
# send blank line (\r) to make sure we get back to gui
send -- "\r"
expect "#"
#set i [expr $i-1];
}