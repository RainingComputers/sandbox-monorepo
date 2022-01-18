import math
from numba import njit
import time


@njit
def divisor_sum(n: int) -> int:
    if n <= 0:
        return 0

    i = 1
    limit = math.sqrt(n)
    the_sum = 0
    while i <= limit:
        if n % i == 0:
            the_sum += i
            if i != limit:
                the_sum += n // i
        i += 1

    return the_sum


@njit
def witness_value(n: int) -> float:
    denominator = n * math.log(math.log(n))
    div_sum = divisor_sum(n)
    return div_sum/denominator


@njit
def best_witness(max_range: int, search_start: int) -> int:
    max_val, max_index = 0.0, search_start

    for i in range(search_start, max_range):
        current_value = witness_value(i)
        if current_value > max_val:
            max_index = i
            max_val = current_value

    return max_index


def main() -> None:
    print("The best witness is:", best_witness(10000000, 5041))


if __name__ == '__main__':
    start = time.time()
    main()
    end = time.time()
    print("Took", end-start)
