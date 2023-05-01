use std::{cmp::Ordering, collections::HashMap, hash::Hash};

use itertools::Itertools;

use crate::{
    parse::ast::{Decl, Declaration, Expr},
    prelude::RuntimeError,
};

use super::{env::SEnv, r#type::Type};

#[derive(Debug, Clone)]
pub enum Value {
    None,
    Break(Box<Value>),
    Int(i64),
    Dec(f64),
    Str(String),
    Bool(bool),
    List(Vec<Value>),
    Map(HashMap<Value, Value>),
    Func(Func),
    Type(Type),
}

#[derive(Debug, Clone)]
pub enum Func {
    Native(NativeFunc),
    Declared {
        decls: Vec<Declaration>,
        ret: Type,
        body: Vec<Expr>,
        env: SEnv,
    },
}

pub type NativeFunc = (String, fn(args: &[Value]) -> Result<Value, RuntimeError>);

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
            Value::Break(v) => {
                let v = v.to_type()?;
                Ok(Type::Break(Box::new(v)))
            }
            Value::Func(Func::Declared {
                decls: args,
                ret,
                body: _,
                env: _,
            }) => {
                let args = args
                    .into_iter()
                    .filter_map(|decl| match decl {
                        Declaration::Const(Decl {
                            name: _,
                            type_name,
                            expr: _,
                        }) => Some(type_name.clone()),
                        Declaration::Var(Decl {
                            name: _,
                            type_name,
                            expr: _,
                        }) => Some(type_name.clone()),
                    })
                    .collect_vec();
                Ok(Type::Func(args, Box::new(ret.clone())))
            }
            Value::Func(Func::Native(_)) => Ok(Type::Func(
                vec![Type::BuiltinVariadic],
                Box::new(Type::None),
            )),
            Value::Type(t) => Ok(t.clone()),
        }
    }

    // returns if all elements in value are of the same type
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
            Value::Break(v) => v.is_homogeneous(),
            Value::Func(Func::Declared { .. }) => true,
            Value::Func(Func::Native(_)) => todo!(),
            Value::Type(_) => true, // are types homogeneous?
        }
    }

    pub fn unbreak(&self) -> Self {
        if let Value::Break(box v) = self {
            v.clone()
        } else {
            self.clone()
        }
    }

    pub fn default(t: Type) -> Self {
        match t {
            Type::None => Value::None,
            Type::BuiltinVariadic => Value::None,
            Type::Int => Value::Int(0),
            Type::Dec => Value::Dec(0.0),
            Type::Str => Value::Str(String::new()),
            Type::Bool => Value::Bool(false),
            Type::List(_) => Value::List(Vec::new()),
            Type::Map(_, _) => Value::Map(HashMap::new()),
            Type::Break(box t) => Value::Break(Box::new(Self::default(t))),
            Type::Func(_, _) => todo!("what to do with default func value"), // TODO: what to do with default func value
            Type::Any => todo!("what to do with any"),
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
            (Self::Break(box l0), Self::Break(box r0)) => l0 == r0,
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
            (Self::Break(box l0), Self::Break(box r0)) => l0.partial_cmp(r0),
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
        match value.unbreak() {
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
        match value.unbreak() {
            Value::Int(i) => Ok(i),
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
        match value.unbreak() {
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
        match value.unbreak() {
            Value::Int(i) => Ok(i as Self),
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
        match value.unbreak() {
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
        match value.unbreak() {
            Value::Int(i) => Ok(i as Self),
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
        match value.unbreak() {
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
        match value.unbreak() {
            Value::Dec(f) => Ok(f),
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
        match value.unbreak() {
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
        match value.unbreak() {
            Value::Bool(b) => Ok(b),
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
        match value.unbreak() {
            Value::Str(s) => Ok(s.clone()),
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
        match value.unbreak() {
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
        match value.unbreak() {
            Value::List(v) => Ok(v.clone()),
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
        match value.unbreak() {
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
        match value.unbreak() {
            Value::Map(m) => Ok(m.clone()),
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
        match value.unbreak() {
            Value::Map(m) => Ok(m.clone()),
            _ => Err(RuntimeError::MismatchedTypes {
                got: vec![t],
                expected: vec![Type::Bool],
            }),
        }
    }
}
