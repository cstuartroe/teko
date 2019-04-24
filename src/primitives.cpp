#include <stdio.h>

bool streq(char *s1, char *s2) {
	if (*s1 == 0 && *s2 == 0) {
		return true;
	} else if (*s1 == *s2) {
		return streq(++s1, ++s2);
	} else {
		return false;
	}
}

int gooddiv(int n1, int n2) {
	int out = 0;
	while(n1 < 0) {
		out--;
		n1 = n1 + n2;
	}
	while(n1 >= n2) {
		out++;
		n1 = n1 - n2;
	}
	return out;
}

int goodmod(int n1, int n2) {
	int d = gooddiv(n1, n2);
	return n1 + (-d*n2);
}

struct FlexArray {
  
}
