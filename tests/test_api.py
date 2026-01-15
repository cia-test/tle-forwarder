import pytest
import requests
import time


BASE_URL = "http://localhost:8000"


@pytest.fixture(scope="session", autouse=True)
def wait_for_service():
    """Wait for the service to be ready before running tests."""
    max_retries = 30
    for i in range(max_retries):
        try:
            response = requests.get(f"{BASE_URL}/health", timeout=2)
            if response.status_code == 200:
                return
        except requests.exceptions.RequestException:
            pass
        time.sleep(1)
    pytest.fail("Service did not become ready in time")


def test_health_endpoint():
    """Test that the health endpoint returns healthy status."""
    response = requests.get(f"{BASE_URL}/health")
    assert response.status_code == 200
    data = response.json()
    assert data["status"] == "healthy"


def test_root_endpoint():
    """Test that the root endpoint returns service information."""
    response = requests.get(f"{BASE_URL}/")
    assert response.status_code == 200
    data = response.json()
    assert data["service"] == "TLE Forwarder"
    assert "endpoints" in data
    assert "/tle" in data["endpoints"]


def test_tle_endpoint_no_params():
    """Test that TLE endpoint requires parameters."""
    response = requests.get(f"{BASE_URL}/tle")
    assert response.status_code == 400


def test_tle_endpoint_with_satellite_id():
    """Test fetching TLE data by satellite ID (ISS)."""
    response = requests.get(f"{BASE_URL}/tle", params={"satellite_id": "25544"})
    assert response.status_code == 200
    assert len(response.text) > 0
    # Basic TLE format validation
    lines = response.text.strip().split("\n")
    assert len(lines) >= 2


def test_tle_endpoint_with_name():
    """Test fetching TLE data by satellite name."""
    response = requests.get(f"{BASE_URL}/tle", params={"name": "ISS"})
    assert response.status_code == 200
    assert len(response.text) > 0


def test_tle_endpoint_with_group():
    """Test fetching TLE data by satellite group."""
    response = requests.get(f"{BASE_URL}/tle", params={"group": "stations"})
    assert response.status_code == 200
    assert len(response.text) > 0


def test_tle_endpoint_invalid_satellite():
    """Test fetching TLE data with invalid satellite ID."""
    response = requests.get(f"{BASE_URL}/tle", params={"satellite_id": "99999999"})
    # Should return 404 or error status
    assert response.status_code in [404, 500, 503]
