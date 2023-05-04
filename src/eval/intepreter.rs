use std::{cell::RefCell, collections::HashMap, rc::Rc};

use itertools::Itertools;

use crate::{
    parse::ast::{
        Assignment, BinaryExpr, BinaryOperator, Decl, Declaration, Expr, IndexKind, Indexing,
        Literal, Operator, Program, UnaryExpr,
    },
    prelude::RuntimeError,
    util::assert_exactly_args,
};

use super::{
    env::{Env, SEnv, ValueEntry},
    r#type::Type,
    value::{Func, Value},
};

pub mod builtin;
mod operators;

pub struct Interpreter {}

impl Interpreter {
    pub fn default_env() -> SEnv {
        Rc::new(RefCell::new(Env::default()))
    }

    pub fn interpret_program(&self, program: &Program, env: &SEnv) -> Result<Value, RuntimeError> {
        let mut last_value = Value::None;
        for expr in program {
            last_value = self.interpret_expr(expr, env)?;
        }
        Ok(last_value)
    }

    pub fn interpret_expr(&self, expr: &Expr, env: &SEnv) -> Result<Value, RuntimeError> {
        match expr {
            Expr::Literal(l) => self.interpret_literal(l, env),
            Expr::Identifier(name) => Ok(env.borrow().get(name.clone())?.v),
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
                return self.interpret_program(exprs, &block_env);
            }
            Expr::Return(Some(box e)) => Ok(Value::Break(Box::new(self.interpret_expr(e, env)?))),
            Expr::Return(None) => Ok(Value::Break(Box::new(Value::None))),
            Expr::If {
                vars,
                condition,
                when_true,
                when_false,
            } => {
                let mut if_env = Env::extend(env.clone());
                for v in vars {
                    self.interpret_declaration(v, &mut if_env)?;
                }
                let condition: bool = self.interpret_expr(&condition, &if_env)?.try_into()?;
                if condition {
                    self.interpret_expr(&when_true, &if_env)
                } else if let Some(when_false) = when_false {
                    self.interpret_expr(&when_false, &if_env)
                } else {
                    Ok(Value::None)
                }
            }
            Expr::For { vars, body } => {
                // init the vars
                let mut loop_env = Env::extend(env.clone());
                for v in vars {
                    self.interpret_declaration(v, &mut loop_env)?;
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
            Expr::ForWhile {
                vars,
                condition,
                body,
            } => {
                // init the vars
                let mut loop_env = Env::extend(env.clone());
                for v in vars {
                    self.interpret_declaration(v, &mut loop_env)?;
                }
                // run the loop
                let mut last_val = Value::None;
                loop {
                    let continue_loop: bool =
                        self.interpret_expr(condition, &loop_env)?.try_into()?;
                    if !continue_loop {
                        return Ok(last_val);
                    }
                    for expr in body {
                        last_val = self.interpret_expr(&expr, &loop_env)?;
                        if let Value::Break(_) = last_val {
                            return Ok(last_val.unbreak());
                        }
                    }
                }
            }
            Expr::ForIn {
                vars,
                box collection,
                body,
            } => {
                let collection = self.interpret_expr(collection, env)?;
                // init the vars
                let mut loop_env = Env::extend(env.clone());
                if vars.len() < 1 {
                    return Err(RuntimeError::NotEnoughArguments {
                        want: 1,
                        got: vars.len(),
                    });
                }
                for v in vars {
                    self.interpret_declaration(v, &mut loop_env)?;
                }
                // check for the first decl to be a var
                let first = vars.first().unwrap();
                if let Declaration::Var(Decl { name, .. }) = first {
                    let mut last_val = Value::None;
                    match collection {
                        Value::Str(s) => {
                            for c in s.chars() {
                                let c = Value::Str(c.into());
                                loop_env.borrow_mut().set(name.clone(), c, false)?;
                                for expr in body {
                                    last_val = self.interpret_expr(&expr, &loop_env)?;
                                    if let Value::Break(_) = last_val {
                                        return Ok(last_val.unbreak());
                                    }
                                }
                            }
                        }
                        Value::List(l) => {
                            for e in l {
                                loop_env.borrow_mut().set(name.clone(), e, false)?;
                                for expr in body {
                                    last_val = self.interpret_expr(&expr, &loop_env)?;
                                    if let Value::Break(_) = last_val {
                                        return Ok(last_val.unbreak());
                                    }
                                }
                            }
                        }
                        Value::Map(m) => {
                            for (k, _) in m {
                                loop_env.borrow_mut().set(name.clone(), k, false)?;
                                for expr in body {
                                    last_val = self.interpret_expr(&expr, &loop_env)?;
                                    if let Value::Break(_) = last_val {
                                        return Ok(last_val.unbreak());
                                    }
                                }
                            }
                        }
                        _ => {
                            return Err(RuntimeError::MismatchedTypes {
                                got: vec![collection.to_type()?],
                                expected: vec![
                                    Type::Str,
                                    Type::List(Box::new(Type::None)),
                                    Type::Map(Box::new(Type::None), Box::new(Type::None)),
                                ],
                            })
                        }
                    }
                    Ok(last_val)
                } else {
                    return Err(RuntimeError::ForInFirstDeclMustBeVar);
                }
            }
            Expr::Break(Some(box e)) => Ok(Value::Break(Box::new(self.interpret_expr(e, env)?))),
            Expr::Break(None) => Ok(Value::Break(Box::new(Value::None))),
            Expr::FunctionDefinition { args, ret, body } => {
                self.interpret_function_definition(args, ret, body, env)
            }
            Expr::TypeName(t) => self.interpret_type_name(t, env),
            Expr::TypeCast(expr, t) => self.interpret_type_cast(expr, t, env),
        }
    }

    fn interpret_declaration(&self, decl: &Declaration, env: &SEnv) -> Result<Value, RuntimeError> {
        match decl {
            Declaration::Const(Decl {
                name,
                type_name: decl_type,
                expr: value,
            }) => {
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
            Declaration::Var(Decl {
                name,
                type_name: decl_type,
                expr: value,
            }) => {
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
        let entry = env.borrow().get(assignment.name.clone())?;
        let mut result = entry.v.clone();
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
        // make sure the collection at the index is set to the result
        let result = if let Some(box index) = &assignment.index {
            let collection = entry.v.clone();
            let index = self.interpret_expr(index, env)?;
            let l_type = collection.unbreak().to_type()?;
            let r_type = index.unbreak().to_type()?;
            match (&l_type, &r_type) {
                (Type::Str, Type::Int) => operators::ssetidx(collection, index, result, env),
                (Type::List(_), Type::Int) => operators::lsetidx(collection, index, result, env),
                (Type::Map(box k, _), t) if t == k || collection.is_empty() => {
                    operators::msetidx(collection, index, result, env)
                }
                _ => Err(RuntimeError::InvalidOperatorOnTypes {
                    op: Operator::BinaryOperator(BinaryOperator::SetIndexing),
                    types: vec![l_type, r_type],
                }),
            }?
        } else {
            result
        };
        env.borrow_mut()
            .set(assignment.name.clone(), result.clone(), false)
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
                let entries = map.clone().into_iter().collect_vec();
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
            .collect_vec();
        let errors = args.iter().filter_map(|r| r.as_ref().err()).collect_vec();
        if errors.len() > 0 {
            return Err(errors[0].clone());
        }
        let args = args
            .into_iter()
            .filter_map(|r| r.ok())
            .collect::<Vec<Value>>();

        let entry = env.borrow().get(name.into())?;
        let ValueEntry { v, t: _, c: _ } = entry;
        let v = v.unbreak();

        match v {
            Value::Func(Func::Native(f)) => f.1(&args), // builtin
            Value::Func(Func::Declared {
                decls,
                ret: _,
                body,
                env: closed_env,
            }) => {
                // user-declared
                assert_exactly_args(decls.len(), args.len())?;
                for (decl, arg) in decls.iter().zip(args.into_iter()) {
                    // fill the variables with the argument values
                    let (Declaration::Const(Decl { name, .. })
                    | Declaration::Var(Decl { name, .. })) = &decl;
                    closed_env.borrow_mut().set(name.into(), arg, true)?;
                }
                // run the function body
                self.interpret_expr(&Expr::Block(body), &closed_env)
            }
            _ => Err(RuntimeError::CallingNonCallable { name: name.into() }),
        }
    }

    fn interpret_minus(&self, e: &UnaryExpr, env: &SEnv) -> Result<Value, RuntimeError> {
        let val = self.interpret_expr(&e.expr, env)?.unbreak();
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

    fn interpret_indexing(&self, e: &Indexing, env: &SEnv) -> Result<Value, RuntimeError> {
        let collection = self.interpret_expr(&e.collection, env)?.unbreak();
        let start_index = match &e.index {
            IndexKind::Single { index } => self.interpret_expr(index, env)?.unbreak(),
            IndexKind::Range { start, end: _ } => self.interpret_expr(start, env)?.unbreak(),
        };
        let end_index = if let IndexKind::Range { start: _, box end } = &e.index {
            Some(self.interpret_expr(end, env)?.unbreak())
        } else {
            None
        };
        let l_type = collection.to_type()?;
        let start_r_type = start_index.to_type()?;
        let end_r_type = end_index.clone().and_then(|e| e.to_type().ok());
        match (&l_type, &start_r_type, &end_r_type) {
            (Type::Str, Type::Int, _) => {
                operators::sgetidx(collection, start_index, end_index, env)
            }
            (Type::List(_), Type::Int, _) => {
                operators::lgetidx(collection, start_index, end_index, env)
            }
            (Type::Map(box k, _), t, None) if t == k => {
                // can't do a range on a map
                operators::mgetidx(collection, start_index, end_index, env)
            }
            _ => Err(RuntimeError::InvalidOperatorOnTypes {
                op: Operator::BinaryOperator(e.op.clone()),
                types: vec![
                    l_type,
                    start_r_type.clone(),
                    end_r_type.unwrap_or(start_r_type),
                ],
            }),
        }
    }

    fn interpret_multiplication(&self, e: &BinaryExpr, env: &SEnv) -> Result<Value, RuntimeError> {
        let left = self.interpret_expr(&e.left, env)?.unbreak();
        let right = self.interpret_expr(&e.right, env)?.unbreak();
        self.multiplication(e.op.clone(), left, right, env)
    }

    fn interpret_addition(&self, e: &BinaryExpr, env: &SEnv) -> Result<Value, RuntimeError> {
        let left = self.interpret_expr(&e.left, env)?.unbreak();
        let right = self.interpret_expr(&e.right, env)?.unbreak();
        self.addition(e.op.clone(), left, right, env)
    }

    fn interpret_comparison(&self, e: &BinaryExpr, env: &SEnv) -> Result<Value, RuntimeError> {
        let left = self.interpret_expr(&e.left, env)?.unbreak();
        let right = self.interpret_expr(&e.right, env)?.unbreak();

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
        let left = self.interpret_expr(&e.left, env)?.unbreak();
        let right = self.interpret_expr(&e.right, env)?.unbreak();

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
            (BinaryOperator::Add, Type::Map(_, _), Type::Map(_, _)) => {
                operators::madd(left, right, env)
            }
            (BinaryOperator::Subtract, Type::Int, Type::Int) => operators::isub(left, right, env),
            (BinaryOperator::Subtract, Type::Dec, Type::Dec) => operators::fsub(left, right, env),
            (BinaryOperator::Subtract, Type::Str, Type::Str) => operators::ssub(left, right, env),
            (BinaryOperator::Subtract, Type::List(_), Type::List(_)) => {
                operators::lsub(left, right, env)
            }
            (BinaryOperator::Subtract, Type::Map(_, _), Type::Map(_, _)) => {
                operators::msub(left, right, env)
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

    fn interpret_function_definition(
        &self,
        args: &[Declaration],
        ret: &Type,
        body: &[Expr],
        env: &SEnv,
    ) -> Result<Value, RuntimeError> {
        // set up any variables in the args
        let mut func_env = Env::extend(env.clone());
        for decl in args {
            self.interpret_declaration(decl, &mut func_env)?;
        }

        Ok(Value::Func(Func::Declared {
            decls: args.into(),
            ret: ret.clone(),
            body: body.into(),
            env: func_env.clone(),
        }))
    }

    fn interpret_type_name(&self, t: &Type, _env: &SEnv) -> Result<Value, RuntimeError> {
        Ok(Value::Type(t.clone()))
    }

    fn interpret_type_cast(
        &self,
        expr: &Expr,
        t: &Type,
        env: &SEnv,
    ) -> Result<Value, RuntimeError> {
        let value = self.interpret_expr(expr, env)?;
        let value = value.unbreak();
        if let Value::Type(t) = self.interpret_type_name(t, env)? {
            self.cast_to_type(&value, &t)
        } else {
            Err(RuntimeError::NotImplemented)
        }
    }

    fn cast_to_type(&self, value: &Value, t: &Type) -> Result<Value, RuntimeError> {
        match (&value, &t) {
            (Value::Int(i), Type::Int) => Ok(Value::Int(*i)),
            (Value::Dec(i), Type::Int) => Ok(Value::Int(*i as i64)),
            (Value::Str(i), Type::Int) => {
                if let Ok(i) = parse_int::parse::<i64>(&i) {
                    Ok(Value::Int(i))
                } else {
                    Err(RuntimeError::StringToValueParseFailed(
                        value.clone(),
                        t.clone(),
                    ))
                }
            }
            (Value::Bool(i), Type::Int) => Ok(Value::Int(*i as i64)),

            (Value::Int(i), Type::Dec) => Ok(Value::Dec(*i as f64)),
            (Value::Dec(i), Type::Dec) => Ok(Value::Dec(*i)),
            (Value::Str(i), Type::Dec) => {
                if let Ok(i) = parse_int::parse::<f64>(&i) {
                    Ok(Value::Dec(i))
                } else {
                    Err(RuntimeError::StringToValueParseFailed(
                        value.clone(),
                        t.clone(),
                    ))
                }
            }
            (Value::Bool(i), Type::Dec) => Ok(Value::Dec(*i as i64 as f64)),

            (Value::Int(i), Type::Str) => Ok(Value::Str(format!("{i}"))),
            (Value::Dec(i), Type::Str) => Ok(Value::Str(format!("{i}"))),
            (Value::Str(i), Type::Str) => Ok(Value::Str(format!("{i}"))),
            (Value::Bool(i), Type::Str) => Ok(Value::Str(format!("{i}"))),

            (Value::Int(i), Type::Bool) => Ok(Value::Bool(*i != 0)),
            (Value::Dec(i), Type::Bool) => Ok(Value::Bool(*i != 0.0)),
            (Value::Str(i), Type::Bool) => Ok(Value::Bool(*i != "")),
            (Value::Bool(i), Type::Bool) => Ok(Value::Bool(*i)),

            (Value::List(l), Type::List(t)) => {
                let mut new_vals = Vec::new();

                for e in l {
                    new_vals.push(self.cast_to_type(e, t)?);
                }

                Ok(Value::List(new_vals))
            }
            (Value::Map(m), Type::Map(kt, vt)) => {
                let mut new_vals = HashMap::new();

                for (k, v) in m {
                    let k = self.cast_to_type(k, kt)?;
                    let v = self.cast_to_type(v, vt)?;
                    new_vals.insert(k, v);
                }

                Ok(Value::Map(new_vals))
            }

            (_, _) => Err(RuntimeError::CannotCastValueToType(
                value.clone(),
                t.clone(),
            )),
        }
    }
}
