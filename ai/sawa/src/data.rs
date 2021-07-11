use serde::{Deserialize, Serialize};
use crate::util::*;
use std::cmp::Ordering;

pub struct Problem {
    pub hole: Vec<Point>,
    pub epsilon: i64,
    pub edges: Vec<(usize, usize)>,
    pub distances: Vec<i64>,
    pub center:Point,
}

#[derive(Serialize, Deserialize, Clone)]
pub struct ProblemSource {
    pub hole: Vec<Point>,
    pub epsilon: i64,
    pub figure: Figure,
}

#[derive(Serialize, Deserialize, Clone)]
pub struct Figure {
    pub edges   :Vec<(usize, usize)>,
    pub vertices:Vec<Point>,
}

#[derive(Serialize, Deserialize, Clone, Copy, Eq, PartialEq)]
pub struct Point(pub i64, pub i64);

#[derive(Serialize, Deserialize, Clone)]
pub struct Answer {
    pub vertices:Vec<Point>,
}

#[derive(Clone, Eq)]
pub struct State {
    pub answer:Vec<Point>,
    pub dislike:i64,
    pub intersected:i64,
    pub unmatched:i64,
}

impl State {
    pub fn new(problem:&Problem, answer:Vec<Point>) -> State {
        State {
            dislike    : get_dislike(problem, &answer),
            intersected: get_intersected(problem, &answer),
            unmatched  : get_unmatched(problem, &answer),
            answer:  answer,
        }
    }
    pub fn is_valid(&self) -> bool {
        self.intersected == 0 && self.unmatched == 0
    }
    pub fn get_score(&self) -> i64 {
        self.dislike + self.intersected * 1000 + self.unmatched * 300
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