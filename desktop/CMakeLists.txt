cmake_minimum_required(VERSION 3.21)
project(Voltaserve C)

set(CMAKE_C_STANDARD 23)

add_executable(voltaserve cli/voltaserve.c)

set_target_properties(voltaserve PROPERTIES LIBRARY_OUTPUT_DIRECTORY "${CMAKE_BINARY_DIR}")