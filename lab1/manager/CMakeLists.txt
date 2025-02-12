cmake_minimum_required(VERSION 3.22)
project(userver-samples-hello_service CXX)

set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

include(${CMAKE_BINARY_DIR}/conan_toolchain.cmake)

find_package(userver REQUIRED core REQUIRED)
userver_setup_environment()

add_library(${PROJECT_NAME}_objs OBJECT
    src/dto/crack_hash_request.hpp
    src/dto/crack_hash_request.cpp

    src/hash_crack_handler.cpp
    src/hash_crack_handler.hpp
)

target_link_libraries(${PROJECT_NAME}_objs userver::core)
target_include_directories(${PROJECT_NAME}_objs PUBLIC 
    src 
    ${CONAN_INCLUDE_DIRS}  # Подключаем инклуды из Conan
)
# 
add_executable(${PROJECT_NAME} main.cpp)
target_link_libraries(${PROJECT_NAME} ${PROJECT_NAME}_objs)
# 
 
# 
add_executable(${PROJECT_NAME}-unittest unittests/say_hello_test.cpp)
target_link_libraries(${PROJECT_NAME}-unittest ${PROJECT_NAME}_objs userver::utest)
add_google_tests(${PROJECT_NAME}-unittest)
# 
 
# 
add_executable(${PROJECT_NAME}_benchmark benchmarks/say_hello_bench.cpp)
target_link_libraries(${PROJECT_NAME}_benchmark ${PROJECT_NAME}_objs userver::ubench)
add_google_benchmark_tests(${PROJECT_NAME}_benchmark)
# 
 
# 
# Makes a virtualenv suitable for running pytests from `test` directory and integrates with `ctest`.
userver_testsuite_add_simple()
