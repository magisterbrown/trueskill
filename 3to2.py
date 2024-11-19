from trueskill import Variable, Gaussian, SumFactor, TruncateFactor, TrueSkill

tmsk = [Variable() for i in range(3)]
tmsk[0].set(Gaussian(4,2))
tmsk[1].set(Gaussian(3.5,1))
tmsk[2].set(Gaussian(3,1.5))

diffs = [Variable() for i in range(len(tmsk)-1)]
sum1 = SumFactor(diffs[0], [tmsk[0], tmsk[1]], [1,-1])
sum2 = SumFactor(diffs[1], [tmsk[1], tmsk[2]], [1,-1])

trunc1 = TruncateFactor(diffs[0], TrueSkill.v_win, TrueSkill.w_win, 0.74)
trunc2 = TruncateFactor(diffs[1], TrueSkill.v_win, TrueSkill.w_win, 0.74)
sumfacs = [sum1, sum2]
truncfacs = [trunc1, trunc2]

suml = len(sumfacs)
for i in range(10):
    for x in range(1):
        print(x)


print(tmsk)
