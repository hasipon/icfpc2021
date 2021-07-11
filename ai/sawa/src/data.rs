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
    pub unmatched:i64,
}

impl State {
    pub fn new(problem:&Problem, answer:Vec<Point>) -> State {
        State {
            dislike     : get_dislike(problem, &answer),
            unmatched   : get_unmatched(problem, &answer),
            not_included: get_not_included(problem, &answer),
            answer:  answer,
        }
    }
    pub fn is_valid(&self) -> bool {
        self.unmatched == 0 && 
        self.not_included == 0

    }
    pub fn get_score(&self) -> i64 {
        let penalty = self.not_included * 30 + self.unmatched * 5;
        self.dislike + penalty + (self.not_included * self.dislike) / 500
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