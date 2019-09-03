import re

def newBackupCommand(args):
    if args.schedule != None:
        command = "velero schedule create"
        schedule_name = args.name
    else:
        command = "velero backup create"
        backup_name = args.name

    # dictionary vars(args)
    command = command + ' ' + args.name

    args_dict = vars(args)
    for key in args_dict:
        if args_dict[key] != None and key != 'name' and key != 'volns' and key != 'app':
            keyhyphen =  re.sub(r'[^a-zA-Z]','-',key)
            value = str(args_dict[key]).strip('\'[]\'')
            if key == 'volume_snapshot_locations':
                command = command + ' --snapshot-volumes' + ' --' + keyhyphen  + ' ' + value
            else:
                command = command + ' --' + keyhyphen  + ' ' + value
            
    return command