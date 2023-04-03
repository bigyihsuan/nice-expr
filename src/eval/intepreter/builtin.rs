use std::io::{self, Read};

use itertools::Itertools;
use unicode_reader::CodePoints;

use crate::{
    eval::value::{Func, NativeFunc, Value},
    parse::ast::{Declaration, Expr},
    prelude::{IOError, RuntimeError},
    util::{self, assert_at_least_args},
};

pub fn format_value(value: &Value) -> String {
    match value {
        Value::None => "None".into(),
        Value::Int(i) => format!("{i}"),
        Value::Dec(d) => format!("{d}"),
        Value::Str(s) => format!("{s}"),
        Value::Bool(b) => format!("{b}"),
        Value::List(l) => {
            let l = l
                .into_iter()
                .map(|e| format_value(&e))
                .collect_vec()
                .join(",");
            format!("[{l}]")
        }
        Value::Map(m) => {
            let m = m
                .into_iter()
                .map(|(k, v)| (format_value(&k), format_value(&v)))
                .map(|(k, v)| format!("{k}:{v}"))
                .collect_vec()
                .join(",");
            format!("<|{m}|>")
        }
        Value::Func(Func::Declared {
            decls: args,
            ret,
            body: _,
        }) => {
            let args = args
                .into_iter()
                .filter_map(|decl| {
                    if let Expr::Declaration(decl) = decl {
                        match decl {
                            Declaration::Const {
                                name,
                                type_name,
                                expr: _,
                            } => Some((name, format!("{type_name}"))),
                            Declaration::Var {
                                name,
                                type_name,
                                expr: _,
                            } => Some((name, format!("{type_name}"))),
                        }
                    } else {
                        None
                    }
                })
                .map(|(n, t)| format!("{n} {t}"))
                .collect_vec()
                .join(",");
            let ret = format!("{ret}");

            format!("func({args}){ret}{{...}}")
        }
        Value::Func(Func::Native(_)) => format!("builtin_function"),
        Value::Break(box v) => format!("{}", format_value(v)),
    }
}

pub fn print(args: &[Value]) -> Result<Value, RuntimeError> {
    for arg in args {
        print!("{}", format_value(&arg));
    }
    Ok(Value::None)
}
pub fn println(args: &[Value]) -> Result<Value, RuntimeError> {
    for arg in args {
        println!("{}", format_value(&arg));
    }
    Ok(Value::None)
}

pub fn len(args: &[Value]) -> Result<Value, RuntimeError> {
    assert_at_least_args(1, args.len())?;
    let val = &args[0];
    match val {
        Value::Str(v) => Ok(Value::Int(v.chars().collect_vec().len() as i64)),
        Value::List(v) => Ok(Value::Int(v.len() as i64)),
        Value::Map(v) => Ok(Value::Int(v.len() as i64)),
        _ => Err(RuntimeError::TakingLenOfLengthless {
            got: val.to_type()?,
        }),
    }
}

pub fn range(args: &[Value]) -> Result<Value, RuntimeError> {
    assert_at_least_args(3, args.len())?;
    let start = &args[0];
    let end = &args[1];
    let step = &args[2];

    let start: isize = start.try_into()?;
    let end: isize = end.try_into()?;
    let step: isize = step.try_into()?;

    if step > 0 {
        Ok(Value::List(
            (start..end)
                .step_by(step as usize)
                .map(|i| Value::Int(i as i64))
                .collect_vec(),
        ))
    } else if step < 0 {
        Ok(Value::List(
            (end..start)
                .step_by(step as usize)
                .rev()
                .map(|i| Value::Int(i as i64))
                .collect_vec(),
        ))
    } else {
        Err(RuntimeError::InvalidRangeStep(step))
    }
}

pub fn inputchar(_: &[Value]) -> Result<Value, RuntimeError> {
    let c = CodePoints::from(io::stdin().bytes())
        .map(|r| r.unwrap())
        .next();
    if let Some(c) = c {
        Ok(Value::Str(String::from(c)))
    } else {
        Err(RuntimeError::IOError(IOError::CouldNotGetChar))
    }
}

pub fn inline(_: &[Value]) -> Result<Value, RuntimeError> {
    let mut str = String::new();
    if let Err(err) = io::stdin().read_line(&mut str) {
        Err(RuntimeError::IOError(IOError::ErrorKind(err.kind())))
    } else {
        Ok(Value::Str(str))
    }
}

pub fn inall(_: &[Value]) -> Result<Value, RuntimeError> {
    let str = CodePoints::from(io::stdin().bytes())
        .map(|r| r.unwrap())
        .collect();
    Ok(Value::Str(str))
}
