
use crate::data::*;
use crate::util::*;
use serde_json::json;
use rand::rngs::SmallRng;
use rand::{Rng, SeedableRng};
use std::collections::BinaryHeap;
use std::mem;

pub fn solve(problem:&Problem) -> State {

    let current = State::new(problem, problem.figure.vertices.clone());
    let mut best = current.clone();
    let mut rng = SmallRng::from_entropy();
    let mut arr0 = Vec::new();
    let mut arr1 = Vec::new();

    let size = 600;
    for i in 0..size {
        let mut vertecies = current.answer.clone();
        if i > 3 { random(problem, &mut vertecies, 1, &mut rng); }
        arr0.push(State::new(problem, vertecies));
    }

    let mut prev_score = 1200003;
    for _ in 0..200 {
        arr0.sort();
        arr0.split_off(size);

        for current in &arr0 {
            for i in 0..2 {
                let score = current.get_score();

                let mut vertecies = current.answer.clone();
                if rng.gen_bool(0.1) || score == prev_score && rng.gen_bool(0.9) { 
                    random(problem, &mut vertecies, 1, &mut rng); 
                }
                if rng.gen_bool(0.6) { fit (problem, &mut vertecies, 1, &mut rng); }
                pull(problem, &mut vertecies, 40, &mut rng);

                let next = State::new(problem, vertecies);
                if 
                    (!best.is_valid() && next.is_valid()) || 
                    (best.get_score() > next.get_score())
                {
                    best = next.clone()
                }
                arr1.push(next);
                prev_score = score;
            }
        }
        mem::swap(&mut arr0, &mut arr1);
    }

    best
}
