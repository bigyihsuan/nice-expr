use crate::eval::{env::SEnv, value::Value, RuntimeError};

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
        left.chars().filter(|c| right.contains(*c)).collect(),
    ))
}

pub fn bnot(left: Value, _env: &SEnv) -> Result<Value, RuntimeError> {
    let b: bool = left.try_into()?;
    Ok(Value::Bool(!b))
}
