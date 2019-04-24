CC = g++
CFLAGS = --std=c++11

all:
	$(CC) $(CFLAGS) src/teko.cpp -o teko
