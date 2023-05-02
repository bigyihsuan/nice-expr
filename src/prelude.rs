use std::{fmt::Display, io, string};

use itertools::Itertools;

use crate::{
    eval::{r#type::Type, value::Value},
    lexer::TokenLocation,
    parse::ast::{BinaryOperator, Operator},
};

pub type Result<T> = std::result::Result<T, SyntaxError>;

#[derive(Debug)]
pub enum SyntaxError {
    InvalidToken {
        token: String,
        filename: Option<std::path::PathBuf>,
        line: usize,
        col: usize,
    },
    UnexpectedToken {
        token: String,
        filename: Option<std::path::PathBuf>,
        line: usize,
        col: usize,
        expected: Vec<&'static str>,
    },
}

impl Display for SyntaxError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Self::InvalidToken {
                token,
                filename,
                line,
                col,
            } => {
                let filename = filename
                    .as_ref()
                    .map(|path| path.as_path().display().to_string())
                    .unwrap_or("<>".to_string());

                write!(f, "invalid token '{token}' @ {filename}:{line}:{col}")
            }
            Self::UnexpectedToken {
                token,
                filename,
                line,
                col,
                expected,
            } => {
                let filename = filename
                    .as_ref()
                    .map(|path| path.as_path().display().to_string())
                    .unwrap_or("<>".to_string());

                write!(
                    f,
                    "unexpected token '{token}' @ {filename}:{line}:{col}\n want {expected:?}"
                )
            }
        }
    }
}

impl std::error::Error for SyntaxError {
    fn source(&self) -> Option<&(dyn std::error::Error + 'static)> {
        None
    }
}

type ParseError = peg::error::ParseError<TokenLocation>;

impl From<ParseError> for SyntaxError {
    fn from(err: ParseError) -> Self {
        let TokenLocation {
            filename,
            linecol: (line, col),
            token,
        } = err.location;
        let expected = err.expected.tokens().collect_vec();

        Self::UnexpectedToken {
            token: token.map(|tok| tok.to_string()).unwrap_or("<>".to_string()),
            filename,
            line,
            col,
            expected,
        }
    }
}

#[derive(Debug, Clone)]
pub enum RuntimeError {
    NotImplemented,
    MismatchedTypes { got: Vec<Type>, expected: Vec<Type> },
    InvalidOperatorOnTypes { op: Operator, types: Vec<Type> },
    IdentifierNotFound(String),
    SettingConst(String),
    DivideByZero,
    InvalidAssignmentOperator(BinaryOperator),
    NotEnoughArguments { want: usize, got: usize },
    IndexingNonIndexable { got: Type },
    TakingLenOfLengthless { got: Type },
    IOError(IOError),
    InvalidRangeStep(isize),
    IndexingCollectionWithZeroElements(Value),
    IndexOutOfBounds(Value, Value, Option<Value>),
    KeyNotFound(Value, Value),
    CallingNonCallable { name: String },
    CannotCastValueToType(Value, Type),
    StringToValueParseFailed(Value, Type),
    ForInFirstDeclMustBeVar,
}

#[derive(Debug, Clone)]
pub enum IOError {
    ErrorKind(io::ErrorKind),
    UtfError(string::FromUtf8Error),
    CouldNotGetChar,
}
