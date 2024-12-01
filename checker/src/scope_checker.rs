use std::process::exit;

use shared::logger::{Logger, LoggerImpl};
use shared::code::{file_code::{FileCode, FileCodeImpl}, function::FunctionImpl};
use shared::token::token::TokenImpl;

fn throw_value_already_defined(name: &String, trace: &String) {
    Logger::err(
        "Value already defined",
        &[
            "Choose a different name for the value"
        ],
        &[
            trace.as_str(),
            format!("The value {} is already defined", name).as_str()
        ],
    );

    exit(1);
}

// Analyzes the source code to determine undefined variables
pub fn analyze_scope(source: &FileCode) {
    let functions = source.get_functions();

    for (function_name, function) in functions {
        for token in function.get_body() {
            let token_type = token.get_token_type();
        }
    }
}