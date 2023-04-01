use std::{cell::RefCell, collections::HashMap, rc::Rc};

use crate::{
    parse::ast::{Assignment, AssignmentOperator, Declaration, Expr, Literal, Program},
    util::assert_at_least_args,
};

use super::{
    env::{Env, SEnv, ValueEntry},
    value::Value,
    RuntimeError,
};

pub struct Interpreter {}

impl Interpreter {
    pub fn default_env() -> SEnv {
        Rc::new(RefCell::new(Env::default()))
    }

    pub fn format_value(&self, value: &Value) -> String {
        match value {
            Value::None => "None".into(),
            Value::Int(i) => format!("{i}"),
            Value::Dec(d) => format!("{d}"),
            Value::Str(s) => format!("{s}"),
            Value::Bool(b) => format!("{b}"),
            Value::List(l) => {
                let l = l
                    .into_iter()
                    .map(|e| self.format_value(&e))
                    .collect::<Vec<String>>()
                    .join(",");
                format!("[{l}]")
            }
            Value::Map(m) => {
                let m = m
                    .into_iter()
                    .map(|(k, v)| (self.format_value(&k), self.format_value(&v)))
                    .map(|(k, v)| format!("{k}:{v}"))
                    .collect::<Vec<String>>()
                    .join(",");
                format!("<|{m}|>")
            }
            Value::Func(_) => todo!(),
        }
    }

    pub fn interpret_program(
        &self,
        program: &Program,
        env: &SEnv,
    ) -> Result<Vec<Value>, RuntimeError> {
        let mut values = Vec::new();
        for expr in program {
            let value = self.interpret_expr(expr, env)?;
            if let Value::None = value {
                continue;
            }
            values.push(value);
        }
        Ok(values)
    }

    pub fn interpret_expr(&self, expr: &Expr, env: &SEnv) -> Result<Value, RuntimeError> {
        match expr {
            Expr::Literal(l) => self.interpret_literal(l, env),
            Expr::Identifier(name) => env
                .borrow()
                .get(name.clone())
                .ok_or_else(|| RuntimeError::IdentifierNotFound(name.clone()))
                .and_then(|ValueEntry { v, c: _, t: _ }| Ok(v)),
            Expr::Declaration(d) => self.interpret_declaration(d, env),
            Expr::Assignment(a) => self.interpret_assignment(a, env),
            Expr::FunctionCall { name, args } => self.interpret_function_call(name, args, env),

            Expr::Unary { op, expr } => todo!("implement unary operators"),
        }
    }

    fn interpret_declaration(&self, decl: &Declaration, env: &SEnv) -> Result<Value, RuntimeError> {
        match decl {
            Declaration::Const {
                name,
                type_name: decl_type,
                expr: value,
            } => {
                let value = self.interpret_expr(value, env)?;
                let val_type = value.to_type()?;

                let inferred_type = val_type.infer_contained_type(decl_type);
                if inferred_type.is_none() {
                    return Err(RuntimeError::MismatchedTypes {
                        got: val_type,
                        expected: decl_type.clone(),
                    });
                }
                let inferred_type = inferred_type.unwrap();

                if inferred_type != decl_type.clone() {
                    return Err(RuntimeError::MismatchedTypes {
                        got: val_type,
                        expected: decl_type.clone(),
                    });
                }
                let result =
                    env.borrow_mut()
                        .def_const(name.clone(), value.clone(), decl_type.clone());
                if let Err(name) = result {
                    Err(RuntimeError::IdentifierNotFound(name.clone()))
                } else {
                    Ok(value)
                }
            }
            Declaration::Var {
                name,
                type_name,
                expr: value,
            } => {
                let value = self.interpret_expr(value, env)?;
                let t = value.to_type()?;
                if t != type_name.clone() {
                    return Err(RuntimeError::MismatchedTypes {
                        got: t,
                        expected: type_name.clone(),
                    });
                }
                let result =
                    env.borrow_mut()
                        .def_var(name.clone(), value.clone(), type_name.clone());
                if let Err(name) = result {
                    Err(RuntimeError::IdentifierNotFound(name.clone()))
                } else {
                    Ok(value)
                }
            }
        }
    }

    pub fn interpret_assignment(
        &self,
        assignment: &Assignment,
        env: &SEnv,
    ) -> Result<Value, RuntimeError> {
        let entry = env
            .borrow()
            .get(assignment.name.clone())
            .ok_or(RuntimeError::IdentifierNotFound(assignment.name.clone()))?;
        let mut result = entry.v;
        let value = self.interpret_expr(assignment.expr.as_ref(), env)?;

        match assignment.op {
            AssignmentOperator::Is => result = value,
            // TODO: other assignment operators
            _ => {
                return Err(RuntimeError::InvalidAssignmentOperator(
                    assignment.op.clone(),
                ))
            }
        }
        env.borrow_mut()
            .set(assignment.name.clone(), result.clone())
    }

    pub fn interpret_literal(&self, literal: &Literal, env: &SEnv) -> Result<Value, RuntimeError> {
        match literal {
            Literal::Int(i) => Ok(Value::Int(*i)),
            Literal::Dec(d) => Ok(Value::Dec(*d)),
            Literal::Str(s) => Ok(Value::Str(s.clone())),
            Literal::Bool(b) => Ok(Value::Bool(*b)),
            Literal::List(l) => {
                let mut list = Vec::new();
                for e in l {
                    list.push(self.interpret_expr(e, &env)?);
                }
                Ok(Value::List(list))
            }
            Literal::Map(m) => {
                let mut map = HashMap::new();
                for (k, v) in m {
                    map.insert(self.interpret_expr(k, &env)?, self.interpret_expr(v, &env)?);
                }
                Ok(Value::Map(map))
            }
        }
    }

    pub fn interpret_function_call(
        &self,
        name: &str,
        args: &[Expr],
        env: &SEnv,
    ) -> Result<Value, RuntimeError> {
        let args = args
            .into_iter()
            .map(|e| self.interpret_expr(e, env))
            .collect::<Vec<_>>();
        let errors = args
            .iter()
            .filter_map(|r| r.as_ref().err())
            .collect::<Vec<_>>();
        if errors.len() > 0 {
            return Err(errors[0].clone());
        }
        let args = args
            .into_iter()
            .filter_map(|r| r.ok())
            .collect::<Vec<Value>>();
        match name {
            "print" => self.builtin_print(&args),
            "println" => self.builtin_println(&args),
            "len" => self.builtin_len(&args),
            "range" => self.builtin_range(&args),
            _ => {
                todo!("user-defined functions: got {}", name);
            }
        }
    }

    fn builtin_print(&self, args: &[Value]) -> Result<Value, RuntimeError> {
        for arg in args {
            print!("{}", self.format_value(&arg));
        }
        Ok(Value::None)
    }
    fn builtin_println(&self, args: &[Value]) -> Result<Value, RuntimeError> {
        for arg in args {
            println!("{}", self.format_value(&arg));
        }
        Ok(Value::None)
    }

    fn builtin_len(&self, args: &[Value]) -> Result<Value, RuntimeError> {
        assert_at_least_args(1, args.len())?;
        let val = &args[0];
        match val {
            Value::Str(v) => Ok(Value::Int(v.chars().collect::<Vec<_>>().len() as i64)),
            Value::List(v) => Ok(Value::Int(v.len() as i64)),
            Value::Map(v) => Ok(Value::Int(v.len() as i64)),
            _ => Err(RuntimeError::TakingLenOfLengthless {
                got: val.to_type()?,
            }),
        }
    }

    fn builtin_range(&self, args: &[Value]) -> Result<Value, RuntimeError> {
        assert_at_least_args(3, args.len())?;
        let start = &args[0];
        let end = &args[1];
        let step = &args[2];

        let start: usize = start.try_into()?;
        let end: usize = end.try_into()?;
        let step: usize = step.try_into()?;

        Ok(Value::List(
            (start..end)
                .step_by(step)
                .map(|i| Value::Int(i as i64))
                .collect(),
        ))
    }
}
