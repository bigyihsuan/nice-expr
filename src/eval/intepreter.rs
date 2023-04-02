use std::{
    cell::RefCell,
    collections::HashMap,
    io::{self, Read},
    rc::Rc,
};

use unicode_reader::CodePoints;

use crate::{
    parse::ast::{
        Assignment, BinaryExpr, BinaryOperator, Declaration, Expr, Literal, Operator, Program,
        UnaryExpr,
    },
    prelude::{IOError, RuntimeError},
    util::assert_at_least_args,
};

use super::{
    env::{Env, SEnv, ValueEntry},
    r#type::Type,
    value::Value,
};

mod operators;

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
            Value::Break(box v) => format!("{}", self.format_value(v)),
        }
    }

    pub fn interpret_program(&self, program: &Program, env: &SEnv) -> Result<(), RuntimeError> {
        for expr in program {
            let value = self.interpret_expr(expr, env)?;
            if let Value::None = value {
                continue;
            }
        }
        Ok(())
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

            Expr::Minus(e) => self.interpret_minus(e, env),
            Expr::Not(e) => self.interpret_not(e, env),
            Expr::Indexing(e) => self.interpret_indexing(e, env),
            Expr::Multiplication(e) => self.interpret_multiplication(e, env),
            Expr::Addition(e) => self.interpret_addition(e, env),
            Expr::Comparison(e) => self.interpret_comparison(e, env),
            Expr::Logical(e) => self.interpret_logical(e, env),

            Expr::Block(exprs) => {
                let block_env = Env::extend(env.clone());
                let mut last_val = Value::None;
                for expr in exprs {
                    last_val = self.interpret_expr(&expr, &block_env)?;
                    if let Value::Break(_) = last_val {
                        return Ok(last_val);
                    }
                }
                return Ok(last_val);
            }
            Expr::Return(Some(box e)) => Ok(Value::Break(Box::new(self.interpret_expr(e, env)?))),
            Expr::Return(None) => Ok(Value::Break(Box::new(Value::None))),
            Expr::If {
                condition,
                when_true,
                when_false,
            } => {
                let condition: bool = self.interpret_expr(&condition, env)?.try_into()?;
                if condition {
                    self.interpret_expr(&when_true, env)
                } else if let Some(when_false) = when_false {
                    self.interpret_expr(&when_false, env)
                } else {
                    Ok(Value::None)
                }
            }
            Expr::For { vars, body } => {
                // init the vars
                let loop_env = Env::extend(env.clone());
                for v in vars {
                    self.interpret_expr(v, &loop_env)?;
                }
                // run the loop
                loop {
                    for expr in body {
                        let last_val = self.interpret_expr(&expr, &loop_env)?;
                        if let Value::Break(_) = last_val {
                            return Ok(last_val.unbreak());
                        }
                    }
                }
            }
            Expr::Break(Some(box e)) => Ok(Value::Break(Box::new(self.interpret_expr(e, env)?))),
            Expr::Break(None) => Ok(Value::Break(Box::new(Value::None))),
        }
    }

    fn interpret_declaration(&self, decl: &Declaration, env: &SEnv) -> Result<Value, RuntimeError> {
        match decl {
            Declaration::Const {
                name,
                type_name: decl_type,
                expr: value,
            } => {
                let v = {
                    if let Some(v) = value {
                        self.interpret_expr(v, env)?
                    } else {
                        Value::default(decl_type.clone())
                    }
                };
                let val_type = v.to_type()?;

                let inferred_type = val_type.infer_contained_type(decl_type);
                if inferred_type.is_none() {
                    return Err(RuntimeError::MismatchedTypes {
                        got: vec![val_type],
                        expected: vec![decl_type.clone()],
                    });
                }
                let inferred_type = inferred_type.unwrap();

                if let Some(_) = value && inferred_type != decl_type.clone() {
                    return Err(RuntimeError::MismatchedTypes {
                        got: vec![val_type],
                        expected: vec![decl_type.clone()],
                    });
                }
                let result = env
                    .borrow_mut()
                    .def_const(name.clone(), v.clone(), decl_type.clone());
                if let Err(name) = result {
                    Err(RuntimeError::IdentifierNotFound(name.clone()))
                } else {
                    Ok(v)
                }
            }
            Declaration::Var {
                name,
                type_name: decl_type,
                expr: value,
            } => {
                let v = {
                    if let Some(v) = value {
                        self.interpret_expr(v, env)?
                    } else {
                        Value::default(decl_type.clone())
                    }
                };
                let val_type = v.to_type()?;

                let inferred_type = val_type.infer_contained_type(decl_type);
                if inferred_type.is_none() {
                    return Err(RuntimeError::MismatchedTypes {
                        got: vec![val_type],
                        expected: vec![decl_type.clone()],
                    });
                }
                let inferred_type = inferred_type.unwrap();

                if let Some(_) = value && inferred_type != decl_type.clone() {
                    return Err(RuntimeError::MismatchedTypes {
                        got: vec![val_type],
                        expected: vec![decl_type.clone()],
                    });
                }
                let result = env
                    .borrow_mut()
                    .def_var(name.clone(), v.clone(), decl_type.clone());
                if let Err(name) = result {
                    Err(RuntimeError::IdentifierNotFound(name.clone()))
                } else {
                    Ok(v)
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

        result = match assignment.op {
            BinaryOperator::Is => value,
            BinaryOperator::And => operators::band(result, value, env)?,
            BinaryOperator::Or => operators::bor(result, value, env)?,
            BinaryOperator::Greater => operators::gt(result, value, env)?,
            BinaryOperator::Less => operators::lt(result, value, env)?,
            BinaryOperator::GreaterEqual => operators::ge(result, value, env)?,
            BinaryOperator::LessEqual => operators::le(result, value, env)?,
            BinaryOperator::Equal => operators::eq(result, value, env)?,
            BinaryOperator::NotEqual => operators::ne(result, value, env)?,
            BinaryOperator::Add | BinaryOperator::Subtract => {
                self.addition(assignment.op.clone(), result, value, env)?
            }
            BinaryOperator::Times | BinaryOperator::Divide | BinaryOperator::Modulo => {
                self.multiplication(assignment.op.clone(), result, value, env)?
            }
            _ => {
                return Err(RuntimeError::InvalidAssignmentOperator(
                    assignment.op.clone(),
                ))
            }
        };
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

                // check that all elements match the same type
                // assume that the type is the first element
                if list.len() > 0 {
                    let first_type = list.get(0).unwrap().to_type()?;
                    for e in list.iter() {
                        let e_type = e.to_type()?;
                        if e_type != first_type {
                            return Err(RuntimeError::MismatchedTypes {
                                got: vec![Type::List(Box::new(e_type))],
                                expected: vec![Type::List(Box::new(first_type))],
                            });
                        }
                    }
                }
                Ok(Value::List(list))
            }
            Literal::Map(m) => {
                let mut map = HashMap::new();
                for (k, v) in m {
                    map.insert(self.interpret_expr(k, &env)?, self.interpret_expr(v, &env)?);
                }

                // check that all elements match the same type
                // assume that the type is the first element
                let entries = map.clone().into_iter().collect::<Vec<(Value, Value)>>();
                if map.len() > 0 {
                    let first = entries.first().unwrap();
                    let ktype = first.0.to_type()?;
                    let vtype = first.1.to_type()?;
                    for e in map.iter() {
                        let ek_type = e.0.to_type()?;
                        let ev_type = e.1.to_type()?;
                        if (&ktype, &vtype) != (&ek_type, &ev_type) {
                            return Err(RuntimeError::MismatchedTypes {
                                got: vec![Type::Map(Box::new(ek_type), Box::new(ev_type))],
                                expected: vec![Type::Map(Box::new(ktype), Box::new(vtype))],
                            });
                        }
                    }
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
            "inputchar" => self.builtin_inputchar(),
            "inputline" => self.builtin_inline(),
            "inputall" => self.builtin_inall(),
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

        let start: isize = start.try_into()?;
        let end: isize = end.try_into()?;
        let step: isize = step.try_into()?;

        if step > 0 {
            Ok(Value::List(
                (start..end)
                    .step_by(step as usize)
                    .map(|i| Value::Int(i as i64))
                    .collect(),
            ))
        } else if step < 0 {
            Ok(Value::List(
                (end..start)
                    .step_by(step as usize)
                    .rev()
                    .map(|i| Value::Int(i as i64))
                    .collect(),
            ))
        } else {
            Err(RuntimeError::InvalidRangeStep(step))
        }
    }

    fn builtin_inputchar(&self) -> Result<Value, RuntimeError> {
        let c = CodePoints::from(io::stdin().bytes())
            .map(|r| r.unwrap())
            .next();
        if let Some(c) = c {
            Ok(Value::Str(String::from(c)))
        } else {
            Err(RuntimeError::IOError(IOError::CouldNotGetChar))
        }
    }

    fn builtin_inline(&self) -> Result<Value, RuntimeError> {
        let mut str = String::new();
        if let Err(err) = io::stdin().read_line(&mut str) {
            Err(RuntimeError::IOError(IOError::ErrorKind(err.kind())))
        } else {
            Ok(Value::Str(str))
        }
    }

    fn builtin_inall(&self) -> Result<Value, RuntimeError> {
        let str = CodePoints::from(io::stdin().bytes())
            .map(|r| r.unwrap())
            .collect();
        Ok(Value::Str(str))
    }

    fn interpret_minus(&self, e: &UnaryExpr, env: &SEnv) -> Result<Value, RuntimeError> {
        let val = self.interpret_expr(&e.expr, env)?;
        let val_type = val.to_type()?;
        match val_type {
            Type::Int => operators::ineg(val, env),
            Type::Dec => operators::fneg(val, env),
            _ => Err(RuntimeError::InvalidOperatorOnTypes {
                op: Operator::UnaryOperator(e.op.clone()),
                types: vec![val_type],
            }),
        }
    }

    fn interpret_not(&self, e: &UnaryExpr, env: &SEnv) -> Result<Value, RuntimeError> {
        let val = self.interpret_expr(&e.expr, env)?;
        let val_type = val.to_type()?;
        match val_type {
            Type::Bool => operators::bnot(val, env),
            _ => Err(RuntimeError::InvalidOperatorOnTypes {
                op: Operator::UnaryOperator(e.op.clone()),
                types: vec![val_type],
            }),
        }
    }

    fn interpret_indexing(&self, e: &BinaryExpr, env: &SEnv) -> Result<Value, RuntimeError> {
        let left = self.interpret_expr(&e.left, env)?;
        let right = self.interpret_expr(&e.right, env)?;
        let l_type = left.to_type()?;
        let r_type = right.to_type()?;
        match (&l_type, &r_type) {
            (Type::Str, Type::Int) => operators::sidx(left, right, env),
            (Type::List(_), Type::Int) => operators::lidx(left, right, env),
            (Type::Map(box k, _), t) if t == k => operators::midx(left, right, env),
            _ => Err(RuntimeError::InvalidOperatorOnTypes {
                op: Operator::BinaryOperator(e.op.clone()),
                types: vec![l_type, r_type],
            }),
        }
    }

    fn interpret_multiplication(&self, e: &BinaryExpr, env: &SEnv) -> Result<Value, RuntimeError> {
        let left = self.interpret_expr(&e.left, env)?;
        let right = self.interpret_expr(&e.right, env)?;
        self.multiplication(e.op.clone(), left, right, env)
    }

    fn interpret_addition(&self, e: &BinaryExpr, env: &SEnv) -> Result<Value, RuntimeError> {
        let left = self.interpret_expr(&e.left, env)?;
        let right = self.interpret_expr(&e.right, env)?;
        self.addition(e.op.clone(), left, right, env)
    }

    fn interpret_comparison(&self, e: &BinaryExpr, env: &SEnv) -> Result<Value, RuntimeError> {
        let left = self.interpret_expr(&e.left, env)?;
        let right = self.interpret_expr(&e.right, env)?;

        match &e.op {
            BinaryOperator::Equal => operators::eq(left, right, env),
            BinaryOperator::NotEqual => operators::ne(left, right, env),
            BinaryOperator::Greater => operators::gt(left, right, env),
            BinaryOperator::GreaterEqual => operators::ge(left, right, env),
            BinaryOperator::Less => operators::lt(left, right, env),
            BinaryOperator::LessEqual => operators::le(left, right, env),
            _ => Err(RuntimeError::InvalidOperatorOnTypes {
                op: Operator::BinaryOperator(e.op.clone()),
                types: vec![left.to_type()?, right.to_type()?],
            }),
        }
    }

    fn interpret_logical(&self, e: &BinaryExpr, env: &SEnv) -> Result<Value, RuntimeError> {
        let left = self.interpret_expr(&e.left, env)?;
        let right = self.interpret_expr(&e.right, env)?;

        match e.op {
            BinaryOperator::And => operators::band(left, right, env),
            BinaryOperator::Or => operators::bor(left, right, env),
            _ => Err(RuntimeError::InvalidOperatorOnTypes {
                op: Operator::BinaryOperator(e.op.clone()),
                types: vec![Type::Bool, Type::Bool],
            }),
        }
    }

    fn addition(
        &self,
        op: BinaryOperator,
        left: Value,
        right: Value,
        env: &SEnv,
    ) -> Result<Value, RuntimeError> {
        let l_type = left.to_type()?;
        let r_type = right.to_type()?;
        match (&op, &l_type, &r_type) {
            (BinaryOperator::Add, Type::Int, Type::Int) => operators::iadd(left, right, env),
            (BinaryOperator::Add, Type::Dec, Type::Dec) => operators::fadd(left, right, env),
            (BinaryOperator::Add, Type::Str, Type::Str) => operators::sadd(left, right, env),
            (BinaryOperator::Add, Type::List(_), Type::List(_)) => {
                operators::ladd(left, right, env)
            }
            (BinaryOperator::Subtract, Type::Int, Type::Int) => operators::isub(left, right, env),
            (BinaryOperator::Subtract, Type::Dec, Type::Dec) => operators::fsub(left, right, env),
            (BinaryOperator::Subtract, Type::Str, Type::Str) => operators::ssub(left, right, env),
            (BinaryOperator::Subtract, Type::List(_), Type::List(_)) => {
                operators::lsub(left, right, env)
            }
            _ => Err(RuntimeError::InvalidOperatorOnTypes {
                op: Operator::BinaryOperator(op.clone()),
                types: vec![l_type, r_type],
            }),
        }
    }

    fn multiplication(
        &self,
        op: BinaryOperator,
        left: Value,
        right: Value,
        env: &SEnv,
    ) -> Result<Value, RuntimeError> {
        let l_type = left.to_type()?;
        let r_type = right.to_type()?;
        match (&op, &l_type, &r_type) {
            (BinaryOperator::Times, Type::Int, Type::Int) => operators::imul(left, right, env),
            (BinaryOperator::Times, Type::Dec, Type::Dec) => operators::fmul(left, right, env),
            (BinaryOperator::Divide, Type::Int, Type::Int) => operators::idiv(left, right, env),
            (BinaryOperator::Divide, Type::Dec, Type::Dec) => operators::fdiv(left, right, env),
            (BinaryOperator::Modulo, Type::Int, Type::Int) => operators::imod(left, right, env),
            _ => Err(RuntimeError::InvalidOperatorOnTypes {
                op: Operator::BinaryOperator(op.clone()),
                types: vec![l_type, r_type],
            }),
        }
    }
}
