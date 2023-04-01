use std::{cell::RefCell, collections::HashMap, hash::Hash, rc::Rc};

use crate::parse::ast::Type;

use super::{env::Env, RuntimeError};

#[derive(Debug, Clone)]
pub enum Value {
    Int(i64),
    Dec(f64),
    Str(String),
    Bool(bool),
    List(Vec<Value>),
    Map(HashMap<Value, Value>),
    // TODO: Func(Function),
}

impl Value {
    pub fn to_type(&self) -> Result<Type, RuntimeError> {
        match self {
            Value::Int(_) => Ok(Type::Int),
            Value::Dec(_) => Ok(Type::Dec),
            Value::Str(_) => Ok(Type::Str),
            Value::Bool(_) => Ok(Type::Bool),
            Value::List(l) => {
                let t = l.get(0).map_or(Ok(Type::None), |e| e.to_type())?;
                Ok(Type::List(Box::new(t)))
            }
            Value::Map(m) => {
                let e = m
                    .iter()
                    .take(1)
                    .unzip::<&Value, &Value, Vec<&Value>, Vec<&Value>>();
                let k =
                    e.0.get(0)
                        .map(|k| k.to_type())
                        .unwrap_or_else(|| Ok(Type::None))?;
                let v =
                    e.1.get(0)
                        .map(|v| v.to_type())
                        .unwrap_or_else(|| Ok(Type::None))?;
                Ok(Type::Map(Box::new(k), Box::new(v)))
            }
        }
    }
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
