cmake_minimum_required(VERSION 3.22)
project(crack-manager CXX)

set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

set(CMAKE_CXX_STANDARD 23)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

# For coverage report
set(CMAKE_CXX_FLAGS "--coverage -fprofile-arcs -ftest-coverage") 
set(CMAKE_EXE_LINKER_FLAGS "--coverage")

# include(build/Relese/generators/conan_toolchain.cmake)

list(APPEND CMAKE_PREFIX_PATH "${CMAKE_BINARY_DIR}")
find_package(userver REQUIRED core REQUIRED)

find_package(sole REQUIRED sole REQUIRED)
find_package(libenvpp REQUIRED libenvpp REQUIRED)

userver_setup_environment()

file(GLOB_RECURSE SOURCES CONFIGURE_DEPENDS ${CMAKE_CURRENT_SOURCE_DIR}/src/*.cpp)
add_library(${PROJECT_NAME}_objs OBJECT ${SOURCES})

message(STATUS "${CXX_STANDARD}" ${CXX_STANDARD})
message(STATUS "${HAS_CXX}" ${HAS_CXX})

target_link_libraries(${PROJECT_NAME}_objs userver::core sole::sole libenvpp::libenvpp)
target_include_directories(${PROJECT_NAME}_objs PUBLIC 
    src
    ${SOLE_INCLUDE_DIRS}
    ${CONAN_INCLUDE_DIRS}
)
# 
add_executable(${PROJECT_NAME} main.cpp)
target_link_libraries(${PROJECT_NAME} ${PROJECT_NAME}_objs)
# 
 
# 
add_executable(${PROJECT_NAME}-unittest unittests/task_scheduler_test.cpp)
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

