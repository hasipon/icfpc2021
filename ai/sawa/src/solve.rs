
use crate::data::*;
use crate::util::*;

pub fn solve(problem:&Problem) -> SolveResult {
    let answer = Vec::new();

    let first_answer = problem.figure.vertices.clone();
    let first = State::new(problem, first_answer);
    let best = first;

    SolveResult {
        valid  : best.is_valid(),
        answer : Answer { vertices: best.answer },
        dislike: best.dislike,
    }
}
