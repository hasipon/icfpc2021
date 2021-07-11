use serde::{Deserialize, Serialize};
use crate::util::*;
use std::cmp::Ordering;

pub struct Problem {
    pub hole: Vec<Point>,
    pub epsilon: i64,
    pub edges: Vec<(usize, usize)>,
    pub distances: Vec<i64>,
    pub center:Point,
    pub left:i64,
    pub right:i64,
    pub top:i64,
    pub bottom:i64,
    pub bonuses: Vec<Point>,
}

#[derive(Serialize, Deserialize, Clone)]
pub struct ProblemSource {
    pub hole: Vec<Point>,
    pub epsilon: i64,
    pub figure: Figure,
    pub bonuses: Vec<Bonus>,
}
#[derive(Serialize, Deserialize, Clone)]
pub struct Bonus {
    pub position:Point,   
}

#[derive(Serialize, Deserialize, Clone)]
pub struct Figure {
    pub edges   :Vec<(usize, usize)>,
    pub vertices:Vec<Point>,
}

#[derive(Serialize, Deserialize, Clone, Eq, PartialEq)]
pub struct Point(pub i64, pub i64);

#[derive(Serialize, Deserialize, Clone)]
pub struct Answer {
    pub vertices:Vec<Point>,
}

#[derive(Clone, Eq)]
pub struct State {
    pub answer:Vec<Point>,
    pub dislike:i64,
    pub not_included:i64,
    pub bonus_count:i64,
    pub unmatched:i64,
    pub len:i64,
}

impl State {
    pub fn new(problem:&Problem, answer:Vec<Point>) -> State {
        State {
            dislike     : get_dislike(problem, &answer),
            unmatched   : get_unmatched(problem, &answer),
            not_included: get_not_included(problem, &answer),
            bonus_count : get_bonus_count(problem, &answer),
            len : (problem.edges.len() + answer.len() + problem.hole.len()) as i64 + (problem.right - problem.left)  + (problem.bottom - problem.top),
            answer:  answer,
        }
    }
    pub fn is_valid(&self) -> bool {
        self.unmatched == 0 && 
        self.not_included == 0
    }
    pub fn get_score(&self) -> i64 {
        let penalty = self.not_included * 50 + self.unmatched * 2 - self.bonus_count * 10;
        self.dislike + penalty * self.len
    }
}

impl Ord for State {
    fn cmp(&self, other: &Self) -> Ordering {
        self.get_score().cmp(&other.get_score())
    }
}

impl PartialOrd for State {
    fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
        Some(self.cmp(other))
    }
}

impl PartialEq for State {
    fn eq(&self, other: &Self) -> bool {
        self.get_score() == other.get_score()
    }
}