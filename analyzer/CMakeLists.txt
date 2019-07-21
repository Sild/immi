cmake_minimum_required(VERSION 3.5)

project(stock_analyzer VERSION 0.0.0)
include(FetchContent)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_REQUIRED_FLAGS -std=c++17)
set(CMAKE_COMPILE_OPTIONS "-Wpedantic;-Wall;-Wextra")
set(CMAKE_CXX_COMPILER "clang++")

set(CSV_REPO https://github.com/ben-strasser/fast-cpp-csv-parser.git)
set(CSV_VERSION master)

FetchContent_Declare(
  csv
  GIT_REPOSITORY ${CSV_REPO}
  GIT_TAG        ${CSV_VERSION}
)
FetchContent_GetProperties(csv)
if(NOT csv_POPULATED)
  FetchContent_Populate(csv)
endif()

INCLUDE_DIRECTORIES(
    src
    ${csv_SOURCE_DIR}
)


add_executable(stock_analyzer
    src/main.cpp
)


install(TARGETS stock_analyzer RUNTIME DESTINATION bin)

set(CPACK_GENERATOR DEB)

set(CPACK_PACKAGE_VERSION_MAJOR "0")
set(CPACK_PACKAGE_VERSION_MINOR "0")
set(CPACK_PACKAGE_VERSION_PATCH "0")

set(CPACK_PACKAGE_CONTACT sildtm@icloud.com)

set(CPACK_DEBIAN_PACKAGE_SHLIBDEPS ON)

include(CPack)
