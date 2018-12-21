def f1(n):
    if n < 0:
        return n
    else:
        return f2(n)

def f2(n):
    return f1(n-1)

f1(5)
