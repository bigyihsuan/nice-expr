use std::{collections::HashMap, hash::Hash};

#[derive(Debug, Clone)]
pub enum Value {
    Int(i64),
    Dec(f64),
    Str(String),
    Bool(bool),
    List(Vec<Value>),
    Map(HashMap<Value, Value>),
    // TODO: Ident(Identifier),
    // TODO: Func(Function),
}

impl PartialEq for Value {
    fn eq(&self, other: &Self) -> bool {
        match (self, other) {
            (Self::Int(l0), Self::Int(r0)) => l0 == r0,
            (Self::Dec(l0), Self::Dec(r0)) => l0 == r0,
            (Self::Str(l0), Self::Str(r0)) => l0 == r0,
            (Self::Bool(l0), Self::Bool(r0)) => l0 == r0,
            (Self::List(l0), Self::List(r0)) => l0 == r0,
            (Self::Map(l0), Self::Map(r0)) => {
                l0.len() == r0.len() && l0.iter().all(|(lk, lv)| r0.get(lk) == Some(lv))
            }
            _ => false,
        }
    }
}
impl Eq for Value {}

impl Hash for Value {
    fn hash<H: std::hash::Hasher>(&self, state: &mut H) {
        core::mem::discriminant(self).hash(state);
    }
}
