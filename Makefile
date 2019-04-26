CC = g++
CFLAGS = --std=c++11 -g

all:
	$(CC) $(CFLAGS) src/teko.cpp -o teko

pretty:
	uncrustify -c .uncrustify --no-backup src/*
	$(CC) $(CFLAGS) src/teko.cpp -o teko
