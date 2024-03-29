use crate::{
    eval::r#type::Type,
    lexer::{tok::Token, TokenStream},
    parse::ast::{
        Assignment, BinaryExpr, BinaryOperator, Decl, Declaration, Expr, IndexKind, Indexing,
        Literal, Program, UnaryExpr, UnaryOperator,
    },
};

peg::parser! {
    pub grammar parser<'source>() for TokenStream<'source>  {
        pub rule program() -> Program
        = stmt()+

        pub rule stmt() -> Expr
        = e:expr() [Token::Semicolon]
        { e }

        pub rule expr() -> Expr = precedence!{
            expr:assignment() { expr }
            expr:declaration_expr() { expr }
            expr:function_call() { expr }
            --
            expr:if_expr() { expr }
            expr:for_expr() { expr }
            expr:for_while_expr() { expr }
            expr:for_in_expr() { expr }
            expr:function_definition() { expr }
            --
            [Token::Return] expr:(@) { Expr::Return(Some(Box::new(expr))) }
            [Token::Return] { Expr::Return(None) }
            --
            [Token::Break] expr:(@) { Expr::Break(Some(Box::new(expr))) }
            [Token::Break] { Expr::Break(None) }
            --
            expr:(@) [Token::As] t:type_name() { Expr::TypeCast(Box::new(expr), t) }
            --
            [Token::Not] expr:(@) {
                Expr::Not(UnaryExpr{op: UnaryOperator::Not, expr: Box::new(expr)})
            }
            --
            left:(@) op:[Token::And | Token::Or] right:@ {
                Expr::Logical(BinaryExpr { left: Box::new(left), op: match op {
                Token::And => BinaryOperator::And,
                Token::Or => BinaryOperator::Or,
                _ => unreachable!(),
            }, right: Box::new(right) })}
            --
            left:(@) op:[Token::Greater | Token::Less | Token::GreaterEqual | Token::LessEqual | Token::Equal | Token::NotEqual] right:@ {
                Expr::Comparison(BinaryExpr{ left: Box::new(left), op: match op {
                    Token::Greater => BinaryOperator::Greater,
                    Token::Less => BinaryOperator::Less,
                    Token::GreaterEqual => BinaryOperator::GreaterEqual,
                    Token::LessEqual => BinaryOperator::LessEqual,
                    Token::Equal => BinaryOperator::Equal,
                    Token::NotEqual => BinaryOperator::NotEqual,
                    _ => unreachable!(),
                }, right: Box::new(right) })
            }
            --
            left:(@) op:[Token::Plus | Token::Minus] right:@ {
                Expr::Addition(BinaryExpr{ left: Box::new(left), op: match op {
                    Token::Plus => BinaryOperator::Add,
                    Token::Minus => BinaryOperator::Subtract,
                    _ => unreachable!(),
                }, right: Box::new(right) })
            }
            --
            left:(@) op:[Token::Star | Token::Slash | Token::Percent] right:@ {
                Expr::Multiplication(BinaryExpr{ left: Box::new(left), op: match op {
                    Token::Star => BinaryOperator::Times,
                    Token::Slash => BinaryOperator::Divide,
                    Token::Percent => BinaryOperator::Modulo,
                    _ => unreachable!(),
                }, right: Box::new(right) })
            }
            --
            [Token::Minus] expr:(@) {
                Expr::Minus(UnaryExpr{op: UnaryOperator::Minus, expr: Box::new(expr)})
            }
            --
            left:@ op:[Token::Underscore] right:index() {
                Expr::Indexing(Indexing{collection: Box::new(left), op: match op {
                    Token::Underscore => BinaryOperator::GetIndexing,
                    _ => unreachable!()
                }, index: right})
            }
            --
            expr:block_expr() { expr }
            --
            expr:type_name_expr() { expr }
            expr:expr_lit_or_ident() { expr }
            --
            [Token::LeftParen] expr:expr() [Token::RightParen] { expr }
        }

        pub rule index() -> IndexKind
        = start:expr() end:([Token::DoubleDot] end:expr() {end})?
        {
            if let Some(end) = end {
                IndexKind::Range{start: Box::new(start), end: Box::new(end)}
            } else {
                IndexKind::Single{index: Box::new(start)}
            }
        }

        pub rule function_definition() -> Expr
        = [Token::Func] [Token::LeftParen] args:(declaration() ** [Token::Comma]) [Token::Comma]? [Token::RightParen] ret:type_name() body:block()
        { Expr::FunctionDefinition {args, ret, body} }

        pub rule if_expr() -> Expr
        = [Token::If] vars:(vars:decl_list() [Token::While] {vars})? condition:expr() [Token::Then] when_true:block_expr() when_false:([Token::Else] when_false:(block_expr() / if_expr()) {Box::new(when_false)})?
        {
            if let Some(vars) = vars {
                Expr::If { vars, condition: Box::new(condition), when_true: Box::new(when_true), when_false: when_false }
            } else {
                Expr::If { vars: vec![], condition: Box::new(condition), when_true: Box::new(when_true), when_false: when_false }
            }
        }

        pub rule for_expr() -> Expr
        = [Token::For] vars:decl_list() body:block() { Expr::For{ vars, body } }

        pub rule for_while_expr() -> Expr
        = [Token::For] vars:decl_list() [Token::While] condition:expr() body:block()
        { Expr::ForWhile{ vars, condition: Box::new(condition), body } }

        pub rule for_in_expr() -> Expr
        = [Token::For] vars:decl_list() [Token::In] collection:expr() body:block()
        { Expr::ForIn{ vars, collection: Box::new(collection), body } }

        rule decl_list() -> Vec<Declaration>
        = vars:(declaration() ** [Token::Comma]) [Token::Comma]? { vars }

        pub rule block_expr() -> Expr
        = block:block() { Expr::Block(block) }

        pub rule block() -> Program
        = [Token::LeftBrace] program:program() [Token::RightBrace] { program }

        pub rule declaration_expr() -> Expr
        = d:declaration() { Expr::Declaration(d) }
        pub rule declaration() -> Declaration
        = declaration_var() / declaration_const()
        pub rule declaration_var() -> Declaration
        = [Token::Var] name:identifier() [Token::Is] type_name:type_name() value:expr()?
        { Declaration::Var(Decl{ name, type_name, expr: value.map(|e| Box::new(e)) })}
        pub rule declaration_const() -> Declaration
        = [Token::Const] name:identifier() [Token::Is] type_name:type_name() value:expr()?
        { Declaration::Const(Decl{ name, type_name, expr: value.map(|e| Box::new(e)) })}

        pub rule assignment() -> Expr
        = [Token::Set] name:identifier() index:([Token::At] index:expr() {Box::new(index)})? op:assignment_operator() value:expr()
        {Expr::Assignment(Assignment { name, index, op, expr: Box::new(value) })}
        pub rule assignment_operator() -> BinaryOperator
        = op:[Token::Is | Token::And | Token::Or | Token::Greater | Token::Less | Token::GreaterEqual | Token::LessEqual
        | Token::Equal | Token::NotEqual | Token::Plus | Token::Minus | Token::Star | Token::Slash | Token::Percent]
        { match op {
            Token::Is => BinaryOperator::Is,
            Token::And => BinaryOperator::And,
            Token::Or => BinaryOperator::Or,
            Token::Greater => BinaryOperator::Greater,
            Token::Less => BinaryOperator::Less,
            Token::GreaterEqual => BinaryOperator::GreaterEqual,
            Token::LessEqual => BinaryOperator::LessEqual,
            Token::Equal => BinaryOperator::Equal,
            Token::NotEqual => BinaryOperator::NotEqual,
            Token::Plus => BinaryOperator::Add,
            Token::Minus => BinaryOperator::Subtract,
            Token::Star => BinaryOperator::Times,
            Token::Slash => BinaryOperator::Divide,
            Token::Percent => BinaryOperator::Modulo,
            _ => unreachable!()
        } }

        pub rule function_call() -> Expr
        = name:identifier() [Token::LeftParen] args:(expr() ** [Token::Comma]) [Token::Comma]? [Token::RightParen]
        { Expr::FunctionCall { name, args } }

        pub rule expr_lit_or_ident() -> Expr
        = l:literal() {l}
        / i:expr_identifier() {i}

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
        = [Token::LeftBracket] l:(expr() ** [Token::Comma]) [Token::Comma]? [Token::RightBracket]
        { Literal::List(l) }
        pub rule literal_map() -> Literal
        = [Token::LeftTriangle] m:(map_element() ** [Token::Comma]) [Token::Comma]? [Token::RightTriangle]
        { let m = m.into_iter().collect(); Literal::Map(m) }

        pub rule map_element() -> (Expr, Expr)
        =  l:expr() [Token::Colon] r:expr()
        { (l,r) }

        pub rule type_name_expr() -> Expr
        = expr:type_name() { Expr::TypeName(expr) }
        pub rule type_name() -> Type
        = simple_type() / compound_type()

        pub rule simple_type() -> Type
        = [Token::NoneTypename] {Type::None}
        / [Token::IntTypename] {Type::Int}
        / [Token::DecTypename] {Type::Dec}
        / [Token::StrTypename] {Type::Str}
        / [Token::BoolTypename] {Type::Bool}

        pub rule compound_type() -> Type
        = [Token::ListTypename] [Token::LeftBracket] t:type_name() [Token::RightBracket]
            {Type::List(Box::new(t))}
        / [Token::MapTypename] [Token::LeftBracket] k:type_name() [Token::RightBracket] v:type_name()
            {Type::Map(Box::new(k), Box::new(v))}
        / [Token::Func] [Token::LeftParen] args:(type_name() ** [Token::Comma]) [Token::Comma]? [Token::RightParen] ret:type_name()
            {Type::Func(args, Box::new(ret))}
    }
}
