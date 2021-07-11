
use crate::data::*;
use std::collections::HashSet;

use std::iter::Iterator;
use rand::{Rng, SeedableRng};

pub fn get_dislike(problem:&Problem, answer:&Vec<Point>) -> i64 {
    let mut result = 0;
    for hole in &problem.hole {
        let mut min = i64::MAX;
        for a in answer {
            let d = get_d(a, &hole);
            if d < min { min = d; }
        }
        result += min;
    }
    result
}

pub fn get_d(a:&Point, b:&Point) -> i64 {
    let x = a.0 - b.0;
    let y = a.1 - b.1;
    x * x + y * y
}


pub fn intersects(a:&Point, b:&Point, c:&Point, d:&Point) -> bool {
    let ax = a.0; let ay = a.1;
    let bx = b.0; let by = b.1;
    let cx = c.0; let cy = c.1;
    let dx = d.0; let dy = d.1;
    let s = (ax - bx) * (cy - ay) - (ay - by) * (cx - ax);
    let t = (ax - bx) * (dy - ay) - (ay - by) * (dx - ax);
    if s * t >= 0 { return false; }
    
    let s  = (cx - dx) * (ay - cy) - (cy - dy) * (ax - cx);
    let t  = (cx - dx) * (by - cy) - (cy - dy) * (bx - cx);
    if s * t >= 0 { return false; }
    true
}

pub fn get_unmatched(problem:&Problem, answer:&Vec<Point>)->i64 {
    let mut result = 0;
    for (ei, edge) in problem.edges.iter().enumerate()
    {
        let ad = get_d(&answer[edge.0], &answer[edge.1]);
        let pd = problem.distances[ei];
        if !check_epsilon(problem, ad, pd) {
            result += 1;
        }
    }
    result
}

pub fn get_not_included(problem:&Problem, answer:&Vec<Point>) -> i64 {
    let mut result = 0;

    // point 
    for a in answer {
        if !includes(problem, a) {
            result += 1;
        }
    }

    // edge
    let mut h0 = problem.hole[problem.hole.len() - 1].clone();
    for h1 in &problem.hole {
        for edge in &problem.edges {
            if intersects(
                &h0, 
                h1, 
                &answer[edge.0],
                &answer[edge.1] 
            ) {
                result += 1;
            }
        }
        h0 = h1.clone();
    }
    result
}

pub fn get_bonus_count(problem:&Problem, answer:&Vec<Point>) -> i64 {
    let mut result = 0;
    // point 
    for bonus in &problem.bonuses {
        for a in answer {
            if *a == *bonus {
                result += 1;
            }
        }
    }
    result
}

pub fn get_not_included_point(problem:&Problem, answer:&Vec<Point>) -> i64 {
    let mut result = 0;

    // point 
    for a in answer {
        if !includes(problem, a) {
            result += 1;
        }
    }
    result
}

pub fn includes(problem:&Problem, point:&Point) -> bool {
    let x = point.0;
    let y = point.1;
    let mut count = 0;
    let mut h0 = problem.hole[problem.hole.len() - 1].clone();
    for h1 in &problem.hole {
        let mut x0 = h0.0 - x;
        let mut y0 = h0.1 - y;
        let mut x1 = h1.0 - x;
        let mut y1 = h1.1 - y;
        
        let cv = x0 * x1 + y0 * y1;
        let sv = x0 * y1 - x1 * y0;
        
        if sv == 0 && cv <= 0 {
            return true;
        }
        
        if y0 < y1 {
            let tmp = x0;
            x0 = x1;
            x1 = tmp;
            let tmp = y0;
            y0 = y1;
            y1 = tmp;
        }
            
        if y1 <= 0 && 0 < y0 {
            let a = x0 * (y1 - y0);
            let b = y0 * (x1 - x0);
            if b < a {
                count += 1;
            }
        }
        h0 = h1.clone();
    }
    return count % 2 != 0;
}

pub fn check_epsilon(problem:&Problem, ad:i64, pd:i64) -> bool {
    let e = problem.epsilon;
    if ad < pd { -(1000000 * ad) <= (e - 1000000) * pd } else { (1000000 * ad) <= (e + 1000000) * pd }
}

pub fn get_center(points:&Vec<Point>) -> Point {
    let mut left   = i64::MAX;
    let mut right  = i64::MIN;
    let mut top    = i64::MAX;
    let mut bottom = i64::MIN;
    for point in points {
        if left   > point.0 { left   = point.0; }
        if right  < point.0 { right  = point.0; }
        if top    > point.1 { top    = point.1; }
        if bottom < point.1 { bottom = point.1; }
    }
    Point((left + right) / 2, (top + bottom) / 2)
}

pub fn lock_points<R: Rng + ?Sized>(locked_points:&mut HashSet<usize>, targets:&Vec<Point>, vertecies:&Vec<Point>, rng:&mut R, rate:f64) {
    for hole in targets {
        if rng.gen_bool(rate) {
            for (pi, p) in vertecies.iter().enumerate() {
                if hole == p {
                    locked_points.insert(pi as usize);
                    break;
                }
            }
        }
    }
}
