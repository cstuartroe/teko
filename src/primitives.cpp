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

// digitstrings are big-endian: the head of the list contains the least significant digit, and its next contains more significant digits
struct digitstring {
	char digit; // a digit 0-9
	bool neg;
	digitstring *next;

	digitstring() {}

	digitstring(char *s) {
		next = 0;
		digitstring *temp;
		if (*s == '-') { neg = true; s++; } else { neg = false; }
		while (*s == '0' && *(s+1) != 0) { s++; } // remove leading zeros
		while (*(s+1) != 0) {
			temp = new digitstring();
			temp->digit = ctoi(*s);
			temp->next = next;
			temp->neg = neg;
			next = temp;
			s++;
		}
		digit = ctoi(*s);

	}

	char* tostr() {
		int l = len();
		char *ptr;
		if (neg) {
			ptr = new char[l+2];
			*ptr = '-';
			copystr(ptr+1, l);
		} else {
			ptr = new char[l+1];
			copystr(ptr, l);
		}
		return ptr;
	}

	void copystr(char *s, int l) {
		if (next != 0) {
			next->copystr(s,l-1);
		}
		s[l-1] = itoc(digit);
		s[l] = 0;
	}

	int len() {
		//printf("%p\n",next);
		if (next == 0) {
			return 1;
		} else {
			return next->len() + 1;
		}
	}

	void add_char(char n) {
		char mult = neg ? -1 : 1;
		n = n*mult;

		digit = digit + n;
		char carry = gooddiv(digit,10);
		digit = goodmod(digit, 10);

		if (carry == 0) { return; }

		if (next == 0) {
			next = new digitstring();
			next->digit = 0;
			next->next = 0;
			next->neg = (carry < 0);
		}
		next->add_char(mult*carry);
		settle_carry();
	}

	// returns the absolute value if result is negative, 0 otherwise
	digitstring *sum_digitstring(digitstring other) {
		digitstring *out = new digitstring();
		if (next != 0 && other.next != 0) {
			out->next = next->sum_digitstring(*other.next);
		} else if (other.next != 0) {
			out->next = other.next->copy();
		} else if (next != 0) {
			out->next = next->copy();
		} else {
			out->next = 0;
		}
		out->digit = digit;
		out->neg = neg;
		out->settle_carry();

		char mult = other.neg ? -1 : 1;
		out->add_char(other.digit*mult);

		return out;
	}

	void settle_carry() {
		if(next == 0) { return; }
		if (next->neg != neg) {
			next->add_char(neg ? -1 : 1);
			neg = !neg;
			digit = goodmod(10-digit, 10);
			next->settle_carry();
		}
		if (next->digit == 0  && next->next == 0) {
			delete next;
			next = 0;
		}
	}

	digitstring *copy() {
		digitstring *out = new digitstring();
		out->digit = digit;
		out->neg = neg;
		out->next = (next == 0) ? 0 : out->next = next->copy();
		return out;
	}
};

void assert_sum(char *nums1, char *nums2, char *answer) {
	digitstring d1 = digitstring(nums1);
	digitstring d2 = digitstring(nums2);

	digitstring result;

	result = *d1.sum_digitstring(d2);
	if (!streq(result.tostr(), answer)) {
		printf("%s + %s gave %s instead of %s\n", nums1, nums2, result.tostr(), answer);
	}

	result = *d2.sum_digitstring(d1);
	if (!streq(result.tostr(), answer)) {
		printf("%s + %s gave %s instead of %s\n", nums2, nums1, result.tostr(), answer);
	}
}

int main(int argc, char *argv[]) {
	char p5556[] = "5556", p12000[] = "12000", n5556[] = "-5556", p17556[] = "17556", p6444[] = "6444",
	     p483[]  = "483",  p6039[]  = "6039",  n5073[] = "-5073";

	assert_sum(p5556,p12000,p17556);
	assert_sum(n5556,p12000,p6444);
	assert_sum(p483, p5556, p6039);
	assert_sum(p483, n5556, n5073);


	char in1[16], in2[16];
	scanf("%15s",in1);
	scanf("%15s",in2);
	digitstring d1 = digitstring(in1), d2 = digitstring(in2);
	digitstring d3 = *d1.sum_digitstring(d2);
	printf("%s\n", d3.tostr());
	return 0;
}