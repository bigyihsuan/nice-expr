use std::{cell::RefCell, collections::HashMap, rc::Rc};

use crate::prelude::RuntimeError;

use super::{r#type::Type, value::Value, Constness};

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

    pub fn make_child(parent: SEnv) -> Self {
        Self {
            parent: Some(parent),
            values: HashMap::new(),
        }
    }

    pub fn identifiers(&self) -> HashMap<String, ValueEntry> {
        self.values.clone()
    }

    pub fn get(&self, name: String) -> Option<ValueEntry> {
        // try this Env first
        if let Some(value) = self.values.get(&name) {
            Some(value.clone())
        } else if let Some(parent) = &self.parent {
            parent.borrow().get(name)
        } else {
            None
        }
    }

    pub fn set(&mut self, name: String, value: Value) -> Result<Value, RuntimeError> {
        if self.values.contains_key(&name.clone()) {
            // disallow setting on const
            let ve = self.values.get(&name).unwrap();
            if let Constness::Const = ve.c {
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
            parent.borrow_mut().set(name, value)
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
}

impl Default for Env {
    fn default() -> Self {
        let env = Self::new();
        // TODO: define builtin functions
        // TODO: env.define("print", value, Constness::Const)
        // TODO: env.define("println", value, Constness::Const)
        env
    }
}

#[derive(Debug, Clone)]
pub struct ValueEntry {
    pub v: Value,
    pub c: Constness,
    pub t: Type,
}
