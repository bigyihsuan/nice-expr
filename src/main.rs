use crate::{args::parse_args, lexer::TokenStream};

mod args;
mod grammar;
mod lexer;
mod parse;
mod prelude;

fn main() {
    let (file, source) = parse_args();
    eprintln!("file: {file:?}\n```\n{source}```\n");

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
                Ok(ast) => println!("{:#?}", ast),
                Err(err) => eprintln!("error during parsing: {err}"),
            }
        }
        Err(err) => err
            .iter()
            .for_each(|err| eprintln!("error during lexing:\n    {err}")),
    }
}
