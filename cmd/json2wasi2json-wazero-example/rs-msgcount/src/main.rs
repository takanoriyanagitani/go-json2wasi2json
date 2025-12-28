use std::io;
use std::process::ExitCode;

use io::Read;
use io::Write;

#[derive(serde::Deserialize)]
struct Input {
    pub message: String,
}

#[derive(serde::Serialize)]
struct Output {
    pub message_length: u32,
}

fn sub() -> Result<(), io::Error> {
    let i = io::stdin();
    let il = i.lock();
    let mut taken = il.take(1048576);
    let mut buf: Vec<u8> = vec![];
    taken.read_to_end(&mut buf)?;

    let i: Input = serde_json::from_slice(&buf)?;
    let o = Output {
        message_length: i.message.len() as u32,
    };

    let mut ol = io::stdout().lock();
    serde_json::to_writer(&mut ol, &o)?;
    ol.flush()
}

fn main() -> ExitCode {
    sub().map(|_| ExitCode::SUCCESS).unwrap_or_else(|e| {
        eprintln!("{e}");
        ExitCode::FAILURE
    })
}
