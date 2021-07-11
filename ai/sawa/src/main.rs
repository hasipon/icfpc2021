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
        //if cleared.contains(&i) { continue; }
        let target = format!("{}", i);
        let file = File::open(format!("../../problems/{}", target))?;
        let reader = BufReader::new(file);
        let problem:ProblemSource = serde_json::from_reader(reader).unwrap();
        let result = solve(&problem);
        
        println!("{}", target);

        let meta = json!({ "valid": result.best.is_valid(), "score": result.best.get_score(), "dislike": result.best.dislike, "bonus": result.best.bonus_count });
        let answer = json!({ "vertices":result.best.answer.clone() });
        
        println!("{}", meta);
        println!("{}", answer);

        if result.best.is_valid() && result.best.dislike < result.best_bonus.dislike {
            println!("best!");
            let mut file = File::create(format!("out/{}-sawa-auto36-{}.json", target, name))?;
            write!(file, "{}", answer);
            let mut file = File::create(format!("out/{}-sawa-auto36-{}.meta", target, name))?;
            write!(file, "{}", meta);
        }
        
        let meta = json!({ "valid": result.best_bonus.is_valid(), "score": result.best_bonus.get_score(), "dislike": result.best_bonus.dislike, "bonus": result.best_bonus.bonus_count });
        let answer = json!({ "vertices":result.best_bonus.answer.clone() });
        
        println!("{}", meta);
        println!("{}", answer);

        if result.best_bonus.is_valid() && result.best_bonus.bonus_count > 0 {
            println!("best_bonus!");
            let mut file = File::create(format!("out/{}-sawa-auto36-bonus-{}.json", target, name))?;
            write!(file, "{}", answer);
            let mut file = File::create(format!("out/{}-sawa-auto36-bonus-{}.meta", target, name))?;
            write!(file, "{}", meta);
        }
    }
    Ok(())
}
