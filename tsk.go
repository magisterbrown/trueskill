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
    if g.pi == 0 {
        return 0
    }
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
func AddGaussian(g1 *Gaussian, g2 *Gaussian) *Gaussian{
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
    for range 8000000 {
        vl := rand.NormFloat64()*g.sigma()+g.mu()
        if(vl>low && vl<high) {
            avg = (vl+idx*avg)/(idx+1)
            avg_squared = (vl*vl+idx*avg_squared)/(idx+1)
            idx++
        }
    }
    return NewGaussian(avg, math.Sqrt(avg_squared-avg*avg))
}
func p(g *Gaussian) string{
    return g.String()
}

func main() {
    outcome_exp:= NewGaussian(-5,6.9)
    fmt.Println(outcome_exp)
    propogated_exp:= propogateDraw(outcome_exp, 0.74)
    fmt.Println(propogated_exp)

    //team_skills := []*Gaussian {NewGaussian(4,2), NewGaussian(3.5,1), NewGaussian(3,1.5)}
    team_skills := []*Gaussian {NewGaussian(1,5.135), NewGaussian(3,5.135), NewGaussian(6,5.135), NewGaussian(14,5.135)}
    is_draw := []bool {false, false, false, false}
    samplers_winner := make([]*Gaussian, len(team_skills))
    samplers_looser := make([]*Gaussian, len(team_skills))
    for i := range team_skills{
        samplers_winner[i] = &Gaussian{pi:0, tau:0}
        samplers_looser[i] = &Gaussian{pi:0, tau:0}
    }
    prts:=func() {
        for i := range team_skills {
            fmt.Printf("%s ", MultGaussian(MultGaussian(samplers_winner[i], samplers_looser[i]), team_skills[i]).String())
        }
        fmt.Printf("\n")
    }
    prts()

    draw_margin := 0.74
    
    for j:=0; j<10; j++ {
        fmt.Printf("Step %d\n", j)
        max_delta := 0.
        // Right team update
        for i:=1; i<len(team_skills)-1; i++ {
            winner_skill := MultGaussian(team_skills[i-1], samplers_winner[i-1])
            looser_skill := MultGaussian(team_skills[i], samplers_looser[i])
            prior := SubGaussian(winner_skill, looser_skill)
            var posterior *Gaussian
            if is_draw[i-1] {
                posterior = propogateDraw(prior, draw_margin)
            } else {
                posterior = propogateWin(prior, draw_margin)
            }

            sampler := SubGaussian(winner_skill, DivGaussian(posterior, prior))
            max_delta = max(max_delta, math.Sqrt(math.Abs(prior.pi-posterior.pi)), math.Abs(prior.tau-posterior.tau))
            samplers_winner[i] = sampler
        }

        // Left team update
        for i:=len(team_skills)-1; i>1; i-- {
            winner_skill := MultGaussian(team_skills[i-1], samplers_winner[i-1])
            looser_skill := MultGaussian(team_skills[i], samplers_looser[i])
            prior := SubGaussian(winner_skill, looser_skill)
            var posterior *Gaussian
            if is_draw[i-1] {
                posterior = propogateDraw(prior, draw_margin)
            } else {
                posterior = propogateWin(prior, draw_margin)
            }

            sampler := AddGaussian(looser_skill, DivGaussian(posterior, prior))
            max_delta = max(max_delta, math.Sqrt(math.Abs(prior.pi-posterior.pi)), math.Abs(prior.tau-posterior.tau))
            samplers_looser[i-1] = sampler
        }
        //prts()
        //fmt.Printf("Max delta left: %f\n", max_delta)
    }
    var posterior *Gaussian
    winner_skill := MultGaussian(team_skills[0], samplers_winner[0])
    looser_skill := MultGaussian(team_skills[1], samplers_looser[1])
    prior := SubGaussian(winner_skill, looser_skill)
    if is_draw[0] {
        posterior = propogateDraw(prior, draw_margin)
    } else {
        posterior = propogateWin(prior, draw_margin)
    }
    sampler := AddGaussian(looser_skill, DivGaussian(posterior, prior))
    samplers_looser[0] = sampler


    last:=len(team_skills)-1
    winner_skill = MultGaussian(team_skills[last-1], samplers_winner[last-1])
    looser_skill = MultGaussian(team_skills[last], samplers_looser[last])
    prior = SubGaussian(winner_skill, looser_skill)
    if is_draw[last-1] {
        posterior = propogateDraw(prior, draw_margin)
    } else {
        posterior = propogateWin(prior, draw_margin)
    }

    sampler = SubGaussian(winner_skill, DivGaussian(posterior, prior))
    samplers_winner[last] = sampler

    fmt.Println(team_skills)
    prts()
   //fmt.Println(NewGaussian(avg, 
}
