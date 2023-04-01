use crate::parse::ast::{AssignmentOperator, Type};

pub mod env;
pub mod intepreter;
pub mod value;

#[derive(Debug, Clone)]
pub enum Constness {
    Const,
    Var,
}

#[derive(Debug, Clone)]
pub enum RuntimeError {
    NotImplemented,
    MismatchedTypes { got: Type, expected: Type },
    IdentifierNotFound(String),
    SettingConst(String),
    DivideByZero,
    InvalidAssignmentOperator(AssignmentOperator),
    NotEnoughArguments { want: usize, got: usize },
    IndexingNonIndexable { got: Type },
    TakingLenOfLengthless { got: crate::parse::ast::Type },
    // TODO: more runtime errors
}
