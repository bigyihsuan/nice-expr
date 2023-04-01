use std::{cell::RefCell, collections::HashMap, rc::Rc};

use crate::parse::ast::{Assignment, AssignmentOperator, Declaration, Expr, Literal, Program};

use super::{
    env::{Env, ValueEntry},
    value::Value,
    RuntimeError,
};

pub struct Interpreter {}

impl Interpreter {
    pub fn default_env() -> Rc<RefCell<Env>> {
        Rc::new(RefCell::new(Env::default()))
    }

    pub fn format_value(&self, value: &Value) -> String {
        match value {
            Value::Int(i) => format!("{i}"),
            Value::Dec(d) => format!("{d}"),
            Value::Str(s) => format!("\"{s}\""),
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
        }
    }

    pub fn interpret_program(
        &self,
        program: &Program,
        env: &Rc<RefCell<Env>>,
    ) -> Result<Vec<Value>, RuntimeError> {
        let mut values = Vec::new();
        for expr in program {
            let value = self.interpret_expr(expr, env)?;
            values.push(value);
        }
        Ok(values)
    }

    pub fn interpret_expr(
        &self,
        expr: &Expr,
        env: &Rc<RefCell<Env>>,
    ) -> Result<Value, RuntimeError> {
        match expr {
            Expr::Literal(l) => self.interpret_literal(l, env),
            Expr::Identifier(name) => env
                .borrow()
                .get(name.clone())
                .ok_or_else(|| RuntimeError::IdentifierNotFound(name.clone()))
                .and_then(|ValueEntry { v, c: _, t: _ }| Ok(v)),
            Expr::Unary { op, expr } => todo!("implement unary operators"),
            Expr::Declaration(d) => self.interpret_declaration(d, env),
            Expr::Assignment(a) => self.interpret_assignment(a, env),
        }
    }

    fn interpret_declaration(
        &self,
        decl: &Declaration,
        env: &Rc<RefCell<Env>>,
    ) -> Result<Value, RuntimeError> {
        match decl {
            Declaration::Const {
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
                        .def_const(name.clone(), value.clone(), type_name.clone());
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
        env: &Rc<RefCell<Env>>,
    ) -> Result<Value, RuntimeError> {
        let mut value = self.interpret_expr(assignment.expr.as_ref(), env)?;
        match assignment.op {
            AssignmentOperator::Invalid => {
                return Err(RuntimeError::InvalidAssignmentOperator(
                    assignment.op.clone(),
                ))
            }
            AssignmentOperator::Is => value = value,
        }
        env.borrow_mut().set(assignment.name.clone(), value.clone())
    }

    pub fn interpret_literal(
        &self,
        literal: &Literal,
        env: &Rc<RefCell<Env>>,
    ) -> Result<Value, RuntimeError> {
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
}
