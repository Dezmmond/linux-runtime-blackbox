#!/usr/bin/env python3
import socket
import time

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.settimeout(3)

try:
    s.connect(("example.com", 80))
    s.sendall(b"GET / HTTP/1.0\r\nHost: example.com\r\n\r\n")
    time.sleep(5)
finally:
    s.close()