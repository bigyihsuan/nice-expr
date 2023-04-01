use crate::lexer::tok::Token;

pub type Program = Vec<Expr>;

#[derive(Debug)]
pub enum Expr {
    Literal(Literal),
    Identifier(String),
    Unary { op: Token, expr: Box<Expr> },
    Declaration(Declaration),
    Assignment(Assignment),
}

#[derive(Debug)]
pub enum Declaration {
    Const {
        name: String,
        type_name: Type,
        expr: Box<Expr>,
    },
    Var {
        name: String,
        type_name: Type,
        expr: Box<Expr>,
    },
}

#[derive(Debug)]
pub struct Assignment {
    pub name: String,
    pub op: AssignmentOperator,
    pub expr: Box<Expr>,
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

#[derive(Debug, Clone, PartialEq)]
pub enum Type {
    None,
    Int,
    Dec,
    Str,
    Bool,
    List(Box<Type>),
    Map(Box<Type>, Box<Type>),
}

#[derive(Debug, Clone)]
pub enum AssignmentOperator {
    Invalid,
    Is,
    // TODO: the other assignment operators
}
