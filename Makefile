CC = g++
CFLAGS = --std=c++11 -g

all:
	uncrustify -c .uncrustify --no-backup src/*
	$(CC) $(CFLAGS) src/teko.cpp -o teko
