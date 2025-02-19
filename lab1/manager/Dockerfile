# Используем Alpine как базовый образ для сборки
FROM debian:bookworm-slim AS builder

# Устанавливаем необходимые зависимости с доступной версией g++
RUN apt update && apt install -y g++ make cmake python3 git build-essential python3.11-venv

# Создаём виртуальное окружение для Python
RUN python3 -m venv /venv

# Создаём рабочую директорию
WORKDIR /app

# Активируем виртуальное окружение и устанавливаем Conan
RUN /venv/bin/pip install conan

# Настроим git для увеличения буфера
RUN git config --global http.postBuffer 524288000

# Клонируем репозиторий с GitHub
RUN git clone https://github.com/userver-framework/userver.git

# Переходим на нужную ветку
RUN cd userver && git checkout v2.7

# Создаём профиль, если его нет
RUN /venv/bin/conan profile detect

# Обновляем профиль для использования GCC 12
# RUN /venv/bin/conan profile update-settings compiler.version=12 default
# RUN /venv/bin/conan profile update-settings compiler.libcxx=libstdc++ default

# Создаём и устанавливаем пакет через Conan
RUN /venv/bin/conan create /app/userver/. --build=missing

# Копируем исходники
COPY . .

# Собираем приложение
RUN mkdir build && cd build && \
    /venv/bin/conan install .. --build=missing && \
    cmake .. && \
    make -j$(nproc)

# Второй этап - минимальный образ
FROM alpine:latest

# Устанавливаем необходимые библиотеки
RUN apk add --no-cache libstdc++

# Копируем скомпилированный бинарник
COPY --from=builder /app/build/my_webserver /usr/local/bin/my_webserver

# Указываем рабочую директорию
WORKDIR /usr/local/bin

# Открываем порт 8080
EXPOSE 8080

# Запускаем приложение
CMD ["./crack-manager -c ../static_config.yaml"]
