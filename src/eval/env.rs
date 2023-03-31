use std::{cell::RefCell, collections::HashMap, hash::Hash, rc::Rc};

use super::{value::Value, Constness, RuntimeError};

pub struct Env {
    parent: Option<Rc<RefCell<Env>>>,
    values: HashMap<String, (Value, Constness)>,
}

impl Env {
    pub fn new() -> Self {
        Self {
            parent: None,
            values: HashMap::new(),
        }
    }

    pub fn make_child(parent: Rc<RefCell<Env>>) -> Self {
        Self {
            parent: Some(parent),
            values: HashMap::new(),
        }
    }

    pub fn get(&self, name: String) -> Option<(Value, Constness)> {
        // try this Env first
        if let Some(value) = self.values.get(&name) {
            Some(value.clone())
        } else if let Some(parent) = &self.parent {
            parent.borrow().get(name)
        } else {
            None
        }
    }

    pub fn set(&mut self, name: String, value: Value) -> Result<(), (String, RuntimeError)> {
        if self.values.contains_key(&name) {
            // disallow setting on const
            let (_, constness) = self.values.get(&name).unwrap();
            if let Constness::Const = constness {
                return Err((name, RuntimeError::SettingConst));
            }
            self.values.insert(name, (value, constness.clone()));
            Ok(())
        } else if let Some(parent) = &self.parent {
            parent.borrow_mut().set(name, value)
        } else {
            Err((name, RuntimeError::IdentifierNotFound))
        }
    }

    pub fn define(
        &mut self,
        name: String,
        value: Value,
        constness: Constness,
    ) -> Result<(), String> {
        if !self.values.contains_key(&name) {
            self.values.insert(name, (value, constness));
            Ok(())
        } else {
            Err(name)
        }
    }
    pub fn def_var(&mut self, name: String, value: Value) -> Result<(), String> {
        self.define(name, value, Constness::Var)
    }
    pub fn def_const(&mut self, name: String, value: Value) -> Result<(), String> {
        self.define(name, value, Constness::Const)
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
