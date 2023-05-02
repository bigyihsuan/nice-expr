use crate::eval::r#type::Type;

pub type Program = Vec<Expr>;

#[derive(Debug, Clone)]
pub enum Expr {
    Literal(Literal),
    Identifier(String),
    Declaration(Declaration),
    Assignment(Assignment),
    FunctionCall {
        name: String,
        args: Vec<Expr>,
    },
    FunctionDefinition {
        args: Vec<Declaration>,
        ret: Type,
        body: Program,
    },

    Minus(UnaryExpr),
    Not(UnaryExpr),
    Indexing(Indexing),
    Multiplication(BinaryExpr),
    Addition(BinaryExpr),
    Comparison(BinaryExpr),
    Logical(BinaryExpr),

    Block(Program),
    Return(Option<Box<Expr>>),
    If {
        condition: Box<Expr>,
        when_true: Box<Expr>,
        when_false: Option<Box<Expr>>,
    },
    For {
        vars: Vec<Declaration>,
        body: Program,
    },
    ForIn {
        vars: Vec<Declaration>,
        collection: Box<Expr>,
        body: Program,
    },
    Break(Option<Box<Expr>>),

    TypeName(Type),
    TypeCast(Box<Expr>, Type),
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
    Const(Decl),
    Var(Decl),
}

#[derive(Debug, Clone)]
pub struct Decl {
    pub name: String,
    pub type_name: Type,
    pub expr: Option<Box<Expr>>,
}

#[derive(Debug, Clone)]
pub struct Assignment {
    pub name: String,
    pub index: Option<Box<Expr>>,
    pub op: BinaryOperator,
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
    GetIndexing,
    SetIndexing,
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
    Is,
}

#[derive(Debug, Clone)]
pub struct Indexing {
    pub collection: Box<Expr>,
    pub op: BinaryOperator,
    pub index: IndexKind,
}

#[derive(Debug, Clone)]
pub enum IndexKind {
    Single { index: Box<Expr> },
    Range { start: Box<Expr>, end: Box<Expr> },
}
