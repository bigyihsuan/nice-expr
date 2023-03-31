pub mod ast;

use crate::{grammar, lexer, prelude::Result};

use ast::Program;

pub fn parse(filename: Option<std::path::PathBuf>, input: &str) -> Result<Program> {
    let stream = lexer::TokenStream::new(filename, input)?;

    println!("{:?}", &stream.tokens());

    let ast = grammar::module_parser::program(&stream)?;
    Ok(ast)
}
