import pytest

# 1. Проверяем, что ручки доступны
@pytest.mark.parametrize("method, path", [
    ("POST", "/api/hash/crack"),
    ("GET", "/api/hash/status"),
    ("PATCH", "/internal/api/manager/hash/crack/request")
])
async def test_handlers_availability(service_client, method, path):
    response = await service_client.request(method, path)
    assert response.status != 404, f"Ручка {path} не найдена"
    assert response.status != 405, f"Метод {method} не поддерживается на {path}"

# 2. Проверка контрактов ручек

async def test_crack_hash_request(service_client):
    """ Проверяем, что запрос на перебор хэша работает корректно """
    request_data = {
        "hash": "5d41402abc4b2a76b9719d911017c592",
        "maxLength": 5
    }

    response = await service_client.post("/api/hash/crack", json=request_data)
    
    assert response.status == 201
    json_resp = response.json()
    
    assert "requestId" in json_resp, "Ответ не содержит requestId"
    assert isinstance(json_resp["requestId"], str), "requestId должен быть строкой"

async def test_get_status(service_client):
    """ Проверяем, что получение статуса задачи работает корректно """
    response = await service_client.get("/api/hash/status", params={"requestId": "-12345"})
    
    assert response.status == 400

async def test_patch_status(service_client):
    ...

# 3. Полный юзкейс пользователя

async def test_full_user_flow(service_client):
    """ Пользователь отправляет запрос, получает статус, обновляет статус, снова получает статус """

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
    assert json_resp["status"] in ["IN_PROGRESS", "PENDING"]

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

