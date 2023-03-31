use std::{cell::RefCell, collections::HashMap, rc::Rc};

use crate::parse::ast::{Expr, Literal, Program};

use super::{env::Env, value::Value, RuntimeError};

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

    pub fn parse_program(
        &self,
        program: &Program,
        env: &Rc<RefCell<Env>>,
    ) -> Result<Vec<Value>, RuntimeError> {
        let mut values = Vec::new();
        for expr in program {
            let value = self.parse_expr(expr, env)?;
            values.push(value);
        }
        Ok(values)
    }

    pub fn parse_expr(&self, expr: &Expr, env: &Rc<RefCell<Env>>) -> Result<Value, RuntimeError> {
        match expr {
            Expr::Literal(l) => Ok(self.parse_literal(l, env)?),
            Expr::Identifier { name } => env
                .borrow()
                .get(name.clone())
                .ok_or_else(|| RuntimeError::IdentifierNotFound)
                .and_then(|(v, c)| Ok(v)),
            Expr::Unary { op, expr } => todo!(),
        }
    }

    pub fn parse_literal(
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
                    list.push(self.parse_expr(e, &env)?);
                }
                Ok(Value::List(list))
            }
            Literal::Map(m) => {
                let mut map = HashMap::new();
                for (k, v) in m {
                    map.insert(self.parse_expr(k, &env)?, self.parse_expr(v, &env)?);
                }
                Ok(Value::Map(map))
            }
        }
    }
}
