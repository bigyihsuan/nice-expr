



pub mod env;
pub mod intepreter;
pub mod r#type;
pub mod value;

#[derive(Debug, Clone)]
pub enum Constness {
    Const,
    Var,
}
