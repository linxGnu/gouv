import socket
import sys

path = "/tmp"
errorLog = open(path + "/stderr2.txt", "w", 1)
errorLog.write("---Starting Error Log---\n")
sys.stderr = errorLog
stdoutLog = open(path + "/stdout2.txt", "w", 1)
stdoutLog.write("---Starting Standard Out Log---\n")
sys.stdout = stdoutLog

# create an ipv4 (AF_INET) socket object using the tcp protocol (SOCK_STREAM)
client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

# connect the client
# client.connect((target, port))
client.connect(('127.0.0.1', 10000))

# send some data (in this case a HTTP GET request)
while True:
    client.send('GET /index.html')

    # receive the response data (4096 is recommended buffer size)
    response = client.recv(4096)

    print response
