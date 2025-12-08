package com.github.litmuschaos.testtools;

public class CpuLoadGenerator {
	private long fib0 = 0;
	private long fib1 = 1;
	private static final long FIB_THRESHOLD = Long.MAX_VALUE / 2;
	

	public void generateLoad() {
		long fib2 = fib0 + fib1;
		if (fib2 > FIB_THRESHOLD) {
			fib0 = 0;
			fib1 = 1;
		} else {
			fib0 = fib1;
			fib1 = fib2;
		}
	}
}

