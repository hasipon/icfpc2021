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
