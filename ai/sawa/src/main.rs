use std::env::args;
use std::fs::File;
use std::io::BufReader;
use serde_json::json;

mod data;
mod solve;
mod util;
use data::*;
use solve::solve;

fn main()  -> std::io::Result<()>  {
    let arg:Vec<String> = args().collect();
    let mut target = "8".to_owned();
    if arg.len() > 1 {
        target = arg[1].to_owned();
    }
    
    let mut file = File::open(format!("../../problems/{}", target))?;
    let reader = BufReader::new(file);
    let problem:Problem = serde_json::from_reader(reader).unwrap();

    let result = solve(&problem);
    
    println!("{}", json!(result.dislike));
    println!("{}", json!(result.get_score()));
    println!("{}", json!(result.is_valid()));
    println!("{}", json!(Answer{ vertices:result.answer }));

    Ok(())
}
