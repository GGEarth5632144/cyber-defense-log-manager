import socket
import time

UDP_IP = "127.0.0.1"
UDP_PORT = 514

messages = [
    # 4.1 Firewall/Syslog
    "<134>Aug 20 12:44:56 fw01 vendor=demo product=ngfw action=deny src=10.0.1.10 dst=8.8.8.8 spt=5353 dpt=53 proto=udp msg=DNS blocked policy=Block-DNS",
    "<134>Aug 20 12:45:10 fw01 vendor=demo product=ngfw action=deny src=10.0.1.10 dst=8.8.8.9 spt=5353 dpt=53 proto=udp msg=DNS blocked",
    "<134>Aug 20 12:45:20 fw01 vendor=demo product=ngfw action=deny src=10.0.1.10 dst=8.8.8.10 spt=5353 dpt=53 proto=udp msg=DNS blocked",
    
    # 4.2 Network (Router Syslog)
    "<190>Aug 20 13:01:02 r1 if=ge-0/0/1 event=link-down mac=aa:bb:cc:dd:ee:ff reason=carrier-loss"
]

print(f"Sending syslog messages to {UDP_IP}:{UDP_PORT}...")
sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

for msg in messages:
    sock.sendto(msg.encode('utf-8'), (UDP_IP, UDP_PORT))
    print(f"Sent: {msg}")
    time.sleep(0.5)

print("Done.")
