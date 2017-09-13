import socket
import sys
from time import ctime

path = "/tmp"
errorLog = open(path + "/stderr.txt", "w", 1)
sys.stderr = errorLog
stdoutLog = open(path + "/stdout.txt", "w", 1)
sys.stdout = stdoutLog

# create an ipv4 (AF_INET) socket object using the tcp protocol (SOCK_STREAM)
client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

# connect the client
# client.connect((target, port))
client.connect(('127.0.0.1', 9999))

# send some data (in this case a HTTP GET request)
count = 0
while count < 100:
    count += 1
    print count

    client.send('GET /index.html')

    # receive the response data (4096 is recommended buffer size)
    response = client.recv(4096)

    print "(" + str(count) + ") " + ctime() + " : " + response + "$"
