#![feature(box_patterns)]
#![feature(let_chains)]

use itertools::Itertools;

use crate::{
    args::parse_args,
    eval::intepreter::{builtin, Interpreter},
    lexer::TokenStream,
    parse::grammar::parser,
};

mod args;
mod eval;
mod lexer;
mod parse;
mod prelude;
mod util;

fn main() {
    let (file, source, debug) = parse_args();

    if debug {
        eprintln!("file: {file:?}\n```\n{source}```\n");
    }

    let token_stream = TokenStream::new(file, &source);
    match token_stream {
        Ok(token_stream) => {
            let tokens = token_stream.tokens();
            if debug {
                for token in tokens {
                    let token = &token.0;
                    eprintln!("{token}");
                }
                eprintln!();
            }

            let ast = parser::program(&token_stream);
            match ast {
                Ok(ast) => {
                    if debug {
                        eprintln!("{ast:#?}");
                        eprintln!();
                    }
                    let interpeter = Interpreter {};
                    let env = Interpreter::default_env();
                    let values = interpeter.interpret_program(&ast, &env);
                    match values {
                        Ok(_) => {
                            if debug {
                                eprintln!();
                            }
                            let mut entries = env.borrow().identifiers().into_iter().collect_vec();
                            entries.sort_by(|(l, _), (r, _)| l.partial_cmp(r).unwrap());
                            for (name, value_entry) in entries {
                                let value = value_entry.v;
                                let con = value_entry.c;
                                let t = value_entry.t;
                                if debug {
                                    eprintln!(
                                        "{name}:{}, {con:?} {t:?}",
                                        builtin::format_value(&value)
                                    )
                                }
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
