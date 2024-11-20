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
    return propogateExpectation(g, draw_margin, math.Inf(1))
}
func propogateDraw(g *Gaussian, draw_margin float64) *Gaussian{
    return propogateExpectation(g, -draw_margin, draw_margin)
}
func propogateExpectation(g *Gaussian, low float64, high float64) *Gaussian{
    avg := 0.
    avg_squared := 0.
    idx := 0.
    for range 800000 {
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

    //team_skills := []*Gaussian {NewGaussian(4,2), NewGaussian(3.5,1), NewGaussian(3,1.5)}
    player_skills := []*Gaussian {NewGaussian(1,5.135), NewGaussian(3,5.135), NewGaussian(6,5.135), NewGaussian(14,5.135)}
	team_ids := []int {0,1,2,1}
	team_places := []int {0,2,1}
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
	fmt.Println(team_places)

	// Team skills orderring
	var prev_min int
	for i := range team_skills {
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
				}
				player_pos[j] = i
			}
		}

	}
	fmt.Println(player_skills)
	fmt.Println(team_ids)
	fmt.Println(team_skills)
	fmt.Println(player_pos)
	fmt.Println(is_draw)
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
            max_delta = max(max_delta, distance(sampler, samplers_winner[i]))
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
            max_delta = max(max_delta, distance(sampler, samplers_looser[i-1]))
            samplers_looser[i-1] = sampler
        }
        //prts()
		fmt.Printf("Max delta %f\n", max_delta)
        if(max_delta<0.02){
            break
        }
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

	//player_samples := make([]*Gaussian, len(player_skills))
	new_player_skills := make([]*Gaussian, len(player_skills))
    //fmt.Println(player_pos)
	for i := range player_skills {
        team_pos := player_pos[i]		
        sampler := MultGaussian(samplers_winner[team_pos], samplers_looser[team_pos])
        fmt.Printf("Player %d\n", i)
        for j := range player_skills {
            if(j!=i && team_order[team_pos] == team_ids[j]) {
                fmt.Printf("  Adding player %d\n", j)
                sampler = AddGaussian(sampler, player_skills[j])
            }
        }
        new_player_skills[i] = MultGaussian(player_skills[i], sampler)
        //player_skill := player_skills[i]
	}
    fmt.Println(player_skills)
    fmt.Println(new_player_skills)

   //fmt.Println(NewGaussian(avg, 
}
