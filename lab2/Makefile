
### `Makefile`

# Makefile для установки зависимостей, сборки и запуска проекта

# Переменные
USERVER_REPO = https://github.com/userver-framework/userver
USERVER_TAG = v2.7

# Установить Conan
install-conan:
	@echo "Устанавливаем Conan..."
	pip install conan

# Клонировать репозиторий Userver и переключиться на тег v2.7
clone-userver:
	@echo "Клонируем репозиторий Userver..."
	git clone $(USERVER_REPO)
	cd userver && git checkout $(USERVER_TAG)

# Собрать Conan пакет для Userver
build-userver-package:
	@echo "Собираем Conan пакет для Userver..."
	cd userver && conan create . --build=missing

# Сборка менеджера
build-manager:
	@echo "Собираем менеджер..."
	mkdir -p manager/build
	cd manager/build
	conan install .. --build=missing
	cmake .. -DCMAKE_TOOLCHAIN_FILE=Release/generators/conan_toolchain.cmake -DCMAKE_POLICY_DEFAULT_CMP0091=NEW -DCMAKE_BUILD_TYPE=Release
	make -j$(nproc)

# Сборка воркера
build-worker:
	@echo "Собираем воркера..."
	mkdir -p worker/build
	cd worker/build
	conan install .. --build=missing
	cmake .. -DCMAKE_TOOLCHAIN_FILE=Release/generators/conan_toolchain.cmake -DCMAKE_POLICY_DEFAULT_CMP0091=NEW -DCMAKE_BUILD_TYPE=Release
	make -j$(nproc)

# Запуск менеджера
run-manager:
	@echo "Запускаем менеджер..."
	cd manager
	export $(grep -v '^#' .env | xargs)
	./build/crack-manager -c ../static_config.yaml

# Запуск воркера
run-worker:
	@echo "Запускаем воркера..."
	./worker/build/crack-worker -c ../static_config.yaml
