import pytest

# Проверяем, что ручки доступны
@pytest.mark.parametrize("method, path", [
    ("POST", "/api/hash/crack"),
    ("GET", "/api/hash/status"),
    ("PATCH", "/internal/api/manager/hash/crack/request"),
    ("GET", "/api/hash/crack/answer")
])
async def test_handlers_availability(service_client, method, path):
    response = await service_client.request(method, path)
    assert response.status != 405, f"Метод {method} не поддерживается на {path}"
