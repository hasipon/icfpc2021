use crate::data::*;
use crate::util::*;
use rand::Rng;
use std::iter::Iterator;

pub fn translate<R: Rng + ?Sized>(problem:&Problem, answer:&mut Vec<Point>, rng: &mut R) {
    let center = get_center(answer);
    
    let dx = rng.gen_range(problem.center.0.min(center.0), problem.center.0.max(center.0) + 1) - center.0;
    let dy = rng.gen_range(problem.center.1.min(center.1), problem.center.1.max(center.1) + 1) - center.1;
    for a in answer {
        a.0 += dx;
        a.1 += dy;
    }
}

pub fn translate_small<R: Rng + ?Sized>(answer:&mut Vec<Point>, rng: &mut R, scale:f64) {
    let d = (scale * 20.0).ceil() as i64;
    let dx = rng.gen_range(-d, d);
    let dy = rng.gen_range(-d, d);
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

pub fn rotate<R: Rng + ?Sized>(problem:&Problem, answer:&mut Vec<Point>, rng: &mut R, scale: f64) {
    let d = rng.gen_range(-std::f64::consts::PI, std::f64::consts::PI) * scale;
    let sin = d.sin();
    let cos = d.cos();
    for a in answer {
        let x = (a.0 - problem.center.0) as f64;
        let y = (a.1 - problem.center.1) as f64;
        a.0 = (x * cos - y * sin) as i64 + problem.center.0;
        a.1 = (x * sin + y * cos) as i64 + problem.center.1;
    }
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
                let v = (adf.sqrt() - pdf.sqrt()) / 3.5;
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
                answer[i].0 = a0.round() as i64;
                answer[i].1 = a1.round() as i64;
            }
        }
    }
}

pub fn fit<R: Rng + ?Sized>(problem:&Problem, answer:&mut Vec<Point>, rng: &mut R, scale:f64) {
    for hole in &problem.hole {
        if rng.gen_bool(scale) { continue; }

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
            let v = (min as f64).sqrt() * rng.gen_range(scale, 1.0);
            let mut a = &answer[target];
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

pub fn random<R: Rng + ?Sized>(problem:&Problem, answer:&mut Vec<Point>, repeat:i64, rng: &mut R) {
    for i in 0..repeat {
        for hole in &problem.hole {
            let i = rng.gen_range(0, answer.len());
            
            let a = &answer[i];
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

pub fn random_include<R: Rng + ?Sized>(problem:&Problem, answer:&mut Vec<Point>, rng: &mut R) {
   for a in answer {
        a.0 = problem.left  .max(a.0);
        a.0 = problem.right .min(a.0);
        a.1 = problem.top   .max(a.1);
        a.1 = problem.bottom.min(a.1);
        if includes(problem, a) {
            continue;
        }

        let x = rng.gen_range(problem.left.max(a.0 - 10), problem.right .min(a.0 + 10));
        let y = rng.gen_range(problem.top .max(a.1 - 10), problem.bottom.min(a.1 + 10));
        if includes(problem, &Point(x, y)) {
            a.0 = x;
            a.1 = y;
        }
        let x = rng.gen_range(problem.left.max(a.0 - 40), problem.right .min(a.0 + 40));
        let y = rng.gen_range(problem.top .max(a.1 - 40), problem.bottom.min(a.1 + 40));
        if includes(problem, &Point(x, y)) {
            a.0 = x;
            a.1 = y;
        }
        let x = rng.gen_range(problem.left.max(a.0 - 80), problem.right .min(a.0 + 80));
        let y = rng.gen_range(problem.top .max(a.1 - 80), problem.bottom.min(a.1 + 80));
        if includes(problem, &Point(x, y)) {
            a.0 = x;
            a.1 = y;
        }
   }
}

pub fn search_include<R: Rng + ?Sized>(problem:&Problem, answer:&mut Vec<Point>, rng: &mut R) {
    for ai in 0..answer.len() {
        let a = answer[ai].clone();
        if rng.gen_bool(0.5) && includes(problem, &a) { continue; }
        let mut p = a.clone();
        let mut related_edges = Vec::new();

        for (ei, edge) in problem.edges.iter().enumerate() {
            if edge.0 == ai { related_edges.push((ei, answer[edge.1].clone())); } 
            if edge.1 == ai { related_edges.push((ei, answer[edge.0].clone())); } 
        }
        for _ in 0..5 {
            p.0 = rng.gen_range(problem.left, problem.right );
            p.1 = rng.gen_range(problem.top , problem.bottom);
            
            let mut success = false;
            for _ in 0..6 {
                let mut failed = false;
                for (ei, point) in &related_edges {
                    let ad = get_d(&p, point);
                    let pd = problem.distances[*ei as usize];
                    let v = ((ad as f64).sqrt() - (pd as f64).sqrt()) * 0.85;
                    if check_epsilon(problem, ad, pd) { continue; }
                    failed = true;
                    let ax = (p.0 - point.0) as f64;
                    let ay = (p.1 - point.1) as f64;
                    let d = ay.atan2(ax);
                    p.0 -= (v * d.cos() + rng.gen_range(-0.5, 0.5)).round() as i64;
                    p.1 -= (v * d.sin() + rng.gen_range(-0.5, 0.5)).round() as i64;
                }

                if !failed {
                    success = true;
                    break;
                }
            }

            if success && includes(problem, &p) {
                answer[ai].0 = p.0;
                answer[ai].1 = p.1;
                break;
            }
        }
        
   }
}
