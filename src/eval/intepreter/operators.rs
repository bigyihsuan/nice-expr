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
pub fn sgetidx(
    left: Value,
    start_index: Value,
    end_index: Option<Value>,
    _env: &SEnv,
) -> Result<Value, RuntimeError> {
    let l: String = left.clone().try_into()?;
    let len = l.chars().count();

    let (range_start, _, take_count) = transform_to_range(len, start_index, end_index)?;

    if len == 0 {
        Err(RuntimeError::IndexingCollectionWithZeroElements(left))
    } else {
        Ok(Value::Str(
            l.chars().skip(range_start).take(take_count).collect(),
        ))
    }
}
pub fn ssetidx(
    collection: Value,
    index: Value,
    element: Value,
    _env: &SEnv,
) -> Result<Value, RuntimeError> {
    let mut c: String = collection.clone().try_into()?;
    let i: isize = index.clone().try_into()?;
    let e: String = element.clone().try_into()?;
    if i >= 0 {
        let i = i as usize;
        if i > c.len() {
            return Err(RuntimeError::IndexOutOfBounds(collection, index, None));
        }
        c.replace_range(i..=i, e.as_str());
        Ok(Value::Str(c))
    } else {
        if -i as usize > c.len() + 1 {
            return Err(RuntimeError::IndexOutOfBounds(collection, index, None));
        }
        c.insert_str(c.len() - (-i as usize) + 1, e.as_str());
        Ok(Value::Str(c))
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
pub fn lgetidx(
    left: Value,
    start_index: Value,
    end_index: Option<Value>,
    _env: &SEnv,
) -> Result<Value, RuntimeError> {
    let l: Vec<Value> = left.clone().try_into()?;
    let len = l.len();

    let (range_start, _, take_count) =
        transform_to_range(len, start_index.clone(), end_index.clone())?;

    if len == 0 {
        Err(RuntimeError::IndexingCollectionWithZeroElements(left))
    } else if let None = end_index && take_count == 1 {
        let ele =
            l.get(range_start)
                .ok_or(RuntimeError::IndexOutOfBounds(left, start_index, None))?;
        Ok(ele.clone())
    } else {
        Ok(Value::List(
            l.into_iter().skip(range_start).take(take_count).collect(),
        ))
    }
}
pub fn lsetidx(
    collection: Value,
    index: Value,
    element: Value,
    _env: &SEnv,
) -> Result<Value, RuntimeError> {
    let mut c: Vec<Value> = collection.clone().try_into()?;
    let i: isize = index.clone().try_into()?;
    let e: Value = element.clone();

    let len = c.len();

    if len == 0 {
        Err(RuntimeError::IndexingCollectionWithZeroElements(collection))
    } else if i >= 0 {
        let i = i as usize;
        if i > c.len() {
            return Err(RuntimeError::IndexOutOfBounds(collection, index, None));
        }
        c[i] = e;
        Ok(Value::List(c))
    } else {
        if -i as usize > c.len() + 1 {
            return Err(RuntimeError::IndexOutOfBounds(collection, index, None));
        }
        let l = c.len();
        c[l - (-i as usize) + 1] = e;
        Ok(Value::List(c))
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
pub fn mgetidx(
    left: Value,
    key: Value,
    _: Option<Value>,
    _env: &SEnv,
) -> Result<Value, RuntimeError> {
    let l: HashMap<Value, Value> = left.clone().try_into()?;
    let r: Value = key.clone();

    if l.len() == 0 {
        Err(RuntimeError::IndexingCollectionWithZeroElements(left))
    } else {
        match l.get(&r) {
            Some(v) => Ok(v.clone()),
            None => Err(RuntimeError::KeyNotFound(left, key)),
        }
    }
}
pub fn msetidx(
    collection: Value,
    key: Value,
    value: Value,
    _env: &SEnv,
) -> Result<Value, RuntimeError> {
    let mut m: HashMap<Value, Value> = collection.clone().try_into()?;
    let k: Value = key.clone();

    if m.len() == 0 {
        Err(RuntimeError::IndexingCollectionWithZeroElements(collection))
    } else {
        m.insert(k, value);
        Ok(Value::Map(m))
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

fn transform_to_range(
    len: usize,
    start: Value,
    end: Option<Value>,
) -> Result<(usize, usize, usize), RuntimeError> {
    let start: isize = start.clone().try_into()?;
    let end: isize = if let Some(end) = &end {
        end.clone().try_into()?
    } else {
        start + 1
    };

    let range_start = if start < 0 {
        len - (-start) as usize
    } else {
        start as usize
    };
    let range_end = if end <= 0 {
        len - (-end) as usize
    } else {
        end as usize
    };

    let take_count = if range_end >= range_start {
        range_end - range_start
    } else {
        0
    };

    // println!("[{start},{end}) -> [{range_start},{range_end}) take {take_count}");
    Ok((range_start, range_end, take_count))
}
