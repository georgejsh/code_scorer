This is a Client which helps to accept the commands available from server
The expected structure of commands
A single Command Structure
struct Task {
    command: String,
    url: String,
}

The structure of arguments of a commands
struct Args {
    name: String,
    typ: String,
    secret:bool,
    values: Optional Vector<String> // the enum values supported
}

The outputs of any commands api
struct Output {
    status: String,
    error: String,
}

The client supports the login using JWT token.