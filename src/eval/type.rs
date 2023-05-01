use std::fmt::Display;

use itertools::Itertools;

#[derive(Debug, Clone, Hash)]
pub enum Type {
    None,
    BuiltinVariadic,
    Int,
    Dec,
    Str,
    Bool,
    List(Box<Type>),
    Map(Box<Type>, Box<Type>),
    Func(Vec<Type>, Box<Type>),
    Break(Box<Type>),
    Any,
}

impl Type {
    pub fn key_type(&self) -> Option<Type> {
        match self {
            Type::None => None,
            Type::BuiltinVariadic => None,
            Type::Int => None,
            Type::Dec => None,
            Type::Str => Some(Type::Int),
            Type::Bool => None,
            Type::List(_) => Some(Type::Int),
            Type::Map(box k, _) => Some(k.clone()),
            Type::Break(box t) => t.key_type(),
            Type::Func(_, _) => None,
            Type::Any => None,
        }
    }

    pub fn element_type(&self) -> Option<Type> {
        match self {
            Type::None => None,
            Type::BuiltinVariadic => None,
            Type::Int => None,
            Type::Dec => None,
            Type::Str => Some(Type::Str),
            Type::Bool => None,
            Type::List(box e) => Some(e.clone()),
            Type::Map(_, box v) => Some(v.clone()),
            Type::Break(box t) => t.element_type(),
            Type::Func(_, _) => None,
            Type::Any => None,
        }
    }

    pub fn infer_contained_type(&self, other: &Self) -> Option<Self> {
        match (self, other) {
            (Type::List(box l), Type::List(box r)) if *l != Type::None && *r == Type::None => {
                Some(Type::List(Box::new(l.clone())))
            }

            (Type::List(box l), Type::List(box r)) if *l == Type::None && *r != Type::None => {
                Some(Type::List(Box::new(r.clone())))
            }
            (Type::Map(box lk, box lv), Type::Map(box rk, box rv)) => {
                match (
                    *lk != Type::None,
                    *lv != Type::None,
                    *rk != Type::None,
                    *rv != Type::None,
                ) {
                    (true, true, true, true) => Some(self.clone()),
                    (true, true, true, false) => Some(self.clone()),
                    (true, true, false, true) => Some(self.clone()),
                    (true, true, false, false) => Some(self.clone()),
                    (true, false, true, true) => Some(other.clone()),
                    (false, true, true, true) => Some(other.clone()),
                    (false, false, true, true) => Some(other.clone()),
                    (true, false, false, true) => {
                        Some(Type::Map(Box::new(lk.clone()), Box::new(rv.clone())))
                    }
                    (false, true, true, false) => {
                        Some(Type::Map(Box::new(rk.clone()), Box::new(lv.clone())))
                    }
                    _ => None,
                }
            }
            (Type::Break(box l), r) => l.infer_contained_type(r),
            (l, Type::Break(box r)) => l.infer_contained_type(r),
            // TODO: Func handled by below?
            (l, _) if *l != Type::None => Some(l.clone()),
            (_, r) if *r != Type::None => Some(r.clone()),
            _ => None,
        }
    }
}

impl PartialEq for Type {
    fn eq(&self, other: &Self) -> bool {
        match (self, other) {
            (Self::List(l0), Self::List(r0)) => l0 == r0,
            (Self::Map(l0, l1), Self::Map(r0, r1)) => l0 == r0 && l1 == r1,
            (Self::Func(la, lr), Self::Func(ra, rr)) => la == ra && lr == rr,
            _ => core::mem::discriminant(self) == core::mem::discriminant(other),
        }
    }
}
impl Eq for Type {}

impl Display for Type {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Type::None => f.write_str("none"),
            Type::BuiltinVariadic => f.write_str("builtin_variadic"),
            Type::Int => f.write_str("int"),
            Type::Dec => f.write_str("dec"),
            Type::Str => f.write_str("str"),
            Type::Bool => f.write_str("bool"),
            Type::List(box e) => f.write_fmt(format_args!("list[{e}]")),
            Type::Map(box k, box v) => f.write_fmt(format_args!("map[{k}]{v}")),
            Type::Func(args, ret) => f.write_fmt(format_args!(
                "func({}){ret}",
                args.into_iter()
                    .map(|t| format!("{t}"))
                    .collect_vec()
                    .join(",")
            )),
            Type::Break(box t) => f.write_fmt(format_args!("{t}")),
            Type::Any => f.write_str("any"),
        }
    }
}
