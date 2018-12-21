public class Ordering {
    static int f1(int n) {
        if (n < 0) {
            return n;
        } else {
            return f2(n);
        }
    }
        
    static int f2(int n) {
        return f1(n-1);
    }
      
    public static void main(String[] args) {
        f1(5);
    }
}
