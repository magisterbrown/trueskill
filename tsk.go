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
func p(g *Gaussian) String{
    return g.String()
}

func main() {
    outcome_exp:= NewGaussian(-5,6.9)
    fmt.Println(outcome_exp)
    propogated_exp:= propogateDraw(outcome_exp, 0.74)
    fmt.Println(propogated_exp)

    //team_skills := []*Gaussian {NewGaussian(4,2), NewGaussian(3.5,1), NewGaussian(3,1.5)}
    team_skills := []*Gaussian {NewGaussian(1,5.135), NewGaussian(3,5.135), NewGaussian(6,5.135), NewGaussian(14,5.135)}
    samplers_winner := make([]*Gaussian, len(team_skills))
    samplers_looser := make([]*Gaussian, len(team_skills))
    fmt.Println(samplers_winner)
    fmt.Println(samplers_looser)

    //skills_buffer := make([]*Gaussian, len(team_skills))
    //skills_buffer_r := make([]*Gaussian, len(team_skills))
    //copy(skills_buffer, team_skills)
    //copy(skills_buffer_r, team_skills)

    //diff_priors := make([]*Gaussian, len(team_skills)-1)
    //for i := range diff_priors {
    //    diff_priors[i] = &Gaussian{pi:0, tau:0}
    //}
    is_draw := []bool {false, false, false, false}

    draw_margin := 0.74
    
    for j:=0; j<10; j++ {
        fmt.Printf("Step %d\n", j)
        fmt.Println(team_skills)
        max_delta := 0.
        // Right team update
        //fmt.Println("Right")
        for i:=1; i<len(team_skills)-1; i++ {
            subr := SubGaussian(team_skills[i-1], team_skills[i])
            prior := subr //MultGaussian(diff_priors[i-1], subr)
            //diff_priors[i-1] = prior
            var posterior *Gaussian
            if is_draw[i-1] {
                posterior = propogateDraw(prior, draw_margin)
            } else {
                posterior = propogateWin(prior, draw_margin)
            }
            //diff_priors[i-1] = posterior

            sampler := DivGaussian(posterior, prior)
            new_skill := MultGaussian(team_skills[i], SubGaussian(team_skills[i-1], sampler))
            max_delta = max(max_delta, math.Sqrt(math.Abs(prior.pi-posterior.pi)), math.Abs(prior.tau-posterior.tau))
            //fmt.Printf("Replaced skill %d from %s -> %s\n", i, team_skills[i], new_skill)
            team_skills[i] = new_skill
        }
        fmt.Printf("Max delta right: %f\n", max_delta)

        //fmt.Println("Left")
        // Left team update
        for i:=len(team_skills)-1; i>1; i-- {
            subr := SubGaussian(team_skills[i-1], team_skills_r[i])
            prior := subr //MultGaussian(diff_priors[i-1], subr)
            //prior := SubGaussian(team_skills[i-1], team_skills[i])
            var posterior *Gaussian
            if is_draw[i-1] {
                posterior = propogateDraw(prior, draw_margin)
            } else {
                posterior = propogateWin(prior, draw_margin)
            }
            //diff_priors[i-1] = posterior

            sampler := DivGaussian(posterior, prior)
            new_skill_r := MultGaussian(team_skills_r[i-1], AddGaussian(team_skills_r[i], sampler))
            team_skills_r[i-1] = new_skill_r
            //fmt.Printf("Replaced skill %d from %s -> %s\n", i-1, team_skills[i-1], new_skill)
            new_skill := MultGaussian(team_skills[i-1], AddGaussian(team_skills_r[i], sampler))
            max_delta = max(max_delta, math.Sqrt(math.Abs(prior.pi-posterior.pi)), math.Abs(prior.tau-posterior.tau))
            team_skills[i-1] = new_skill
        }
        fmt.Printf("Max delta left: %f\n", max_delta)
    }
    fmt.Println(draw_margin)
    fmt.Println(team_skills)
   //fmt.Println(NewGaussian(avg, 
}
