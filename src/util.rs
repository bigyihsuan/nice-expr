use crate::{eval::value::NativeFunc, prelude::RuntimeError};

pub fn assert_at_least_args(want: usize, got: usize) -> Result<(), RuntimeError> {
    if got < want {
        Err(RuntimeError::NotEnoughArguments { want, got })
    } else {
        Ok(())
    }
}
pub fn assert_exactly_args(want: usize, got: usize) -> Result<(), RuntimeError> {
    if got != want {
        Err(RuntimeError::NotEnoughArguments { want, got })
    } else {
        Ok(())
    }
}
