use std::fmt::Display;

pub type Result<T> = std::result::Result<T, Vec<SyntaxError>>;

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
