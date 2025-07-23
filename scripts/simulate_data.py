import requests
import random
import time
import uuid
from datetime import datetime, timezone

API_URL = "http://127.0.0.1:8080/telemetry"  # Change port if needed

def generate_telemetry(trip_id):
    return {
        "lat": round(random.uniform(12.90, 13.10), 6),
        "lon": round(random.uniform(77.50, 77.70), 6),
        "speed": round(random.uniform(20, 80), 2),  # speed in km/h
        "timestamp": datetime.now(timezone.utc).isoformat(),
        "trip_id": trip_id
    }

def simulate_data(num_trips=3, points_per_trip=10, delay=1):
    health_url = "http://127.0.0.1:8080/health"
    resp = requests.get(health_url)
    if resp.status_code != 200:
        print("Health check failed")
        return
    print("Health check passed")
    trip_ids = [str(uuid.uuid4()) for _ in range(num_trips)]
    print(f"Simulating {num_trips} trips with {points_per_trip} data points each...")

    for trip_id in trip_ids:
        for _ in range(points_per_trip):
            data = generate_telemetry(trip_id)
            try:
                resp = requests.post(API_URL, json=data)
                if resp.status_code == 200:
                    print(f"Sent: {data}")
                else:
                    print(f"Failed to send data: {resp.text}")
            except Exception as e:
                print(f"Error sending data: {e}")
            time.sleep(delay)  # wait before next data point

if __name__ == "__main__":
    simulate_data()