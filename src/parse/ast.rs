use crate::lexer::tok::Token;

pub type Program = Vec<Expr>;

#[derive(Debug)]
pub enum Expr {
    Literal(Literal),
    Identifier { name: String },
    Unary { op: Token, expr: Box<Expr> },
}

#[derive(Debug)]
pub enum Literal {
    Int(i64),
    Dec(f64),
    Str(String),
    Bool(bool),
    List(Vec<Expr>),
    Map(Vec<(Expr, Expr)>),
}
