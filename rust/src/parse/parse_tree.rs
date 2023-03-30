use std::collections::HashMap;

pub type Program = Vec<Expr>;

#[derive(Debug)]
pub enum Expr {
    Literal(Literal),
}

#[derive(Debug)]
pub enum Literal {
    Int(i64),
    Dec(f64),
    Str(String),
    Bool(bool),
    List(Vec<Literal>),
    Map(HashMap<Literal, Literal>),
}
