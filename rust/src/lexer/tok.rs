use std::fmt::Display;

use logos::Logos;
use snailquote::unescape;

#[derive(Logos, Debug, Clone, PartialEq)]
pub enum Token {
    // ignore whitespace
    #[error]
    #[regex(r"[ \t\r\n\f]+", logos::skip)]
    Error,

    // comments are double-slashes up to newline
    #[regex(r"//[^\n]*\n?", logos::skip)]
    Comment,

    // simple literals
    #[regex("[0-9]+", |lex| parse_int::parse::<i64>(lex.slice()))]
    IntLit(i64),
    #[regex("[0-9]+.[0-9]+", |lex| parse_int::parse::<f64>(lex.slice()))]
    DecLit(f64),
    // strings are anything that's not escaped quote
    #[regex("\"(?:[^\"]|\\\\\")*\"", |lex| unescape(lex.slice()))]
    StrLit(String),
    #[token("true", |lex| lex.slice().parse::<bool>().unwrap())]
    TrueBoolLit(bool),
    #[token("false", |lex| lex.slice().parse::<bool>().unwrap())]
    FalseBoolLit(bool),

    // symbols
    #[token("[")]
    LeftBracket,
    #[token("]")]
    RightBracket,
    #[token("{")]
    LeftBrace,
    #[token("}")]
    RightBrace,
    #[token("(")]
    LeftParen,
    #[token(")")]
    RightParen,
    #[token("<|")]
    LeftTriangle,
    #[token("|>")]
    RightTriangle,
    #[token("+")]
    Plus,
    #[token("-")]
    Minus,
    #[token("*")]
    Star,
    #[token("/")]
    Slash,
    #[token("&")]
    Percent,
    #[token("+=")]
    PlusEqual,
    #[token("-=")]
    MinusEqual,
    #[token("*=")]
    StarEqual,
    #[token("/=")]
    SlashEqual,
    #[token("%=")]
    PercentEqual,
    #[token("=")]
    Equal,
    #[token(">")]
    Greater,
    #[token(">=")]
    GreaterEqual,
    #[token("<")]
    Less,
    #[token("<=")]
    LessEqual,
    #[token(",")]
    Comma,
    #[token(":")]
    Colon,
    #[token(";")]
    Semicolon,
    #[token("_")]
    Underscore,
    // keywords
    #[token("and")]
    And,
    #[token("or")]
    Or,
    #[token("not")]
    Not,
    #[token("var")]
    Var,
    #[token("const")]
    Const,
    #[token("set")]
    Set,
    #[token("is")]
    Is,
    #[token("for")]
    For,
    #[token("break")]
    Break,
    #[token("return")]
    Return,
    #[token("func")]
    Func,
    #[token("if")]
    If,
    #[token("then")]
    Then,
    #[token("else")]
    Else,
    // type keywords
    #[token("int")]
    IntTypename,
    #[token("dec")]
    DecTypename,
    #[token("str")]
    StrTypename,
    #[token("bool")]
    BoolTypename,
    #[token("list")]
    ListTypename,
    #[token("map")]
    MapTypename,

    // identifiers can only have letters
    #[regex("[a-zA-Z]+", |lex| lex.slice().to_string())]
    Ident(String),
}

impl Display for Token {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        f.write_fmt(format_args!("{:?}", self))
    }
}
