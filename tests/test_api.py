import unittest
import urllib.request
import json
import time

BASE_URL = "http://localhost:8080"

class TestLogManagerAPI(unittest.TestCase):
    def setUp(self):
        # Authenticate
        req = urllib.request.Request(f"{BASE_URL}/login", method="POST")
        req.add_header('Content-Type', 'application/json')
        data = json.dumps({"username": "admin", "password": "admin123"}).encode('utf-8')
        try:
            with urllib.request.urlopen(req, data=data) as res:
                self.token = json.loads(res.read().decode('utf-8'))['token']
        except Exception:
            self.skipTest("API is not running")

    def test_1_ingest_json(self): #case 1
        sample = {
            "tenant": "demoA",
            "source": "api_test",
            "event_type": "test_event",
            "@timestamp": time.strftime('%Y-%m-%dT%H:%M:%SZ', time.gmtime())
        }
        req = urllib.request.Request(f"{BASE_URL}/ingest", method="POST")
        req.add_header('Content-Type', 'application/json')
        data = json.dumps(sample).encode('utf-8')
        with urllib.request.urlopen(req, data=data) as res:
            self.assertEqual(res.status, 200)

    def test_2_fetch_stats(self):#case 2
        req = urllib.request.Request(f"{BASE_URL}/api/stats?time=24h&tenant=all")
        req.add_header('Authorization', f'Bearer {self.token}')
        with urllib.request.urlopen(req) as res:
            self.assertEqual(res.status, 200)
            data = json.loads(res.read().decode('utf-8'))
            self.assertTrue('total_logs' in data)

    def test_3_fetch_logs(self):#case 3
        req = urllib.request.Request(f"{BASE_URL}/api/logs")
        req.add_header('Authorization', f'Bearer {self.token}')
        with urllib.request.urlopen(req) as res:
            self.assertEqual(res.status, 200)
            data = json.loads(res.read().decode('utf-8'))
            self.assertIsInstance(data, list)

if __name__ == '__main__':
    unittest.main()
