use std::{collections::HashMap, process::exit};

use code::{token::TokenImpl, token_type::TokenType, types::ParamTypeImpl};
use logger::{Logger, LoggerImpl};

use shared::code::function::{Function, FunctionImpl};

pub fn check_main_function(functions: &HashMap<String, HashMap<String, Function>>) {
    let mut main_function_optional: Option<Function> = None;

    for (_, file_functions) in functions.iter() {
        for (function_name, function) in file_functions.iter() {
            if function_name == "main" {
                if main_function_optional.is_some() {
                    Logger::err(
                        "Multiple main functions found",
                        &[
                            "There can only be one main function"
                        ],
                        &[
                            main_function_optional.unwrap().get_trace().as_str(),
                            file_functions.get("main").unwrap().get_trace().as_str()
                        ]
                    );

                    exit(1);
                }

                main_function_optional = Some(function.clone());
            }
        }
    }

    let main_function = main_function_optional.unwrap();

    // Main functions should be public in order for it to be transpiled
    // as a function and not an embedded lambda function
    if !main_function.is_public() {
        Logger::err(
            "Main function must be public",
            &[
                "The main function must be public"
            ],
            &[
                main_function.get_trace().as_str()
            ]
        );

        exit(1);
    }
    
    if main_function.get_return_type().get_raw_tokens().get(0).unwrap().get_token_type() != TokenType::Nothing {
        Logger::err(
            "Invalid return type for main function",
            &[
                "The return type of the main function must be Nothing"
            ],
            &[
                main_function.get_trace().as_str()
            ]
        );

        exit(1);
    }
}