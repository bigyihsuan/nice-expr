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
    SettingConst,
    IdentifierNotFound,
    DivideByZero,
    IndexingNonIndexable,
    MismatchedTypes,
    // TODO: more runtime errors
}
