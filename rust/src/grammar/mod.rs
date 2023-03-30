use crate::{
    lexer::{tok::Token, TokenStream},
    parse::parse_tree::{Expr, Literal, Program},
};

peg::parser! {
    pub grammar module_parser<'source>() for TokenStream<'source>  {
        pub rule program() -> Program
        = expr()*

        pub rule expr() -> Expr
        = lit:literal()
        { Expr::Literal(lit) }

        pub rule literal() -> Literal
        = literal_int()
        / literal_dec()
        / literal_str()
        / literal_bool()
        // / literal_list()
        // / literal_map()

        pub rule literal_int() -> Literal
        = [Token::IntLit(i)]
        { Literal::Int(*i) }
        pub rule literal_dec() -> Literal
        = [Token::DecLit(i)]
        { Literal::Dec(*i) }
        pub rule literal_str() -> Literal
        = [Token::StrLit(i)]
        { Literal::Str(i.clone()) }
        pub rule literal_bool() -> Literal
        = [Token::TrueBoolLit(i) | Token::FalseBoolLit(i)]
        { Literal::Bool(*i) }
        // pub rule literal_list() -> Literal
        // = [Token::IntLit(i)]
        // { Literal::Int(i) }
        // pub rule literal_map() -> Literal
        // = [Token::IntLit(i)]
        // { Literal::Int(i) }
    }
}
