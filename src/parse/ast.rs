use crate::eval::r#type::Type;

pub type Program = Vec<Expr>;

#[derive(Debug, Clone)]
pub enum Expr {
    Literal(Literal),
    Identifier(String),
    Declaration(Declaration),
    Assignment(Assignment),
    FunctionCall { name: String, args: Vec<Expr> },

    Minus(UnaryExpr),
    Not(UnaryExpr),
    Indexing(BinaryExpr),
    Multiplication(BinaryExpr),
    Addition(BinaryExpr),
    Comparison(BinaryExpr),
    Logical(BinaryExpr),
}

#[derive(Debug, Clone)]
pub struct UnaryExpr {
    pub op: UnaryOperator,
    pub expr: Box<Expr>,
}

#[derive(Debug, Clone)]
pub struct BinaryExpr {
    pub left: Box<Expr>,
    pub op: BinaryOperator,
    pub right: Box<Expr>,
}

#[derive(Debug, Clone)]
pub enum Declaration {
    Const {
        name: String,
        type_name: Type,
        expr: Option<Box<Expr>>,
    },
    Var {
        name: String,
        type_name: Type,
        expr: Option<Box<Expr>>,
    },
}

#[derive(Debug, Clone)]
pub struct Assignment {
    pub name: String,
    pub op: AssignmentOperator,
    pub expr: Box<Expr>,
}

#[derive(Debug, Clone)]
pub enum Literal {
    Int(i64),
    Dec(f64),
    Str(String),
    Bool(bool),
    List(Vec<Expr>),
    Map(Vec<(Expr, Expr)>),
}

#[derive(Debug, Clone)]
pub enum Operator {
    UnaryOperator(UnaryOperator),
    BinaryOperator(BinaryOperator),
}

#[derive(Debug, Clone)]
pub enum UnaryOperator {
    Minus,
    Not,
}

#[derive(Debug, Clone)]
pub enum BinaryOperator {
    Indexing,
    Times,
    Divide,
    Modulo,
    Add,
    Subtract,
    Greater,
    Less,
    GreaterEqual,
    LessEqual,
    Equal,
    NotEqual,
    And,
    Or,
}

#[derive(Debug, Clone)]
pub enum AssignmentOperator {
    Invalid,
    Is,
    // TODO: the other assignment operators
}
