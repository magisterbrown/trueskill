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
func distance(g1 *Gaussian, g2 *Gaussian) float64{
    return max(math.Sqrt(math.Abs(g1.pi-g2.pi)), math.Abs(g1.tau-g2.tau))
}

func propogateWin(g *Gaussian, draw_margin float64) *Gaussian{
    return propogateExpectationA(g, draw_margin, math.Inf(1))
}
func propogateDraw(g *Gaussian, draw_margin float64) *Gaussian{
    return propogateExpectationA(g, -draw_margin, draw_margin)
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

func pdf(x float64) float64 {
    return math.Exp(-(x*x)/2)/math.Sqrt(2*math.Pi);
}
func cdf(x float64) float64 {
    if x==math.Inf(1) {
        return 1.;
    } else if x==math.Inf(-1) {
        return 0.;
    }
    return (1+math.Erf(x/math.Sqrt(2)))/2;
}
func ppf(x float64) float64 {
    return math.Erfinv(2*x-1)*math.Sqrt(2);
}

func propogateExpectationA(g *Gaussian, low float64, high float64) *Gaussian{
    alpha := (low-g.mu())/g.sigma();
    beta := (high-g.mu())/g.sigma();
    
    Z := cdf(beta) - cdf(alpha)
    
    new_mu := g.mu()+g.sigma()*(pdf(alpha)-pdf(beta))/Z
    bmul := beta*pdf(beta);
    if beta == math.Inf(1) {
        bmul = 0;
    }
    amul := alpha*pdf(alpha);
    if alpha == math.Inf(-1) {
        amul = 0;
    }

    new_sigma := math.Sqrt(1 - (bmul-amul)/Z - math.Pow((pdf(alpha)-pdf(beta))/Z, 2))*g.sigma();

    return NewGaussian(new_mu,new_sigma)
}



var MU = 25.
var SIGMA = MU / 3
var BETA = SIGMA / 2
var TAU = SIGMA / 100
var DRAW_PROBABILITY = .10

func main() {
    input_skills := []*Gaussian {NewGaussian(1,3), NewGaussian(3,3), NewGaussian(4,3) , NewGaussian(6,3), NewGaussian(14,3)}
	team_ids := []int {0,1,1,2,3}
	team_places := []int {0,1,2,3}
    output_skills := true_skill(input_skills, team_ids, team_places)
    fmt.Println(input_skills);
    fmt.Println(output_skills);
}
func true_skill(input_skills []*Gaussian, team_ids []int, team_places []int) []*Gaussian {



    prior_skills := make([]*Gaussian, len(input_skills))

    for i := range input_skills {
        prior_skills[i] = NewGaussian(input_skills[i].mu(), math.Sqrt(input_skills[i].sigma()*input_skills[i].sigma() + TAU*TAU))
    }

    likelihood_skills := make([]*Gaussian, len(input_skills))
    for i := range input_skills {
        likelihood_skills[i] = NewGaussian(prior_skills[i].mu(), math.Sqrt(prior_skills[i].sigma()*input_skills[i].sigma() + BETA*BETA))
    }

    player_skills := likelihood_skills
    
	if(len(team_ids) != len(player_skills)){
		panic("Some players are not assigned a team")
	}
	max_team_id := -1
	for i := range team_ids {
		if(team_ids[i] < 0){
			panic(fmt.Sprintf("Team %d can not have negative id", i))
		}
		if(team_ids[i] > len(team_places)-1){
			panic(fmt.Sprintf("Team with id %d does not have a place", team_ids[i]))
		}
		if(team_ids[i] > max_team_id) {
			max_team_id = team_ids[i]
		}
	}
	team_skills := make([]*Gaussian, max_team_id+1)
    team_order := make([]int, len(team_skills))
	player_pos := make([]int, len(player_skills))
	is_draw := make([]bool, len(team_skills)-1)
	team_sizes := make([]float64, len(team_skills))
    

	// Team skills orderring
	var prev_min int
	for i := range team_skills {
        team_sizes[i] = 1
		min_place := math.MaxInt
		best_idx := -1
		for j := range team_places {
			if(team_places[j]<min_place){
				best_idx = j
				min_place = team_places[j]
			}
		}
		if(i>0) {
			is_draw[i-1] = prev_min == min_place
		}
		prev_min = min_place 
		team_places[best_idx] = math.MaxInt
        team_order[i] = best_idx
		for j := range team_ids {
			if(team_ids[j] == best_idx) {
				if(team_skills[i] == nil){
					team_skills[i] = player_skills[j]
				} else {
					team_skills[i] = AddGaussian(team_skills[i],player_skills[j])
                    team_sizes[i]+=1
				}
				player_pos[j] = i
			}
		}

	}
	draw_margins := make([]float64, len(team_skills)-1)

    for i := range draw_margins {
       draw_margins[i] = ppf((DRAW_PROBABILITY+1)/2)*BETA*math.Sqrt(team_sizes[i]+team_sizes[i+1]);
    }

    samplers_winner := make([]*Gaussian, len(team_skills))
    samplers_looser := make([]*Gaussian, len(team_skills))
    get_head_sampler := func(idx int) (*Gaussian, *Gaussian, *Gaussian) {
            winner_skill := MultGaussian(team_skills[idx-1], samplers_winner[idx-1])
            looser_skill := MultGaussian(team_skills[idx], samplers_looser[idx])
            prior := SubGaussian(winner_skill, looser_skill)
            var posterior *Gaussian
            if is_draw[idx-1] {
                posterior = propogateDraw(prior, draw_margins[idx-1])
            } else {
                posterior = propogateWin(prior, draw_margins[idx-1])
            }
            return DivGaussian(posterior, prior), winner_skill, looser_skill
    }
    for i := range team_skills{
        samplers_winner[i] = &Gaussian{pi:0, tau:0}
        samplers_looser[i] = &Gaussian{pi:0, tau:0}
    }
    
    
    for j:=0; j<10; j++ {
        max_delta := 0.
        // Right team update
        for i:=1; i<len(team_skills)-1; i++ {
            head_sampler, winner_skill, _ := get_head_sampler(i)

            sampler := SubGaussian(winner_skill, head_sampler)
            max_delta = max(max_delta, distance(sampler, samplers_winner[i]))
            samplers_winner[i] = sampler
        }

        // Left team update
        for i:=len(team_skills)-1; i>1; i-- {
            head_sampler, _, looser_skill := get_head_sampler(i)

            sampler := AddGaussian(looser_skill, head_sampler)
            max_delta = max(max_delta, distance(sampler, samplers_looser[i-1]))
            samplers_looser[i-1] = sampler
        }

        if(max_delta<0.002){
            break
        }
    }
    head_sampler, _, looser_skill := get_head_sampler(1)
    sampler := AddGaussian(looser_skill, head_sampler)
    samplers_looser[0] = sampler

    last:=len(team_skills)-1
    head_sampler, winner_skill, _ := get_head_sampler(last)
    sampler = SubGaussian(winner_skill, head_sampler)
    samplers_winner[last] = sampler

	new_player_skills := make([]*Gaussian, len(player_skills))
	for i := range player_skills {
        team_pos := player_pos[i]		
        sampler := MultGaussian(samplers_winner[team_pos], samplers_looser[team_pos])
        for j := range player_skills {
            if(j!=i && team_order[team_pos] == team_ids[j]) {
                sampler = SubGaussian(sampler, player_skills[j])
            }
        }
        perf_sampler := NewGaussian(sampler.mu(), math.Sqrt(sampler.sigma()*sampler.sigma() + BETA*BETA))
        new_player_skills[i] = MultGaussian(prior_skills[i], perf_sampler)
	}

    return new_player_skills
}
