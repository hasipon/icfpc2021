
use crate::data::*;

pub fn get_dislike(problem:&Problem, answer:&Vec<Point>) -> i64 {
    let result = 0;
    for hole in problem.hole {
        let min = i64::MAX;
        for a in answer {
            let d = get_d(a, b);
            if d < min { min = d; }
        }
        
    }
    result
}

pub fn get_d(a:Point, b:Point) -> i64 {
    let x = a.0 - b.0;
    let y = a.1 - b.1;
    x * x + y * y
}