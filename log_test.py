class CPUState:
    def __init__(self, a, x, y, p, sp) -> None:
        self.a = a
        self.x = x
        self.y = y
        self.p = p
        self.sp = sp

    def __str__(self) -> str:
        return f"A: {self.a:2x} X: {self.x:2x} Y: {self.y:2x} P: {self.p:2x} SP: {self.sp:2x}"

    def __eq__(self, value: object) -> bool:
        return self.a == value.a and self.x == value.x and self.y == value.y and self.p == value.p and self.sp == value.sp

nestest_states = []

with open("nestest.log") as nestest_log_file:
    for line in nestest_log_file:
        a = int(line[50:52], 16)
        x = int(line[55:57], 16)
        y = int(line[60:62], 16)
        p = int(line[65:67], 16)
        sp = int(line[71:73], 16)
        nestest_states.append(CPUState(a, x, y, p, sp))

cpu_states = []
with open("logs.txt") as log_file:
    for line in log_file:
        a = int(line[25:27].strip(), 16)
        x = int(line[30:32].strip(), 16)
        y = int(line[35:37].strip(), 16)
        p = int(line[40:42].strip(), 16)
        sp = int(line[46:48].strip(), 16)
        cpu_states.append(CPUState(a, x, y, p, sp))


for i in range(min(len(nestest_states), len(cpu_states))):
    try:
        assert nestest_states[i] == cpu_states[i]
    except:
        print(f"Differences in line {i + 1}:")
        print(f"nestest: {nestest_states[i]}")
        print(f"logs:    {cpu_states[i]}")
        break