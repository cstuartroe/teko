CC = g++
CFLAGS = --std=c++11 -g

all:
	$(CC) $(CFLAGS) src/teko.cpp -o teko
