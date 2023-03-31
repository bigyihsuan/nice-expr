use crate::{
    lexer::{tok::Token, TokenStream},
    parse::ast::{Expr, Literal, Program},
};

peg::parser! {
    pub grammar module_parser<'source>() for TokenStream<'source>  {
        pub rule program() -> Program
        = expr()+

        pub rule expr() -> Expr
        = literal() / identifier() / unary_expr()

        pub rule unary_expr() -> Expr
        = op:[Token::Not | Token::Minus] expr:expr()
        { Expr::Unary{op: op.clone(), expr: Box::new(expr)}}


        pub rule identifier() -> Expr
        = [Token::Ident(name)]
        {Expr::Identifier{name: name.clone()}}

        pub rule literal() -> Expr
        = l:(literal_int()
        / literal_dec()
        / literal_str()
        / literal_bool()
        / literal_list()
        / literal_map())
        { Expr::Literal(l) }

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
        pub rule literal_list() -> Literal
        = [Token::LeftBracket] l:(literal() ** [Token::Comma]) [Token::Comma]? [Token::RightBracket]
        { Literal::List(l) }
        pub rule literal_map() -> Literal
        = [Token::LeftTriangle] m:(literal_map_element() ** [Token::Comma]) [Token::Comma]? [Token::RightTriangle]
        { let m = m.into_iter().collect(); Literal::Map(m) }

        pub rule literal_map_element() -> (Expr, Expr)
        =  l:literal() [Token::Colon] r:literal()
        { (l,r) }
    }
}
