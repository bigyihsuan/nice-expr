use std::{collections::HashMap, hash::Hash};

use crate::lexer::tok::Token;

pub type Program = Vec<Expr>;

#[derive(Debug, PartialEq, Eq)]
pub enum Expr {
    Literal(Literal),
    Identifier { name: Token },
    Unary { op: Token, expr: Box<Expr> },
}

impl Hash for Expr {
    fn hash<H: std::hash::Hasher>(&self, state: &mut H) {
        core::mem::discriminant(self).hash(state);
    }
}

#[derive(Debug)]
pub enum Literal {
    Int(i64),
    Dec(f64),
    Str(String),
    Bool(bool),
    List(Vec<Expr>),
    Map(HashMap<Expr, Expr>),
}

impl PartialEq for Literal {
    fn eq(&self, other: &Self) -> bool {
        match (self, other) {
            (Self::Int(l0), Self::Int(r0)) => l0 == r0,
            (Self::Dec(l0), Self::Dec(r0)) => l0 == r0,
            (Self::Str(l0), Self::Str(r0)) => l0 == r0,
            (Self::Bool(l0), Self::Bool(r0)) => l0 == r0,
            (Self::List(l0), Self::List(r0)) => l0 == r0,
            (Self::Map(l0), Self::Map(r0)) => {
                l0.len() == r0.len()
                    && l0
                        .iter()
                        .zip(r0.iter())
                        .all(|((lk, lv), (rk, rv))| lk == rk && lv == rv)
            }
            _ => false,
        }
    }
}

impl Eq for Literal {}

impl Hash for Literal {
    fn hash<H: std::hash::Hasher>(&self, state: &mut H) {
        core::mem::discriminant(self).hash(state);
    }
}
