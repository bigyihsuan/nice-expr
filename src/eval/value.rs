use std::{cell::RefCell, cmp::Ordering, collections::HashMap, hash::Hash, rc::Rc};

use crate::{parse::ast::Expr, prelude::RuntimeError};

use super::{env::SEnv, r#type::Type};

#[derive(Debug, Clone)]
pub enum Value {
    None,
    Int(i64),
    Dec(f64),
    Str(String),
    Bool(bool),
    List(Vec<Value>),
    Map(HashMap<Value, Value>),
    Func(Function),
}

impl Value {
    pub fn to_type(&self) -> Result<Type, RuntimeError> {
        match self {
            Value::None => Ok(Type::None),
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
            Value::Func(_) => todo!(),
        }
    }

    pub fn is_homogeneous(&self) -> bool {
        match self {
            Value::None => true,
            Value::Int(_) => true,
            Value::Dec(_) => true,
            Value::Str(_) => true,
            Value::Bool(_) => true,
            Value::List(l) => l
                .into_iter()
                .map(|e| e.is_homogeneous())
                .reduce(|acc, e| acc && e)
                .unwrap_or(false),
            Value::Map(m) => m
                .into_iter()
                .map(|e| e.0.is_homogeneous() && e.1.is_homogeneous())
                .reduce(|acc, e| acc && e)
                .unwrap_or(false),
            Value::Func(_) => todo!(),
        }
    }

    pub fn default(t: Type) -> Value {
        match t {
            Type::None => Value::None,
            Type::Int => Value::Int(0),
            Type::Dec => Value::Dec(0.0),
            Type::Str => Value::Str(String::new()),
            Type::Bool => Value::Bool(false),
            Type::List(_) => Value::List(Vec::new()),
            Type::Map(_, _) => Value::Map(HashMap::new()),
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

impl PartialOrd for Value {
    fn partial_cmp(&self, other: &Self) -> Option<std::cmp::Ordering> {
        match (self, other) {
            (Self::None, Self::None) => Some(Ordering::Equal),
            (Self::Int(l0), Self::Int(r0)) => l0.partial_cmp(r0),
            (Self::Dec(l0), Self::Dec(r0)) => l0.partial_cmp(r0),
            (Self::Str(l0), Self::Str(r0)) => l0.partial_cmp(r0),
            (Self::Bool(l0), Self::Bool(r0)) => l0.partial_cmp(r0),
            _ => None,
        }
    }
}

impl Hash for Value {
    fn hash<H: std::hash::Hasher>(&self, state: &mut H) {
        core::mem::discriminant(self).hash(state);
    }
}

impl TryFrom<Value> for i64 {
    type Error = RuntimeError;

    fn try_from(value: Value) -> Result<Self, Self::Error> {
        let t = value.to_type()?;
        match value {
            Value::Int(i) => Ok(i),
            _ => Err(RuntimeError::MismatchedTypes {
                got: vec![t],
                expected: vec![Type::Bool],
            }),
        }
    }
}
impl TryFrom<&Value> for i64 {
    type Error = RuntimeError;

    fn try_from(value: &Value) -> Result<Self, Self::Error> {
        let t = value.to_type()?;
        match value {
            Value::Int(i) => Ok(*i),
            _ => Err(RuntimeError::MismatchedTypes {
                got: vec![t],
                expected: vec![Type::Bool],
            }),
        }
    }
}

impl TryFrom<Value> for usize {
    type Error = RuntimeError;

    fn try_from(value: Value) -> Result<Self, Self::Error> {
        let t = value.to_type()?;
        match value {
            Value::Int(i) => Ok(i as Self),
            _ => Err(RuntimeError::MismatchedTypes {
                got: vec![t],
                expected: vec![Type::Int],
            }),
        }
    }
}
impl TryFrom<&Value> for usize {
    type Error = RuntimeError;

    fn try_from(value: &Value) -> Result<Self, Self::Error> {
        let t = value.to_type()?;
        match value {
            Value::Int(i) => Ok(*i as Self),
            _ => Err(RuntimeError::MismatchedTypes {
                got: vec![t],
                expected: vec![Type::Int],
            }),
        }
    }
}

impl TryFrom<Value> for isize {
    type Error = RuntimeError;

    fn try_from(value: Value) -> Result<Self, Self::Error> {
        let t = value.to_type()?;
        match value {
            Value::Int(i) => Ok(i as Self),
            _ => Err(RuntimeError::MismatchedTypes {
                got: vec![t],
                expected: vec![Type::Int],
            }),
        }
    }
}
impl TryFrom<&Value> for isize {
    type Error = RuntimeError;

    fn try_from(value: &Value) -> Result<Self, Self::Error> {
        let t = value.to_type()?;
        match value {
            Value::Int(i) => Ok(*i as Self),
            _ => Err(RuntimeError::MismatchedTypes {
                got: vec![t],
                expected: vec![Type::Int],
            }),
        }
    }
}

impl TryFrom<Value> for f64 {
    type Error = RuntimeError;

    fn try_from(value: Value) -> Result<Self, Self::Error> {
        let t = value.to_type()?;
        match value {
            Value::Dec(f) => Ok(f),
            _ => Err(RuntimeError::MismatchedTypes {
                got: vec![t],
                expected: vec![Type::Bool],
            }),
        }
    }
}
impl TryFrom<&Value> for f64 {
    type Error = RuntimeError;

    fn try_from(value: &Value) -> Result<Self, Self::Error> {
        let t = value.to_type()?;
        match value {
            Value::Dec(f) => Ok(*f),
            _ => Err(RuntimeError::MismatchedTypes {
                got: vec![t],
                expected: vec![Type::Bool],
            }),
        }
    }
}

impl TryFrom<Value> for bool {
    type Error = RuntimeError;

    fn try_from(value: Value) -> Result<Self, Self::Error> {
        let t = value.to_type()?;
        match value {
            Value::Bool(b) => Ok(b),
            _ => Err(RuntimeError::MismatchedTypes {
                got: vec![t],
                expected: vec![Type::Bool],
            }),
        }
    }
}
impl TryFrom<&Value> for bool {
    type Error = RuntimeError;

    fn try_from(value: &Value) -> Result<Self, Self::Error> {
        let t = value.to_type()?;
        match value {
            Value::Bool(b) => Ok(*b),
            _ => Err(RuntimeError::MismatchedTypes {
                got: vec![t],
                expected: vec![Type::Bool],
            }),
        }
    }
}

impl TryFrom<Value> for String {
    type Error = RuntimeError;

    fn try_from(value: Value) -> Result<Self, Self::Error> {
        let t = value.to_type()?;
        match value {
            Value::Str(s) => Ok(s),
            _ => Err(RuntimeError::MismatchedTypes {
                got: vec![t],
                expected: vec![Type::Bool],
            }),
        }
    }
}
impl TryFrom<&Value> for String {
    type Error = RuntimeError;

    fn try_from(value: &Value) -> Result<Self, Self::Error> {
        let t = value.to_type()?;
        match value {
            Value::Str(s) => Ok(s.clone()),
            _ => Err(RuntimeError::MismatchedTypes {
                got: vec![t],
                expected: vec![Type::Bool],
            }),
        }
    }
}

impl TryFrom<Value> for Vec<Value> {
    type Error = RuntimeError;

    fn try_from(value: Value) -> Result<Self, Self::Error> {
        let t = value.to_type()?;
        match value {
            Value::List(v) => Ok(v),
            _ => Err(RuntimeError::MismatchedTypes {
                got: vec![t],
                expected: vec![Type::Bool],
            }),
        }
    }
}
impl TryFrom<&Value> for Vec<Value> {
    type Error = RuntimeError;

    fn try_from(value: &Value) -> Result<Self, Self::Error> {
        let t = value.to_type()?;
        match value {
            Value::List(v) => Ok(v.clone()),
            _ => Err(RuntimeError::MismatchedTypes {
                got: vec![t],
                expected: vec![Type::Bool],
            }),
        }
    }
}

impl TryFrom<Value> for HashMap<Value, Value> {
    type Error = RuntimeError;

    fn try_from(value: Value) -> Result<Self, Self::Error> {
        let t = value.to_type()?;
        match value {
            Value::Map(m) => Ok(m),
            _ => Err(RuntimeError::MismatchedTypes {
                got: vec![t],
                expected: vec![Type::Bool],
            }),
        }
    }
}
impl TryFrom<&Value> for HashMap<Value, Value> {
    type Error = RuntimeError;

    fn try_from(value: &Value) -> Result<Self, Self::Error> {
        let t = value.to_type()?;
        match value {
            Value::Map(m) => Ok(m.clone()),
            _ => Err(RuntimeError::MismatchedTypes {
                got: vec![t],
                expected: vec![Type::Bool],
            }),
        }
    }
}

impl TryFrom<Value> for Function {
    type Error = RuntimeError;

    fn try_from(value: Value) -> Result<Self, Self::Error> {
        let t = value.to_type()?;
        match value {
            Value::Func(f) => Ok(f),
            _ => Err(RuntimeError::MismatchedTypes {
                got: vec![t],
                expected: vec![Type::Bool],
            }),
        }
    }
}
impl TryFrom<&Value> for Function {
    type Error = RuntimeError;

    fn try_from(value: &Value) -> Result<Self, Self::Error> {
        let t = value.to_type()?;
        match value {
            Value::Func(f) => Ok(f.clone()),
            _ => Err(RuntimeError::MismatchedTypes {
                got: vec![t],
                expected: vec![Type::Bool],
            }),
        }
    }
}

#[derive(Debug, Clone)]
pub struct Function {
    env: SEnv,
    args: Vec<Value>,
    block: Expr,
}
