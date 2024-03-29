use std::io::{self, Read};

use itertools::Itertools;
use unicode_reader::CodePoints;

use crate::{
    eval::{
        r#type::Type,
        value::{Func, Value},
    },
    parse::ast::{Decl, Declaration},
    prelude::{IOError, RuntimeError},
    util::assert_at_least_args,
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
            env: _,
        }) => {
            let args = args
                .into_iter()
                .filter_map(|decl| match decl {
                    Declaration::Const(Decl {
                        name,
                        type_name,
                        expr: _,
                    }) => Some((name, format!("{type_name}"))),
                    Declaration::Var(Decl {
                        name,
                        type_name,
                        expr: _,
                    }) => Some((name, format!("{type_name}"))),
                })
                .map(|(n, t)| format!("{n} {t}"))
                .collect_vec()
                .join(",");
            let ret = format!("{ret}");

            format!("func({args}){ret}{{...}}")
        }
        Value::Func(Func::Native(_)) => format!("builtin_function"),
        Value::Break(box v) => format!("{}", format_value(v)),
        Value::Type(t) => format!("{t}"),
    }
}

pub fn print(args: &[Value]) -> Result<Value, RuntimeError> {
    for arg in args {
        print!("{}", format_value(&arg));
    }
    Ok(Value::None)
}
pub fn println(args: &[Value]) -> Result<Value, RuntimeError> {
    if args.len() > 0 {
        for arg in args {
            println!("{}", format_value(&arg));
        }
    } else {
        println!();
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
pub fn repeat(args: &[Value]) -> Result<Value, RuntimeError> {
    assert_at_least_args(2, args.len())?;
    let ele = &args[0];
    let count = &args[1];
    let count: usize = count.try_into()?;
    let list = vec![ele.clone(); count];
    Ok(Value::List(list))
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

pub fn inputline(_: &[Value]) -> Result<Value, RuntimeError> {
    let mut str = String::new();
    if let Err(err) = io::stdin().read_line(&mut str) {
        Err(RuntimeError::IOError(IOError::ErrorKind(err.kind())))
    } else {
        Ok(Value::Str(str))
    }
}

pub fn inputall(_: &[Value]) -> Result<Value, RuntimeError> {
    let str = CodePoints::from(io::stdin().bytes())
        .map(|r| r.unwrap())
        .collect();
    Ok(Value::Str(str))
}

pub fn char(args: &[Value]) -> Result<Value, RuntimeError> {
    assert_at_least_args(1, args.len())?;
    let ele = &args[0];
    if let Value::Int(i) = ele {
        if *i > 0 {
            Ok(Value::Str(
                char::from_u32(*i as u32)
                    .expect("i should be positive")
                    .to_string(),
            ))
        } else {
            Err(RuntimeError::IntToStrDomainError { got: *i })
        }
    } else {
        Err(RuntimeError::MismatchedTypes {
            got: vec![ele.to_type()?],
            expected: vec![Type::Int],
        })
    }
}

pub fn ord(args: &[Value]) -> Result<Value, RuntimeError> {
    assert_at_least_args(1, args.len())?;
    let ele = &args[0];
    if let Value::Str(s) = ele {
        if s.len() == 1 {
            Ok(Value::Int(s.chars().next().unwrap() as i64))
        } else {
            Err(RuntimeError::OrdOnlyAllowsSingleCharacterStrings { got: s.clone() })
        }
    } else {
        Err(RuntimeError::MismatchedTypes {
            got: vec![ele.to_type()?],
            expected: vec![Type::Int],
        })
    }
}
