
use crate::data::*;
use crate::util::*;
use crate::operation::*;
use rand::rngs::SmallRng;
use rand::{Rng, SeedableRng};
use std::collections::HashSet;
use std::mem;
use serde_json::json;

pub fn solve(source:&ProblemSource) -> SolveResult {
    let mut distances = Vec::new();
    for edge in &source.figure.edges {
        distances.push(get_d(&source.figure.vertices[edge.0], &source.figure.vertices[edge.1]));
    }
    let mut bonuses = Vec::new();
    for bonus in &source.bonuses {
        bonuses.push(bonus.position.clone());
    }

    let mut left   = i64::MAX;
    let mut right  = i64::MIN;
    let mut top    = i64::MAX;
    let mut bottom = i64::MIN;
    for point in &source.hole {
        if left   > point.0 { left   = point.0; }
        if right  < point.0 { right  = point.0; }
        if top    > point.1 { top    = point.1; }
        if bottom < point.1 { bottom = point.1; }
    }
    let center = Point((left + right) / 2, (top + bottom) / 2);
    let problem = Problem {
        hole: source.hole.clone(),
        edges: source.figure.edges.clone(),
        epsilon: source.epsilon,
        center,
        left,
        right,
        top,
        bottom,
        distances,
        bonuses
    };
    let current = State::new(&problem, source.figure.vertices.clone());
    let mut best       = current.clone();
    let mut best_bonus = current.clone();
    let mut rng = SmallRng::from_entropy();
    let mut arr0 = Vec::new();
    let mut arr1 = Vec::new();

    let size = 3000;
    let mut locked_points = HashSet::new();

    for i in 0..size {
        let mut vertecies = current.answer.clone();
        if rng.gen_bool(0.2) { translate(&problem, &mut vertecies, &mut rng, &locked_points); }
        if rng.gen_bool(0.2) { inverse_x(&problem, &mut vertecies, &locked_points); }
        if rng.gen_bool(0.2) { inverse_y(&problem, &mut vertecies, &locked_points); }
        if i > 50 { 
            random(&problem, &mut vertecies, 1, &mut rng, &locked_points); 
        }
        arr0.push(State::new(&problem, vertecies));
    }

    let mut prev_score = 1200003;
    let mut prev_dislike = 1200003;
    let repeat = 190;

    for i in 0..repeat {
        arr0.sort();
        arr0.split_off(size);
        if arr1.len() > size / 5 { arr1.split_off(size / 5); }
        let scale = (repeat - i) as f64 / repeat as f64;
        //println!("{}: {} {} {}", i, arr0[0].is_valid(), arr0[0].dislike, arr0[0].get_score());

        for current in &arr0 {
            for _ in 0..2 {
                let score = current.get_score();
                let dislike = current.dislike;
                
                let mut vertecies = current.answer.clone();

                locked_points.clear();
                if rng.gen_bool(0.5) {
                    let rate = if rng.gen_bool(0.5) { 1.0 - scale * 0.2 } else { rng.gen_range(0.1, 1.0) };
                    lock_points(
                        &mut locked_points, 
                        &problem.hole, 
                        &vertecies, 
                        &mut rng, 
                        rate
                    );
                }
                //if rng.gen_bool(0.1) {
                //    let rate = if rng.gen_bool(0.5) { 1.0 - scale * 0.2 } else { rng.gen_range(0.1, 1.0) };
                //    lock_points(
                //        &mut locked_points, 
                //        &problem.bonuses, 
                //        &vertecies, 
                //        &mut rng, 
                //        rate
                //    );
                //}

                if rng.gen_bool(0.3 * scale) || score == prev_score && dislike == prev_dislike && rng.gen_bool(1.0 * scale) { 
                    random(&problem, &mut vertecies, 1, &mut rng, &locked_points); 
                }

                if rng.gen_bool(0.3 * scale) { translate(&problem, &mut vertecies, &mut rng, &locked_points); }
                else if rng.gen_bool(0.1) { translate_small(&mut vertecies, &mut rng, scale, &locked_points); }
                if rng.gen_bool(0.2 * scale) { inverse_x(&problem, &mut vertecies, &locked_points); }
                if rng.gen_bool(0.2 * scale) { inverse_y(&problem, &mut vertecies, &locked_points); }
                if rng.gen_bool(0.2 * scale) { rotate(&problem, &mut vertecies, &mut rng, scale, &locked_points); }
                

                if rng.gen_bool(0.5 ) { random_include(&problem, &mut vertecies, &mut rng, &locked_points); }
                if rng.gen_bool(0.7 ) { fit           (&problem.hole   , &mut vertecies, &mut rng, scale); }
                if rng.gen_bool(0.01) { fit           (&problem.bonuses, &mut vertecies, &mut rng, 1.0); }
                if rng.gen_bool(0.1 ) { random_small  (&mut vertecies, 1, &mut rng, &locked_points); }

                pull(&problem, &mut vertecies, 40, &mut rng, &locked_points);
                if rng.gen_bool(0.5) && get_unmatched(&problem, &mut vertecies) == 0 {
                    search_include(&problem, &mut vertecies, &mut rng, &locked_points);
                }
                
                let next = State::new(&problem, vertecies);
                if next.bonus_count > 0 {
                    if 
                        (!best_bonus.is_valid() && next.is_valid()) || 
                        (best_bonus.is_valid() == next.is_valid() && best_bonus.get_score() >= next.get_score())
                    {
                        best_bonus = next.clone()
                    }
                } else {
                    if 
                        (!best.is_valid() && next.is_valid()) || 
                        (best.is_valid() == next.is_valid() && best.get_score() >= next.get_score())
                    {
                        best = next.clone()
                    }
                }
                arr1.push(next);
                prev_score = score;
                prev_dislike = current.dislike;
            }
        }
        mem::swap(&mut arr0, &mut arr1);
    }

    SolveResult{ best, best_bonus }
}
