use crate::{
    lexer::{tok::Token, TokenStream},
    parse::ast::{
        Assignment, AssignmentOperator, BinaryExpr, BinaryOperator, Declaration, Expr, Literal,
        Program, Type, UnaryExpr, UnaryOperator,
    },
};

peg::parser! {
    pub grammar module_parser<'source>() for TokenStream<'source>  {
        pub rule program() -> Program
        = stmt()+

        pub rule stmt() -> Expr
        = e:expr() [Token::Semicolon]
        {e}

        pub rule expr() -> Expr = precedence!{
            expr:function_call() { expr }
            --
            expr:assignment() { expr }
            --
            expr:declaration() { expr }
            --
            expr:literal() { expr }
            --
            expr:expr_identifier() { expr }
            --
            [Token::Not] expr:(@) {
                Expr::Not(UnaryExpr{op: UnaryOperator::Not, expr: Box::new(expr)})
            }
            --
            left:(@) op:[Token::And | Token::Or] right:@ {
                Expr::Logical(BinaryExpr { left: Box::new(left), op: match op {
                Token::And => BinaryOperator::And,
                Token::Or => BinaryOperator::Or,
                _ => todo!(),
            }, right: Box::new(right) })}
            --
            left:(@) op:[Token::Greater | Token::Less | Token::GreaterEqual | Token::LessEqual | Token::Equal] right:@ {
                Expr::Comparison(BinaryExpr{ left: Box::new(left), op: match op {
                    Token::Greater => BinaryOperator::Greater,
                    Token::Less => BinaryOperator::Less,
                    Token::GreaterEqual => BinaryOperator::GreaterEqual,
                    Token::LessEqual => BinaryOperator::LessEqual,
                    Token::Equal => BinaryOperator::Equal,
                    _ => todo!(),
                }, right: Box::new(right) })
            }
            --
            left:(@) op:[Token::Plus | Token::Minus ] right:@ {
                Expr::Addition(BinaryExpr{ left: Box::new(left), op: match op {
                    Token::Plus => BinaryOperator::Add,
                    Token::Minus => BinaryOperator::Subtract,
                    _ => todo!(),
                }, right: Box::new(right) })
            }
            --
            left:(@) op:[Token::Star | Token::Slash | Token::Percent] right:@ {
                Expr::Multiplication(BinaryExpr{ left: Box::new(left), op: match op {
                    Token::Star => BinaryOperator::Times,
                    Token::Slash => BinaryOperator::Divide,
                    Token::Percent => BinaryOperator::Modulo,
                    _ => todo!(),
                }, right: Box::new(right) })
            }
            --
            [Token::Minus] expr:(@) {
                Expr::Minus(UnaryExpr{op: UnaryOperator::Minus, expr: Box::new(expr)})
            }
            --
            left:(@) op:[Token::Underscore] right:@ {
                Expr::Indexing(BinaryExpr{left: Box::new(left), op: match op {
                    Token::Underscore => BinaryOperator::Indexing,
                    _ => todo!()
                }, right: Box::new(right)})
            }
            --
            [Token::LeftParen] expr:expr() [Token::RightParen] { expr }
        }

        pub rule declaration() -> Expr
        = declaration_var() / declaration_const()
        pub rule declaration_var() -> Expr
        = [Token::Var] name:identifier() [Token::Is] type_name:type_name() value:expr()
        { Expr::Declaration(Declaration::Var { name, type_name, expr: Box::new(value) })}
        pub rule declaration_const() -> Expr
        = [Token::Const] name:identifier() [Token::Is] type_name:type_name() value:expr()
        { Expr::Declaration(Declaration::Const { name, type_name, expr: Box::new(value) })}

        pub rule assignment() -> Expr
        = [Token::Set] name:identifier() op:assignment_operator() value:expr()
        {Expr::Assignment(Assignment { name, op, expr: Box::new(value) })}
        pub rule assignment_operator() -> AssignmentOperator
        = op:[Token::Is]
        { match op {
            Token::Is => AssignmentOperator::Is,
            _ => AssignmentOperator::Invalid
        } }

        pub rule function_call() -> Expr
        = name:identifier() [Token::LeftParen] args:(expr() ** [Token::Comma]) [Token::Comma]? [Token::RightParen]
        { Expr::FunctionCall { name, args } }

        pub rule expr_identifier() -> Expr
        = name:identifier()
        { Expr::Identifier(name) }
        pub rule identifier() -> String
        = [Token::Ident(name)]
        { name.clone() }

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

        pub rule type_name() -> Type
        = simple_type() / compound_type()

        pub rule simple_type() -> Type
        = i:[Token::IntTypename] {Type::Int}
        / i:[Token::DecTypename] {Type::Dec}
        / i:[Token::StrTypename] {Type::Str}
        / i:[Token::BoolTypename] {Type::Bool}

        pub rule compound_type() -> Type
        = [Token::ListTypename] [Token::LeftBracket] t:type_name() [Token::RightBracket] {Type::List(Box::new(t))}
        / [Token::MapTypename] [Token::LeftBracket] k:type_name() [Token::RightBracket] v:type_name() {Type::Map(Box::new(k), Box::new(v))}
    }
}
