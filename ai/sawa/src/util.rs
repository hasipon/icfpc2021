
use crate::data::*;
use rand::Rng;
use std::iter::Iterator;

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
pub fn get_intersected(problem:&Problem, answer:&Vec<Point>) -> i64 {
    let mut result = 0;
    let mut h0 = problem.hole[problem.hole.len() - 1];
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

pub fn check_epsilon(problem:&Problem, ad:i64, pd:i64) -> bool {
    let e = problem.epsilon;
    if ad < pd { -(1000000 * ad) <= (e - 1000000) * pd } else { (1000000 * ad) <= (e + 1000000) * pd }
}

pub fn pull<R: Rng + ?Sized>(problem:&Problem, answer:&mut Vec<Point>, repeat:i64, rng: &mut R) {
    for _ in 0..repeat
    {
        let mut count      = Vec::new();
        let mut velocities = Vec::new();
        for _ in 0..answer.len() {
            count.push(0);
            velocities.push((0.0, 0.0));
        }
        let mut matched = true;
        for (ei, edge) in problem.edges.iter().enumerate()
        {
            let ad = get_d(&answer[edge.0], &answer[edge.1]);
            let pd = problem.distances[ei];
            
            if !check_epsilon(problem, ad, pd) {
                count[edge.0] += 1; 
                count[edge.1] += 1; 
                let adf = ad as f64;
                let pdf = pd as f64;
                let v = (adf.sqrt() - pdf.sqrt()) / 5.0;
                let ax = (answer[edge.0].0 - answer[edge.1].0) as f64;
                let ay = (answer[edge.0].1 - answer[edge.1].1) as f64;
                let d = ay.atan2(ax);
                velocities[edge.0].0 -= v * d.cos();
                velocities[edge.0].1 -= v * d.sin();
                velocities[edge.1].0 += v * d.cos();
                velocities[edge.1].1 += v * d.sin();
                matched = false;
            }
        }
        if matched { break; }
        for i in 0..answer.len()
        {
            let v = velocities[i];
            let c = count[i];
            if c != 0 {
                if c == 1 && rng.gen_bool(0.1)  { continue; }
                let a0:f64 = answer[i].0 as f64 + (v.0 / (c + 1) as f64) + rng.gen_range(-0.5, 0.5);
                let a1:f64 = answer[i].1 as f64 + (v.1 / (c + 1) as f64) + rng.gen_range(-0.5, 0.5);
                answer[i] = Point(a0.round() as i64, a1.round() as i64);
            }
        }
    }
}

pub fn fit<R: Rng + ?Sized>(problem:&Problem, answer:&mut Vec<Point>, repeat:i64, rng: &mut R) {
    for _ in 0..repeat
    {
        for hole in &problem.hole {
            let mut min = i64::MAX;
            let mut target = 0;
            for i in 0..answer.len() {
                let d = get_d(&answer[i], hole);
                if 
                    d < min &&
                    (d == 0 || d + 20 < min || rng.gen_bool(0.5))
                {
                    min = d;
                    target = i;
                }
            }
            if min > 0 {
                let v = (min as f64).sqrt();
                let mut a = answer[target];
                let dx = (a.0 - hole.0) as f64;
                let dy = (a.1 - hole.1) as f64;
                let d = dy.atan2(dx);
                answer[target] = Point(
                    (a.0 as f64 - v * d.cos()).round() as i64,
                    (a.1 as f64 - v * d.sin()).round() as i64
                );
            }
        }
    }
}

pub fn random<R: Rng + ?Sized>(problem:&Problem, answer:&mut Vec<Point>, repeat:i64, rng: &mut R) {
    for i in 0..repeat {
        for hole in &problem.hole {
            let i = rng.gen_range(0, answer.len());
            
            let a = answer[i];
            let dx = (a.0 - hole.0) as f64;
            let dy = (a.1 - hole.1) as f64;
            if dx != 0.0 || dy != 0.0 {
                let v = (dx * dx + dy * dy).sqrt();
                let d = dy.atan2(dx);
                answer[i] = Point(
                    ((a.0 as f64 - v * d.cos()) * rng.gen_range(0.0, 1.0) * rng.gen_range(0.0, 1.0) + rng.gen_range(-0.5, 0.5)).round() as i64,
                    ((a.1 as f64 - v * d.sin()) * rng.gen_range(0.0, 1.0) * rng.gen_range(0.0, 1.0) + rng.gen_range(-0.5, 0.5)).round() as i64
                );
            }
        }
    }
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

pub fn translate<R: Rng + ?Sized>(problem:&Problem, answer:&mut Vec<Point>, rng: &mut R) {
    let center = get_center(answer);
    
    let dx = rng.gen_range(problem.center.0.min(center.0) - 10, problem.center.0.max(center.0) + 10) - center.0;
    let dy = rng.gen_range(problem.center.1.min(center.1) - 10, problem.center.1.max(center.1) + 10) - center.1;
    for a in answer {
        a.0 += dx;
        a.1 += dy;
    }
}
pub fn inverse_x(problem:&Problem, answer:&mut Vec<Point>) {
    for a in answer {
        a.0 = problem.center.0 * 2 - a.0;
    }
}

pub fn inverse_y(problem:&Problem, answer:&mut Vec<Point>) {
    for a in answer {
        a.1 = problem.center.1 * 2 - a.1;
    }    
}