import numpy as np
from trueskill import Gaussian
from scipy.stats import norm
SZ = 10000000
prior = Gaussian(1.5, 9)
post = Gaussian(3.1, 5.2)
print(post*prior)
v1 = np.random.normal(loc=prior.mu, scale=prior.sigma, size=SZ)
probs = norm.pdf(v1, loc=post.mu, scale=post.sigma)
gen = np.random.random_sample(size=probs.shape)
print(gen.mean())
print(gen.std())
ava = v1[probs>gen]
print(ava.shape)
print(ava.mean())
print(ava.std())
#v2 = np.random.normal(loc=4.612, scale=3.241, size=SZ)
