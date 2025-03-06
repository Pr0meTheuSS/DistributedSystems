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
