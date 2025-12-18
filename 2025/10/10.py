import sys
import os
import time

import z3 # install z3-solver


class Machine:
    def __init__(self, machineStr: str) -> None:
        fields = machineStr.split(" ")
        if len(fields) < 3:
            raise Exception("invalid machine str")

        btnStrs = [btn[1:-1] for btn in fields[1:-1]]
        joltStr = fields[-1][1:-1]

        self.joltage = list(map(int, joltStr.split(",")))
        self.buttons = [list(map(int, btn.split(","))) for btn in btnStrs]

        m = len(self.joltage)
        n = len(self.buttons)
        mat = [[0] * n for _ in range(m)]
        for i, btn in enumerate(self.buttons):
            for pos in btn:
                mat[pos][i] = 1
        self.btnMatrix = mat

    def solve(self):
        n = len(self.buttons)
        coef = self.btnMatrix

        s = z3.Optimize()
        vars = []

        for i in range(n):
            var = z3.Int(f"b{i}")
            vars.append(var)
            s.add(var >= 0)

        for j, joltage in enumerate(self.joltage):
            s.add(sum(coef[j][b] * vars[b] for b in range(n)) == joltage)

        presses = z3.Sum(*vars)
        s.minimize(presses)

        if s.check() != z3.sat:
            return None

        m = s.model().eval(presses)
        return int(f"{m}")


def sove_machines_in_file(inpFilePath: str):
    total = 0
    solved = 0
    totalPresses = 0
    with open(inpFilePath) as file:
        for line in file:
            total += 1
            m = Machine(line.rstrip())
            res = m.solve()
            if res is not None:
                solved += 1
                totalPresses += res
    print(f"solved: {solved}/{total}")
    print(f"min presses: {totalPresses}")


def main():
    if len(sys.argv) != 2:
        print("usage: python 10.py <input-file-path>")
        sys.exit(1)
    fp = sys.argv[1]
    if not os.path.isfile(fp):
        print("error: input file not found")
        sys.exit(2)
    start = time.perf_counter()
    sove_machines_in_file(fp)
    print(f"time {time.perf_counter() - start} seconds")


if __name__ == "__main__":
    main()
