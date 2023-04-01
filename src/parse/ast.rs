use crate::lexer::tok::Token;

pub type Program = Vec<Expr>;

#[derive(Debug)]
pub enum Expr {
    Literal(Literal),
    Identifier(String),
    Unary { op: Token, expr: Box<Expr> },
    Declaration(Declaration),
    Assignment(Assignment),
    FunctionCall { name: String, args: Vec<Expr> },
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

impl Type {
    pub fn are_containers(&self, other: &Self) -> bool {
        match (self, other) {
            (Type::List(_), Type::List(_)) => true,
            (Type::Map(_, _), Type::Map(_, _)) => true,
            _ => false,
        }
    }

    pub fn infer_contained_type(&self, other: &Self) -> Option<Self> {
        match (self, other) {
            (Type::List(l), Type::List(r))
                if *l != Box::new(Type::None) && *r == Box::new(Type::None) =>
            {
                Some(Type::List(l.clone()))
            }

            (Type::List(l), Type::List(r))
                if *l == Box::new(Type::None) && *r != Box::new(Type::None) =>
            {
                Some(Type::List(r.clone()))
            }
            (Type::Map(lk, lv), Type::Map(rk, rv)) => {
                match (
                    *lk != Box::new(Type::None),
                    *lv != Box::new(Type::None),
                    *rk != Box::new(Type::None),
                    *rv != Box::new(Type::None),
                ) {
                    (true, true, true, true) => Some(self.clone()),
                    (true, true, true, false) => Some(self.clone()),
                    (true, true, false, true) => Some(self.clone()),
                    (true, true, false, false) => Some(self.clone()),
                    (true, false, true, true) => Some(other.clone()),
                    (false, true, true, true) => Some(other.clone()),
                    (false, false, true, true) => Some(other.clone()),
                    (true, false, false, true) => Some(Type::Map(lk.clone(), rv.clone())),
                    (false, true, true, false) => Some(Type::Map(rk.clone(), lv.clone())),
                    _ => None,
                }
            }
            (l, _) if *l != Type::None => Some(l.clone()),
            (_, r) if *r != Type::None => Some(r.clone()),
            _ => None,
        }
    }
}

#[derive(Debug, Clone)]
pub enum AssignmentOperator {
    Invalid,
    Is,
    // TODO: the other assignment operators
}
