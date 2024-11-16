from trueskill import Gaussian
import matplotlib.pyplot as plt
from matplotlib.widgets import Slider
from scipy.stats import norm
import numpy as np

xs = np.arange(-30, 30, 0.1)
def pdf_norm(gaus):
    return norm.pdf(xs, gaus.mu, gaus.sigma)
old_exp = Gaussian(mu=-5.000, sigma=6.909)
ys = pdf_norm(old_exp)
global new_exp
new_sigma = 6.909
new_mu = -5
new_exp = Gaussian(mu=new_mu, sigma=new_sigma)
fig, (ax1, mean, std) = plt.subplots(3,1,  height_ratios= [8, 1, 1])
ax1.plot(xs, ys, label="old exp")
ax1.set_ylim(0, 1.3)
new_plot, = ax1.plot(xs, pdf_norm(new_exp), label="new exp")
div_res, = ax1.plot(xs, pdf_norm(new_exp)/pdf_norm(old_exp), label="direct")
real_res, = ax1.plot(xs, pdf_norm(new_exp/old_exp), label="real")

#new_exp = Gaussian(mu=val, sigma=new_exp.sigma)
#new_exp = Gaussian(mu=new_exp.mu, sigma=val)
means = Slider(mean, 'mean', -10, 10)
stds = Slider(std, 'std', 0.5, 7)

def update(val):
    new_exp = Gaussian(mu=means.val, sigma=stds.val)
    new_plot.set_ydata(pdf_norm(new_exp))
    div_res.set_ydata(pdf_norm(new_exp)/pdf_norm(old_exp))
    real_res.set_ydata(pdf_norm(new_exp/old_exp))
means.on_changed(update)
stds.on_changed(update)
ax1.legend()
#ln, = ax.plot(range(5))
#plt.ion()

plt.show()

