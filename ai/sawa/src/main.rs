use std::env::args;
use std::fs::File;
use std::io::BufReader;
use std::io::Write;
use serde_json::json;

mod data;
mod solve;
mod util;
mod operation;
use data::*;
use solve::solve;

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

    for i in start..end {
        if cleared.contains(&i) { continue; }
        let target = format!("{}", i);
        let file = File::open(format!("../../problems/{}", target))?;
        let reader = BufReader::new(file);
        let problem:ProblemSource = serde_json::from_reader(reader).unwrap();
        let result = solve(&problem);
        
        let meta = json!({ "valid": result.is_valid(), "score": result.get_score(), "dislike": result.dislike, "bonus": result.bonus_count });
        let answer = json!({ "vertices":result.answer.clone() });
        
        println!("{}", target);
        println!("{}", meta);
        println!("{}", answer);
        if result.is_valid() {
            let mut file = File::create(format!("out/{}-sawa-auto34-{}.json", target, name))?;
            write!(file, "{}", answer);
            let mut file = File::create(format!("out/{}-sawa-auto34-{}.meta", target, name))?;
            write!(file, "{}", meta);
        }
    }
    Ok(())
}
