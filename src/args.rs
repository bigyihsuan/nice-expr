use std::{
    fs,
    io::{self, Read},
    path::PathBuf,
};

use clap::{ArgGroup, Parser};

#[derive(Debug, Parser)]
#[command(name = "nice-expr")]
#[command(author = "bigyhsuan")]
#[command(version = "0.0.0")]
#[clap(group(ArgGroup::new("source").required(false).multiple(false).args(&["file", "code"])))]
struct Args {
    /// Exclusive with -c. Take input source code from a file.
    /// If both -f and -c are missing, take code from STDIN.
    #[arg(short, long, value_name = "FILE")]
    file: Option<PathBuf>,
    /// Exclusive with -f. Take input source code from the next command-line argument.
    /// If both -f and -c are missing, take code from STDIN.
    #[arg(short, long, value_name = "CODE")]
    code: Option<String>,
    /// Output debug information.
    #[arg(short, long)]
    debug: bool,
}

pub fn parse_args() -> (Option<std::path::PathBuf>, String, bool) {
    let args = Args::parse();
    let (file, source) = match (args.file, args.code) {
        (None, Some(code)) => (None, code),
        (Some(file), None) => (
            Some(file.clone()),
            fs::read_to_string(file.clone()).unwrap_or_else(|err| {
                panic!(
                    "could not read from file `{}`: {}",
                    file.clone().as_path().display(),
                    err
                )
            }),
        ),
        (Some(_), Some(_)) => panic!("only 1 of `--file` or `--code` is allowed"),
        (None, None) => {
            let mut s = String::new();
            io::stdin()
                .read_to_string(&mut s)
                .unwrap_or_else(|err| panic!("could not read from stdin: {}", err));
            (None, s)
        }
    };
    if source.ends_with("\n") {
        (file, source, args.debug)
    } else {
        (file, source + "\n", args.debug)
    }
}
