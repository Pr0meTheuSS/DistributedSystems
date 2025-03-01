import pytest

async def test_full_user_flow(service_client):
    """ Пользователь отправляет запрос, получает статус """
    
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

    update_data = {
        "requestId": request_id,
        "status": "READY"
    }
    response = await service_client.patch("/internal/api/manager/hash/crack/request", json=update_data)
    assert response.status == 204

    response = await service_client.get("/api/hash/status", params={"requestId": request_id})
    assert response.status == 200
    json_resp = response.json()
    assert json_resp["requestId"] == request_id
    assert json_resp["status"] == "READY"

    
    response = await service_client.get("/api/hash/crack/answer", params={"requestId": request_id})
    assert response.status == 404

    request_data = {
        "requestId": request_id,
        "partNumber": 0,
        "words": ["ab"]
    }

    response = await service_client.post("/internal/api/manager/hash/crack/request", json=request_data)
    assert response.status == 204

    response = await service_client.get("/api/hash/crack/answer", params={"requestId": request_id})
    assert response.status == 200
    json_resp = response.json()
    assert "ab" in json_resp["words"]
