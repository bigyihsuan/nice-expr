use std::{io, string};

use crate::parse::ast::AssignmentOperator;

use self::r#type::Type;

pub mod env;
pub mod intepreter;
pub mod r#type;
pub mod value;

#[derive(Debug, Clone)]
pub enum Constness {
    Const,
    Var,
}
