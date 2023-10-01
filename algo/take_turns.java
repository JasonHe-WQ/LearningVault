import java.util.Scanner;
import java.util.concurrent.locks.Condition;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;

public class Main {
    private final int n;
    private final int k;
    private final Lock lock = new ReentrantLock();
    private final Condition[] conditions;
    private int currentThread = 0;
    private int counter = 1;

    public Main(int n, int k) {
        this.n = n;
        this.k = k;
        this.conditions = new Condition[n];
        for (int i = 0; i < n; i++) {
            conditions[i] = lock.newCondition();
        }
    }

    public void print(int threadNumber) {
        lock.lock();
        try {
            while (counter <= k) {
                while (threadNumber != currentThread) {
                    conditions[threadNumber].await();
                }
                if (counter <= k) {
                    System.out.print((threadNumber + 1) + ":" + counter++ + " ");
                }
                currentThread = (currentThread + 1) % n;
                conditions[currentThread].signal();
            }
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        } finally {
            lock.unlock();
        }
    }

    public static void main(String[] args) {
        Scanner scanner = new Scanner(System.in);
        int n = scanner.nextInt();
        int k = scanner.nextInt();
        Main printer = new Main(n, k);

        for (int i = 0; i < n; i++) {
            final int threadNumber = i;
            new Thread(() -> printer.print(threadNumber)).start();
        }
    }
}
