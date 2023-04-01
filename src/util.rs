use crate::eval::RuntimeError;

pub fn assert_at_least_args(want: usize, got: usize) -> Result<(), RuntimeError> {
    if got < want {
        Err(RuntimeError::NotEnoughArguments { want, got })
    } else {
        Ok(())
    }
}
