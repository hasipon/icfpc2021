use std::env::args;
use std::fs::File;
use std::io::BufReader;
use std::io::Write;
use serde_json::json;

mod data;
mod solve;
mod util;
use data::*;
use solve::solve;

fn main()  -> std::io::Result<()>  {
    let arg:Vec<String> = args().collect();
    let cleared = vec![
        4, 11, 12,13,15,16,17,18,20,21,22,23,24,25,26,34,35,38,39,41,43,
        46,49,51,52,53,54,55,59,63,70,72,73,76
    ];
    for i in 1..80 {
        if cleared.contains(&i) { continue; }
        let target = format!("{}", i);
        println!("{}", target);
        let file = File::open(format!("../../problems/{}", target))?;
        let reader = BufReader::new(file);
        let problem:ProblemSource = serde_json::from_reader(reader).unwrap();
        let result = solve(&problem);
        

        let meta = json!({ "valid": result.is_valid(), "score": result.get_score(), "dislike": result.dislike });
        let answer = json!({ "vertices":result.answer.clone() });
        println!("{}", meta);
        println!("{}", answer);
        if result.is_valid() {
            let mut file = File::create(format!("out/{}-sawa-auto19.json", target))?;
            write!(file, "{}", answer);
            let mut file = File::create(format!("out/{}-sawa-auto19.meta", target))?;
            write!(file, "{}", meta);
        }
    }
    Ok(())
}
