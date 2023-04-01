use crate::parse::ast::Type;

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
    MismatchedTypes { got: Type, expected: Type },
    // assignments
    IdentifierNotFound(String),
    SettingConst(String),
    // operators
    DivideByZero,
    IndexingNonIndexable,
    // TODO: more runtime errors
}
