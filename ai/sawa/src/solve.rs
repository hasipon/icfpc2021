
use crate::data::*;
use crate::util::*;
use serde_json::json;
use rand::rngs::SmallRng;
use rand::{Rng, SeedableRng};
use std::collections::BinaryHeap;
use std::mem;

pub fn solve(source:&ProblemSource) -> State {
    let mut distances = Vec::new();
    for edge in &source.figure.edges {
        distances.push(get_d(&source.figure.vertices[edge.0], &source.figure.vertices[edge.1]));
    }
    let center = get_center(&source.hole);
    let problem = Problem {
        hole: source.hole.clone(),
        edges: source.figure.edges.clone(),
        epsilon: source.epsilon,
        center,
        distances,
    };
    let current = State::new(&problem, source.figure.vertices.clone());
    let mut best = current.clone();
    let mut rng = SmallRng::from_entropy();
    let mut arr0 = Vec::new();
    let mut arr1 = Vec::new();

    let size = 1200;
    for i in 0..size {
        let mut vertecies = current.answer.clone();
        if rng.gen_bool(0.2) { translate(&problem, &mut vertecies, &mut rng); }
        if rng.gen_bool(0.2) { inverse_x(&problem, &mut vertecies); }
        if rng.gen_bool(0.2) { inverse_y(&problem, &mut vertecies); }
        if i > 50 { 
            random(&problem, &mut vertecies, 1, &mut rng); 
        }
        arr0.push(State::new(&problem, vertecies));
    }

    let mut prev_score = 1200003;
    let mut repeat = 150;
    for i in 0..repeat {
        arr0.sort();
        arr0.split_off(size);
        let scale = (repeat - i) as f64 / repeat as f64;

        for current in &arr0 {
            for _ in 0..2 {
                let score = current.get_score();

                let mut vertecies = current.answer.clone();
                if rng.gen_bool(0.3 * scale) || score == prev_score && rng.gen_bool(1.0 * scale) { 
                    random(&problem, &mut vertecies, 1, &mut rng); 
                }
                if rng.gen_bool(0.3 * scale) { translate(&problem, &mut vertecies, &mut rng); }
                if rng.gen_bool(0.2 * scale) { inverse_x(&problem, &mut vertecies); }
                if rng.gen_bool(0.2 * scale) { inverse_y(&problem, &mut vertecies); }
                if rng.gen_bool(0.2 * scale) { rotate(&problem, &mut vertecies, &mut rng, scale); }
                if rng.gen_bool(0.6) { fit (&problem, &mut vertecies, 1, &mut rng); }

                pull(&problem, &mut vertecies, 40, &mut rng);
                

                let next = State::new(&problem, vertecies);
                if 
                    (!best.is_valid() && next.is_valid()) || 
                    (best.is_valid() == next.is_valid() && best.get_score() >= next.get_score())
                {
                    println!("{}: {} {} {}", i, next.is_valid(), next.dislike, next.get_score());
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
