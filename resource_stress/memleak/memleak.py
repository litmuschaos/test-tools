#***Max-Out Memory**
giant_string = ''
while True:
    with open('long_file.txt', 'r') as f:
        for line in f:
            giant_string += line
