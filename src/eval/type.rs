#[derive(Debug, Clone, Hash)]
pub enum Type {
    None,
    Int,
    Dec,
    Str,
    Bool,
    List(Box<Type>),
    Map(Box<Type>, Box<Type>),
    Break(Box<Type>),
}

impl Type {
    pub fn element_type(&self) -> Option<Type> {
        match self {
            Type::None => None,
            Type::Int => None,
            Type::Dec => None,
            Type::Str => Some(Type::Str),
            Type::Bool => None,
            Type::List(box e) => Some(e.clone()),
            Type::Map(_, box v) => Some(v.clone()),
            Type::Break(box t) => t.element_type(),
        }
    }

    pub fn key_type(&self) -> Option<Type> {
        match self {
            Type::None => None,
            Type::Int => None,
            Type::Dec => None,
            Type::Str => Some(Type::Int),
            Type::Bool => None,
            Type::List(_) => Some(Type::Int),
            Type::Map(box k, _) => Some(k.clone()),
            Type::Break(box t) => t.key_type(),
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
            _ => core::mem::discriminant(self) == core::mem::discriminant(other),
        }
    }
}
impl Eq for Type {}
