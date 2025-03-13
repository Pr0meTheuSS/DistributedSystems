import pytest

async def test_get_status_invalid_id(service_client):
    """ Проверяем, что получение статуса с неправильным requestId возвращает ошибку """
    response = await service_client.get("/api/hash/status", params={"requestId": "-12345"})
    assert response.status == 400

async def test_get_status_valid_request(service_client):
    """ Проверяем получение статуса для существующего requestId """
    request_data = {
        "hash": "5d41402abc4b2a76b9719d911017c592",
        "maxLength": 5
    }
    response = await service_client.post("/api/hash/crack", json=request_data)
    assert response.status == 201
    json_resp = response.json()
    request_id = json_resp["requestId"]

    response = await service_client.get("/api/hash/status", params={"requestId": request_id})
    assert response.status == 200
    json_resp = response.json()
    assert json_resp["requestId"] == request_id
    assert json_resp["status"] in ["IN_PROGRESS", "IN_QUEUE"]
