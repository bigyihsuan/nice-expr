#![feature(box_patterns)]
#![feature(let_chains)]

use crate::{args::parse_args, eval::intepreter::Interpreter, lexer::TokenStream};

mod args;
mod eval;
mod grammar;
mod lexer;
mod parse;
mod prelude;
mod util;

fn main() {
    let (file, source) = parse_args();
    println!("file: {file:?}\n```\n{source}```\n");

    let token_stream = TokenStream::new(file, &source);
    match token_stream {
        Ok(token_stream) => {
            let tokens = token_stream.tokens();
            for token in tokens {
                let token = &token.0;
                println!("{token}");
            }
            println!();

            let ast = grammar::module_parser::program(&token_stream);
            match ast {
                Ok(ast) => {
                    println!("{ast:#?}");
                    println!();
                    let interpeter = Interpreter {};
                    let env = Interpreter::default_env();
                    let values = interpeter.interpret_program(&ast, &env);
                    match values {
                        Ok(values) => {
                            for value in values {
                                println!("{}", interpeter.format_value(&value));
                            }
                            println!();

                            for (name, value_entry) in env.borrow().identifiers() {
                                let value = value_entry.v;
                                let con = value_entry.c;
                                let t = value_entry.t;
                                println!(
                                    "{name}:{}, {con:?} {t:?}",
                                    interpeter.format_value(&value)
                                )
                            }
                        }
                        Err(err) => eprintln!("error during evaluation: {err:?}"),
                    }
                }
                Err(err) => eprintln!("error during parsing: {err}"),
            }
        }
        Err(err) => eprintln!("error during lexing:\n    {err}"),
    }
}
