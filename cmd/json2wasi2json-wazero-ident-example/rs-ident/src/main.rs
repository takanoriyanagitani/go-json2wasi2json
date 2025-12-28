use std::io;
use std::process::ExitCode;

use io::Write;

fn sub() -> Result<(), io::Error> {
    let i = io::stdin();
    let mut il = i.lock();

    let o = io::stdout();
    let mut ol = o.lock();

    io::copy(&mut il, &mut ol)?;
    ol.flush()
}

fn main() -> ExitCode {
    sub().map(|_| ExitCode::SUCCESS).unwrap_or_else(|e| {
        eprintln!("{e}");
        ExitCode::FAILURE
    })
}
