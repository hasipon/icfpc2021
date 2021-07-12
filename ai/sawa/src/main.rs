use std::env::args;
use std::fs::File;
use std::io::BufReader;
use std::io::Write;
use serde_json::json;
use std::collections::HashMap;

mod data;
mod solve;
mod util;
mod operation;
use data::*;
use solve::solve;
use std::fs::DirEntry;

fn main()  -> std::io::Result<()>  {
    let arg:Vec<String> = args().collect();
    let cleared = vec![
        4, 11, 12,13,15,16,17,18,20,21,22,23,24,25,26,34,35,38,39,41,43,
        46,47,49,51,52,53,54,55,59,63,65,70,72,73,75,76,77,80,84,90,97,106
    ];
    let mut name  = "x".to_owned();
    let mut start = 1;
    let mut end   = 133;
    if arg.len() == 4 {
        name  = arg[1].to_owned();
        start = arg[2].parse().unwrap();
        end   = arg[3].parse().unwrap();
    }

    let mut inputs = HashMap::new();

    for j in 0..1000000 {
        let i = (j + start) % 132 + 1;

        if cleared.contains(&i) { continue; }
        let target = format!("{}", i);
        
        if !inputs.contains_key(&i) {
            let file = File::open(format!("../../problems/{}", target))?;
            let reader = BufReader::new(file);
            let problem:ProblemSource = serde_json::from_reader(reader).unwrap();
            let mut vertices = Vec::new();
            vertices.push(problem.figure.vertices.clone());
            
            for file in std::fs::read_dir("../../solutions")? {
                read_vertices(&mut vertices, &file?, i, problem.figure.vertices.len());
            }
            inputs.insert(i, (problem, vertices));
        }
        let input = inputs.get(&i).unwrap();
        let result = solve(&input.0, &input.1);
        
        println!("{}", target);

        let meta = json!({ "valid": result.best.is_valid(), "score": result.best.get_score(), "dislike": result.best.dislike, "bonus": result.best.bonus_count });
        let answer = json!({ "vertices":result.best.answer.clone() });
        
        println!("{}", meta);
        println!("{}", answer);

        if result.best.is_valid() {
            println!("best!");
            let mut file = File::create(format!("out/{}-sawa-auto41-{}-{}.json", target, j, name))?;
            write!(file, "{}", answer);
            let mut file = File::create(format!("out/{}-sawa-auto41-{}-{}.meta", target, j, name))?;
            write!(file, "{}", meta);
        }
        
        let meta = json!({ "valid": result.best_bonus.is_valid(), "score": result.best_bonus.get_score(), "dislike": result.best_bonus.dislike, "bonus": result.best_bonus.bonus_count });
        let answer = json!({ "vertices":result.best_bonus.answer.clone() });
        
        println!("{}", meta);
        println!("{}", answer);

        if result.best_bonus.is_valid() && result.best_bonus.bonus_count > 0 {
            println!("best_bonus!");
            let mut file = File::create(format!("out/{}-sawa-auto41-bonus-{}-{}.json", target, j, name))?;
            write!(file, "{}", answer);
            let mut file = File::create(format!("out/{}-sawa-auto41-bonus-{}-{}.meta", target, j, name))?;
            write!(file, "{}", meta);
        }
    }
    Ok(())
}

fn read_vertices(vertices:&mut Vec<Vec<Point>>, path:&DirEntry, index:usize, len:usize) -> std::io::Result<()> {
    let file_name = path.file_name().into_string().unwrap();
    let prefix:String = format!("{}-", index);
    
    if path.file_type()?.is_file() && file_name.starts_with(&prefix) {
        let file = File::open(path.path())?;
        let reader = BufReader::new(file);
        let answer:Answer = serde_json::from_reader(reader)?;
        if answer.vertices.len() == len {
            vertices.push(answer.vertices);
        }
    }
    Ok(())
}
