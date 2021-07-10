use serde::{Deserialize, Serialize};
use crate::util::*;

#[derive(Serialize, Deserialize, Clone)]
pub struct Problem {
    pub hole: Vec<Point>,
    pub epsilon: i64,
    pub figure: Figure,
}

#[derive(Serialize, Deserialize, Clone)]
pub struct Figure {
    pub edges   :Vec<(i64, i64)>,
    pub vertices:Vec<Point>,
}

#[derive(Serialize, Deserialize, Clone)]
pub struct Point(pub i64, pub i64);

#[derive(Serialize, Deserialize, Clone)]
pub struct Answer {
    pub vertices:Vec<Point>,
}

pub struct SolveResult {
    pub answer:Answer,
    pub dislike:i64,
    pub valid  :bool,
}

#[derive(Clone)]
pub struct State {
    pub answer:Vec<Point>,
    pub dislike:i64,
    pub interrupted:i64,
    pub unmatched:i64,
}

impl State {
    pub fn new(problem:&Problem, answer:Vec<Point>) -> State {
        State {
            dislike    : get_dislike(problem, &answer),
            interrupted: get_interrupted(problem, &answer),
            unmatched  : get_unmatched(problem, &answer),
            answer:  answer,
        }
    }
    pub fn is_valid(&self) -> bool {
        self.interrupted == 0 && self.unmatched == 0
    }
}
