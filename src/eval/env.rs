use std::{cell::RefCell, collections::HashMap, rc::Rc};

use crate::eval::intepreter::builtin;
use crate::prelude::RuntimeError;

use super::{
    r#type::Type,
    value::{Func, Value},
    Constness,
};

pub type SEnv = Rc<RefCell<Env>>;

#[derive(Debug, Clone)]
pub struct Env {
    parent: Option<SEnv>,
    values: HashMap<String, ValueEntry>,
}

impl Env {
    pub fn new() -> Self {
        Self {
            parent: None,
            values: HashMap::new(),
        }
    }

    pub fn extend(parent: SEnv) -> Rc<RefCell<Self>> {
        Rc::new(RefCell::new(Self {
            parent: Some(parent.clone()),
            values: HashMap::new(),
        }))
    }

    pub fn identifiers(&self) -> HashMap<String, ValueEntry> {
        if let Some(parent) = &self.parent {
            let parent_values = parent.borrow().identifiers();
            self.values
                .clone()
                .into_iter()
                .chain(parent_values.into_iter())
                .collect()
        } else {
            self.values.clone()
        }
    }

    pub fn get(&self, name: String) -> Result<ValueEntry, RuntimeError> {
        // try this Env first
        if let Some(value) = self.values.get(&name) {
            Ok(value.clone())
        } else if let Some(parent) = &self.parent {
            parent.borrow().get(name)
        } else {
            Err(RuntimeError::IdentifierNotFound(name.clone()))
        }
    }

    pub fn set(
        &mut self,
        name: String,
        value: Value,
        ignore_constness: bool,
    ) -> Result<Value, RuntimeError> {
        if self.values.contains_key(&name.clone()) {
            // disallow setting on const
            let ve = self.values.get(&name).unwrap();
            if let Constness::Const = ve.c && !ignore_constness {
                return Err(RuntimeError::SettingConst(name));
            }
            let t = value.to_type()?;
            if ve.t == t.clone() {
                let new_entry = ValueEntry {
                    v: value.clone(),
                    c: ve.c.clone(),
                    t,
                };
                self.values.insert(name, new_entry);
                Ok(value)
            } else {
                Err(RuntimeError::MismatchedTypes {
                    got: vec![t],
                    expected: vec![ve.t.clone()],
                })
            }
        } else if let Some(parent) = &self.parent {
            parent.borrow_mut().set(name, value, ignore_constness)
        } else {
            Err(RuntimeError::IdentifierNotFound(name))
        }
    }

    pub fn define(
        &mut self,
        name: String,
        value: Value,
        constness: Constness,
        ident_type: Type,
    ) -> Result<Value, String> {
        if !self.values.contains_key(&name) {
            self.values.insert(
                name,
                ValueEntry {
                    v: value.clone(),
                    c: constness,
                    t: ident_type,
                },
            );
            Ok(value)
        } else {
            Err(name)
        }
    }
    pub fn def_var(
        &mut self,
        name: String,
        value: Value,
        type_name: Type,
    ) -> Result<Value, String> {
        self.define(name, value, Constness::Var, type_name)
    }
    pub fn def_const(
        &mut self,
        name: String,
        value: Value,
        type_name: Type,
    ) -> Result<Value, String> {
        self.define(name, value, Constness::Const, type_name)
    }

    pub fn undefine(&mut self, name: String) {
        if self.values.contains_key(&name) {
            self.values.remove(&name);
        } else if let Some(parent) = &self.parent {
            parent.borrow_mut().undefine(name);
        }
    }
}

impl Default for Env {
    fn default() -> Self {
        let mut env = Self::new();
        let _ = env.define(
            "print".into(),
            Value::Func(Func::Native(builtin::print)),
            Constness::Const,
            Type::Func(vec![Type::BuiltinVariadic], Box::new(Type::None)),
        );
        let _ = env.define(
            "printline".into(),
            Value::Func(Func::Native(builtin::println)),
            Constness::Const,
            Type::Func(vec![Type::BuiltinVariadic], Box::new(Type::None)),
        );
        let _ = env.define(
            "len".into(),
            Value::Func(Func::Native(builtin::len)),
            Constness::Const,
            Type::Func(vec![Type::BuiltinVariadic], Box::new(Type::Int)),
        );
        let _ = env.define(
            "range".into(),
            Value::Func(Func::Native(builtin::range)),
            Constness::Const,
            Type::Func(
                vec![Type::Int, Type::Int, Type::Int],
                Box::new(Type::List(Box::new(Type::Int))),
            ),
        );
        let _ = env.define(
            "inputchar".into(),
            Value::Func(Func::Native(builtin::inputchar)),
            Constness::Const,
            Type::Func(vec![], Box::new(Type::Str)),
        );
        let _ = env.define(
            "inputline".into(),
            Value::Func(Func::Native(builtin::inline)),
            Constness::Const,
            Type::Func(vec![], Box::new(Type::Str)),
        );
        let _ = env.define(
            "inputall".into(),
            Value::Func(Func::Native(builtin::inall)),
            Constness::Const,
            Type::Func(vec![], Box::new(Type::Str)),
        );
        env
    }
}

#[derive(Debug, Clone)]
pub struct ValueEntry {
    pub v: Value,
    pub c: Constness,
    pub t: Type,
}
