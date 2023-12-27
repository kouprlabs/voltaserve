cmake_minimum_required(VERSION 3.10)
project(Voltaserve C)

set(CMAKE_C_STANDARD 11)

add_library(voltaserve SHARED lib/voltaserve.c)

set_target_properties(voltaserve PROPERTIES LIBRARY_OUTPUT_DIRECTORY "${CMAKE_BINARY_DIR}")