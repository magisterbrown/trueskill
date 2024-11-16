package main

import (
    "fmt"
    "math"
    "math/rand"
)

type Gaussian struct {
    pi, tau float64
}
func (g Gaussian) mu() float64 {
    return g.tau/g.pi
}
func (g Gaussian) sigma() float64 {
    return 1/math.Sqrt(g.pi)
}
func (g Gaussian) String() string {
    return fmt.Sprintf("Gaussian(mu: %f, sigma: %f)", g.mu(), g.sigma())
}
func NewGaussian(mu float64, sigma float64) *Gaussian{
    pi := 1/(sigma*sigma)
    return &Gaussian{
        pi: pi,
        tau: mu*pi,
    }
}
func AddGaussian(g1 Gaussian, g2 Gaussian) *Gaussian{
    return NewGaussian(g1.mu()+g2.mu(), math.Sqrt(g1.sigma()*g1.sigma()+g2.sigma()*g2.sigma()))
}
func SubGaussian(g1 *Gaussian, g2 *Gaussian) *Gaussian{
    return NewGaussian(g1.mu()-g2.mu(), math.Sqrt(g1.sigma()*g1.sigma()+g2.sigma()*g2.sigma()))
}
func MultGaussian(g1 *Gaussian, g2 *Gaussian) *Gaussian{
    return &Gaussian{
        pi: g1.pi+g2.pi,
        tau: g1.tau+g2.tau,
    }
}
func DivGaussian(g1 *Gaussian, g2 *Gaussian) *Gaussian{
    return &Gaussian{
        pi: g1.pi-g2.pi,
        tau: g1.tau-g2.tau,
    }
}

func propogateWin(g *Gaussian, draw_margin float64) *Gaussian{
    return propogateExpectation(g, draw_margin, math.Inf(1))
}
func propogateDraw(g *Gaussian, draw_margin float64) *Gaussian{
    return propogateExpectation(g, -draw_margin, draw_margin)
}
func propogateExpectation(g *Gaussian, low float64, high float64) *Gaussian{
    avg := 0.
    avg_squared := 0.
    idx := 0.
    for range 2000000 {
        vl := rand.NormFloat64()*g.sigma()+g.mu()
        if(vl>low && vl<high) {
            avg = (vl+idx*avg)/(idx+1)
            avg_squared = (vl*vl+idx*avg_squared)/(idx+1)
            idx++
        }
    }
    return NewGaussian(avg, math.Sqrt(avg_squared-avg*avg))
}

func main() {
    outcome_exp:= NewGaussian(-5,6.9)
    fmt.Println(outcome_exp)
    propogated_exp:= propogateDraw(outcome_exp, 0.74)
    fmt.Println(propogated_exp)

    team_skills := []*Gaussian {NewGaussian(2,2), NewGaussian(1,1), NewGaussian(3,1.5)}
    is_draw := []bool {true, false}

    draw_margin := 0.74
    prior_outcomes := []*Gaussian {}
    posterior_outcomes :=  []*Gaussian {}

    // Right team update
    for i:=1; i<len(team_skills); i++ {
        prior := SubGaussian(team_skills[i-1], team_skills[i])
        var posterior *Gaussian
        if is_draw[i-1] {
            posterior = propogateWin(prior, draw_margin)
        } else {
            posterior = propogateDraw(prior, draw_margin)
        }

        prior_outcomes = append(prior_outcomes, prior)
        posterior_outcomes = append(posterior_outcomes, posterior)
        sampler := DivGaussian(posterior, prior)
        team_skills[i] = MultGaussian(team_skills[i], SubGaussian(team_skills[i-1], sampler))
    }
    fmt.Println(draw_margin)
    fmt.Println(prior_outcomes)
    fmt.Println(posterior_outcomes)
    fmt.Println(team_skills)
   //fmt.Println(NewGaussian(avg, 
}
