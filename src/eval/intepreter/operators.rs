use std::collections::HashMap;

use itertools::Itertools;

use crate::{
    eval::{env::SEnv, value::Value},
    prelude::RuntimeError,
};

pub fn ineg(left: Value, _env: &SEnv) -> Result<Value, RuntimeError> {
    let i: i64 = left.try_into()?;
    Ok(Value::Int(-i))
}
pub fn iadd(left: Value, right: Value, _env: &SEnv) -> Result<Value, RuntimeError> {
    let left: i64 = left.try_into()?;
    let right: i64 = right.try_into()?;
    Ok(Value::Int(left + right))
}
pub fn isub(left: Value, right: Value, _env: &SEnv) -> Result<Value, RuntimeError> {
    let left: i64 = left.try_into()?;
    let right: i64 = right.try_into()?;
    Ok(Value::Int(left - right))
}
pub fn imul(left: Value, right: Value, _env: &SEnv) -> Result<Value, RuntimeError> {
    let left: i64 = left.try_into()?;
    let right: i64 = right.try_into()?;
    Ok(Value::Int(left * right))
}
pub fn idiv(left: Value, right: Value, _env: &SEnv) -> Result<Value, RuntimeError> {
    let left: i64 = left.try_into()?;
    let right: i64 = right.try_into()?;
    if right == 0 {
        Err(RuntimeError::DivideByZero)
    } else {
        Ok(Value::Int(left / right))
    }
}
pub fn imod(left: Value, right: Value, _env: &SEnv) -> Result<Value, RuntimeError> {
    let left: i64 = left.try_into()?;
    let right: i64 = right.try_into()?;
    if right == 0 {
        Err(RuntimeError::DivideByZero)
    } else {
        Ok(Value::Int(left % right))
    }
}

pub fn fneg(left: Value, _env: &SEnv) -> Result<Value, RuntimeError> {
    let f: f64 = left.try_into()?;
    Ok(Value::Dec(-f))
}
pub fn fadd(left: Value, right: Value, _env: &SEnv) -> Result<Value, RuntimeError> {
    let left: f64 = left.try_into()?;
    let right: f64 = right.try_into()?;
    Ok(Value::Dec(left + right))
}
pub fn fsub(left: Value, right: Value, _env: &SEnv) -> Result<Value, RuntimeError> {
    let left: f64 = left.try_into()?;
    let right: f64 = right.try_into()?;
    Ok(Value::Dec(left - right))
}
pub fn fmul(left: Value, right: Value, _env: &SEnv) -> Result<Value, RuntimeError> {
    let left: f64 = left.try_into()?;
    let right: f64 = right.try_into()?;
    Ok(Value::Dec(left * right))
}
pub fn fdiv(left: Value, right: Value, _env: &SEnv) -> Result<Value, RuntimeError> {
    let left: f64 = left.try_into()?;
    let right: f64 = right.try_into()?;
    if right == 0.0 {
        Err(RuntimeError::DivideByZero)
    } else {
        Ok(Value::Dec(left / right))
    }
}

pub fn sadd(left: Value, right: Value, _env: &SEnv) -> Result<Value, RuntimeError> {
    let left: String = left.try_into()?;
    let right: String = right.try_into()?;
    Ok(Value::Str(format!("{left}{right}")))
}
pub fn ssub(left: Value, right: Value, _env: &SEnv) -> Result<Value, RuntimeError> {
    let left: String = left.try_into()?;
    let right: String = right.try_into()?;
    Ok(Value::Str(
        left.chars().filter(|c| !right.contains(*c)).collect(),
    ))
}
pub fn sidx(left: Value, right: Value, _env: &SEnv) -> Result<Value, RuntimeError> {
    let l: String = left.clone().try_into()?;
    let r: isize = right.clone().try_into()?;

    let len = l.chars().count();

    if len == 0 {
        Err(RuntimeError::IndexingCollectionWithZeroElements(left))
    } else if r >= 0 {
        match l.chars().skip(r as usize).next() {
            Some(c) => Ok(Value::Str(c.to_string())),
            None => Err(RuntimeError::IndexOutOfBounds(left, right)),
        }
    } else {
        match l.chars().rev().skip((-r) as usize - 1).next() {
            Some(c) => Ok(Value::Str(c.to_string())),
            None => Err(RuntimeError::IndexOutOfBounds(left, right)),
        }
    }
}

pub fn bnot(left: Value, _env: &SEnv) -> Result<Value, RuntimeError> {
    let b: bool = left.try_into()?;
    Ok(Value::Bool(!b))
}
pub fn band(left: Value, right: Value, _env: &SEnv) -> Result<Value, RuntimeError> {
    let left: bool = left.try_into()?;
    let right: bool = right.try_into()?;

    Ok(Value::Bool(left && right))
}
pub fn bor(left: Value, right: Value, _env: &SEnv) -> Result<Value, RuntimeError> {
    let left: bool = left.try_into()?;
    let right: bool = right.try_into()?;

    Ok(Value::Bool(left || right))
}

pub fn ladd(left: Value, right: Value, _env: &SEnv) -> Result<Value, RuntimeError> {
    let left: Vec<Value> = left.try_into()?;
    let right: Vec<Value> = right.try_into()?;

    let l = left
        .into_iter()
        .chain(right.into_iter())
        .collect::<Vec<_>>();
    let result = Value::List(l.clone());
    if result.is_homogeneous() {
        Ok(result)
    } else {
        Err(RuntimeError::MismatchedTypes {
            got: l
                .into_iter()
                .map(|e| e.to_type())
                .filter_map(|e| e.ok())
                .unique()
                .collect(),
            expected: vec![result.to_type()?],
        })
    }
}
pub fn lsub(left: Value, right: Value, _env: &SEnv) -> Result<Value, RuntimeError> {
    let left: Vec<Value> = left.try_into()?;
    let right: Vec<Value> = right.try_into()?;
    Ok(Value::List(
        left.into_iter().filter(|c| !right.contains(c)).collect(),
    ))
}
pub fn lidx(left: Value, right: Value, _env: &SEnv) -> Result<Value, RuntimeError> {
    let l: Vec<Value> = left.clone().try_into()?;
    let r: isize = right.clone().try_into()?;

    let len = l.len();

    if len == 0 {
        Err(RuntimeError::IndexingCollectionWithZeroElements(left))
    } else if r >= 0 {
        match l.into_iter().skip(r as usize).next() {
            Some(e) => Ok(e),
            None => Err(RuntimeError::IndexOutOfBounds(left, right)),
        }
    } else {
        match l.into_iter().rev().skip((-r) as usize - 1).next() {
            Some(e) => Ok(e),
            None => Err(RuntimeError::IndexOutOfBounds(left, right)),
        }
    }
}

pub fn madd(left: Value, right: Value, _env: &SEnv) -> Result<Value, RuntimeError> {
    let left: HashMap<Value, Value> = left.try_into()?;
    let right: HashMap<Value, Value> = right.try_into()?;

    let l = left
        .into_iter()
        .chain(right.into_iter())
        .collect::<HashMap<_, _>>();
    let result = Value::Map(l.clone());
    if result.is_homogeneous() {
        Ok(result)
    } else {
        Err(RuntimeError::MismatchedTypes {
            got: l
                .into_iter()
                .map(|(k, v)| (k.to_type(), v.to_type()))
                .filter_map(|(k, v)| k.ok().and(v.ok()))
                .unique()
                .collect(),
            expected: vec![result.to_type()?],
        })
    }
}
pub fn msub(left: Value, right: Value, _env: &SEnv) -> Result<Value, RuntimeError> {
    let left: HashMap<Value, Value> = left.try_into()?;
    let right: HashMap<Value, Value> = right.try_into()?;
    Ok(Value::Map(
        left.into_iter()
            .filter(|(k, _)| right.get(k).is_none())
            .collect(),
    ))
}
pub fn midx(left: Value, right: Value, _env: &SEnv) -> Result<Value, RuntimeError> {
    let l: HashMap<Value, Value> = left.clone().try_into()?;
    let r: Value = right.clone();

    if l.len() == 0 {
        Err(RuntimeError::IndexingCollectionWithZeroElements(left))
    } else {
        match l.get(&r) {
            Some(v) => Ok(v.clone()),
            None => Err(RuntimeError::KeyNotFound(left, right)),
        }
    }
}

pub fn eq(left: Value, right: Value, _env: &SEnv) -> Result<Value, RuntimeError> {
    Ok(Value::Bool(left == right))
}
pub fn ne(left: Value, right: Value, _env: &SEnv) -> Result<Value, RuntimeError> {
    Ok(Value::Bool(left != right))
}
pub fn gt(left: Value, right: Value, _env: &SEnv) -> Result<Value, RuntimeError> {
    Ok(Value::Bool(left > right))
}
pub fn ge(left: Value, right: Value, _env: &SEnv) -> Result<Value, RuntimeError> {
    Ok(Value::Bool(left >= right))
}
pub fn lt(left: Value, right: Value, _env: &SEnv) -> Result<Value, RuntimeError> {
    Ok(Value::Bool(left < right))
}
pub fn le(left: Value, right: Value, _env: &SEnv) -> Result<Value, RuntimeError> {
    Ok(Value::Bool(left <= right))
}
