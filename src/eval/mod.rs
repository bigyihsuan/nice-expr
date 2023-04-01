use crate::parse::ast::{AssignmentOperator, Type};

pub mod env;
pub mod intepreter;
pub mod value;

#[derive(Debug, Clone)]
pub enum Constness {
    Const,
    Var,
}

#[derive(Debug)]
pub enum RuntimeError {
    NotImplemented,
    MismatchedTypes { got: Type, expected: Type },
    // assignments
    IdentifierNotFound(String),
    SettingConst(String),
    // operators
    DivideByZero,
    IndexingNonIndexable,
    InvalidAssignmentOperator(AssignmentOperator),
    // TODO: more runtime errors
}
