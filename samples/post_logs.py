import urllib.request
import json
import time

url = "http://localhost:8080/ingest"

# Function to dynamically set current time to avoid retention deletion in demo
def get_now():
    return time.strftime('%Y-%m-%dT%H:%M:%SZ', time.gmtime())

samples = [
    # 4.3 HTTP API
    {
        "tenant": "demoA",
        "source": "api",
        "event_type": "app_login_failed",
        "user": "alice",
        "ip": "203.0.113.7",
        "reason": "wrong_password",
        "@timestamp": get_now()
    },
    # 4.4 CrowdStrike
    {
        "tenant": "demoA",
        "source": "crowdstrike",
        "event_type": "malware_detected",
        "host": "WIN10-01",
        "process": "powershell.exe",
        "severity": 8,
        "sha256": "abc...",
        "action": "quarantine",
        "@timestamp": get_now()
    },
    # 4.5 AWS CloudTrail
    {
        "tenant": "demoB",
        "source": "aws",
        "cloud": {"service": "iam", "account_id": "123456789012", "region": "ap-southeast-1"},
        "event_type": "CreateUser",
        "user": "admin",
        "@timestamp": get_now(),
        "raw": {"eventName": "CreateUser", "requestParameters": {"userName": "temp-user"}}
    },
    # 4.6 Microsoft 365 Audit
    {
        "tenant": "demoB",
        "source": "m365",
        "event_type": "UserLoggedIn",
        "user": "bob@demo.local",
        "ip": "198.51.100.23",
        "status": "Success",
        "workload": "Exchange",
        "@timestamp": get_now()
    },
    # 4.7 Microsoft AD/Windows Security
    {
        "tenant": "demoA",
        "source": "ad",
        "event_id": 4625,
        "event_type": "LogonFailed",
        "user": "demo\\eve",
        "host": "DC01",
        "ip": "203.0.113.77",
        "logon_type": 3,
        "@timestamp": get_now()
    }
]

for sample in samples:
    req = urllib.request.Request(url, method="POST")
    req.add_header('Content-Type', 'application/json')
    data = json.dumps(sample).encode('utf-8')

    print(f"Sending {sample['source']} log...")
    try:
        with urllib.request.urlopen(req, data=data) as response:
            print(response.read().decode('utf-8'))
    except Exception as e:
        print("Failed:", e)
    time.sleep(0.5)
