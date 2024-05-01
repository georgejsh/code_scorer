use std::env;
use reqwest;
use serde::Deserialize;
use std::io::{self, BufRead};
use std::io::Write;
use rpassword::read_password;
use std::collections::HashMap;
use sha256::{digest, try_digest};
use std::fs;
use std::fs::File;
use std::path::Path;
use crate::io::Error;
use reqwest::{Client, Url, multipart};
use reqwest::header::{HeaderMap, HeaderValue, AUTHORIZATION, CONTENT_TYPE};
use serde_json::json;
use rsa::{Pkcs1v15Encrypt, RsaPrivateKey};
use rsa::{RsaPublicKey, pkcs1::DecodeRsaPublicKey};
use crate::reqwest::header;
use chrono;
use queue::Queue;
 use std::borrow::Cow;
//use tokio::fs::File;
//use tokio_util::codec::{BytesCodec, FramedRead};

//let mut URL:&str = "http://localhost:3333/help";
const HOME:&str = "~/.client/";
#[derive(Debug, Deserialize)]
struct Task {
    command: String,
    url: String,
}

#[derive(Debug, Deserialize)]
struct Output {
    status: String,
    error: String,
}

#[derive(Debug, Deserialize,Clone)]
struct Args {
    name: String,
    typ: String,
    secret:bool,
    values: Option<Vec<String>>
}
const SALT:&str = "RcHeDY24MxXV5tZUua3WwPCEmSgsrjpLQKhTzNkybG8AfvF96J";
use std::fs::read_to_string;
const public_rsa:&str = "-----BEGIN RSA PUBLIC KEY-----
MIIBigKCAYEAth7ktUwWlNJO5YRTcqbvf1G8rzrFuPm0VsSfgL8g7sAxsCRL4fbT
OX8KAipgHM01YcYXEW25e8RFPHoWEBdTW1MK5yoZ/3S2ttc/A98LgyBRQrfYMOu7
MBbgD4XkqDOJh9aNIeidPmF56l0seCajmdHQRM/sGIW1+e3bqkFtlAYqq5bzARHv
gZJCYrcKTqBx8CuuvIa97/4qUa9gsR6KZXUmdJb8D2bEnu+TPSyor2su5C5E7NWw
72ZVWWuOTjsBM+Dz0p5kAzZMwhGZD1N+8ihm1Fh2BgeBBq49uCwfk7HqGNDfm0mO
hej7rg9nEDGvpC5GUqHTZt70fWVR0Cb5xhx2BILZUEhW3BvM0/NRKjBDjbOpNDG/
i/7tcLOgP2iglh4rj/JE8EvAquO68jlcqV5aRoL2eqsjwmmA5V6zRRJsZvMihSes
f5i9XXRfbYs/cU6DvzO5SEw7zloCaU6j0HfSlEHK1/5Mfj9MgSbhHsuKTnsqgYHw
5+vYlPV+Zh3rAgMBAAE=
-----END RSA PUBLIC KEY-----";
fn getArgs(url:&str)->Vec<Args>{
    if let Ok(resp) = reqwest::blocking::get(url.to_owned()+"help"){
        if let Ok(out)= resp.text(){
            if let Ok(output)  = serde_json::from_str::<Output>(&out) {
                if let Ok(json) =
                serde_json::from_str::<Vec<Args>>(&output.status){
                    return json;
                }
            }
        }
    }
    return vec![]
}
fn getBaseURL()->String{
    if let Ok(url) = read_to_string(&(HOME.to_owned()+&"url")) {
        //println!("JWT taken {}",tjwt);
        return url;
        // Consumes the iterator, returns an (Optional) String
    }else{
        panic!("No Base URL present!!");
    }
}
fn writeBaseURL(out:&str){
    let mut file = File::create(HOME.to_owned()+"url").unwrap_or_else(|error| {
        panic!("URL file create error!!");
    });
    file.write_all(out.as_bytes()).unwrap_or_else(|error| {
        panic!("URL file save error!!");
    });;
}
fn runTask(url:&str,options:Vec<String>){
   // println!("{:#?} {:#?}", url,options);
    let resp = reqwest::blocking::get(url.to_owned()+"help").unwrap_or_else(|error| {
        panic!("Network not reachable");
    }
    ).text().unwrap_or_else(|error| {
     "".to_string()
    });
    let output : Output = serde_json::from_str(&resp).expect("Network not reachable");
    if(output.error !=""){
        panic!(" {}",output.error );
    }
    let mut json: Vec<Args> =
        serde_json::from_str(&output.status).expect("JSON was not well-formatted");
    
    let mut jwt :String = "".to_string();
    if options.len()>1 && options[1]=="help"{
    
        print!("Usage: {} " ,options[0]);
        for opt in &json{
            if opt.typ.contains("..."){
                print!("[");
            }
            if opt.typ=="string"{
                print!("{} ",opt.name);
            }else if opt.typ=="file"{
                print!("{}<file> ",opt.name);
            }else if opt.typ=="enum"{
                print!("{:?} ",opt.values.as_ref().expect("No list given"));
            }
            
            if opt.typ.contains("..."){
                print!("]... ");
            }
        } 
        println!("");
        return    
    }
    if options[0] != "login" && options[0] != "logout" {
        if let Ok(tjwt) = read_to_string(&(HOME.to_owned()+&"jwt")) {
            //println!("JWT taken {}",tjwt);
            jwt=tjwt.clone();
            // Consumes the iterator, returns an (Optional) String
        }else {
            panic!("Please Log In!!");
        }
    }
    //println!("{:#?}", json);
    let mut req=1;
    let mut url_to_send:String=url.to_string();
    //let mut query= HashMap::new();
    let mut form = reqwest::blocking::multipart::Form::new();
    let mut rng = rand::thread_rng();
    let pub_key = RsaPublicKey::from_pkcs1_pem(public_rsa).unwrap();
    //let slice =Cell::from_mut(&mut json);
    //let mut json_inp = Cow::from(&json[..]);
    let mut q= Queue::new();
    for i in 1..json.len(){
        q.queue(json[i].clone()).unwrap();
    }
    while(!q.is_empty()) {
        let opt=q.dequeue().unwrap();
        if options.len()==req {
            let mut num_arg =1;
            if opt.typ.contains("..."){
                print!("Number of arguments for  {} : ",opt.name);
                std::io::stdout().flush().unwrap();
                let stdin = io::stdin();
                let mut iterator = stdin.lock().lines();
                num_arg = iterator.next().unwrap().unwrap().trim().parse().expect("Input not an integer");;
            }
            for i in 0..num_arg {
                let mut name_suffix="".to_string();
                if opt.typ.contains("..."){
                    name_suffix=(i+1).to_string();
                }
                print!("Enter value for {} : ",opt.name);
                std::io::stdout().flush().unwrap();
                let mut line1="".to_string();
                if opt.typ == "file"{
                    let stdin = io::stdin();
                    let mut iterator = stdin.lock().lines();
                    let local_filepath = iterator.next().unwrap().unwrap();
                    
                    //let file_fs = fs::read(local_filepath).unwrap_or_else(|error| {
                    //    panic!("File not present");
                    //});
                    //let part = reqwest::blocking::multipart::Part::bytes(file_fs);
                    //form=form.part(opt.name.clone()+&name_suffix,part);
                    let part = reqwest::blocking::multipart::Part::file(local_filepath).unwrap_or_else(|error| {
                            panic!("File not present");
                        });;
                    form=form.part(opt.name.clone()+&name_suffix,part);
                }
                else if opt.secret{
                    line1 = read_password().unwrap();
                    let t_n=chrono::offset::Local::now().timestamp_nanos().to_string();
                    line1 = line1 + SALT;
                    line1  =digest(line1) +"_"+ &t_n;
                    //println!("{:?}",line1);
                    
                    let enc = pub_key.encrypt(&mut rng, Pkcs1v15Encrypt, &line1.as_bytes()).expect("failed to encrypt");
                    line1=hex::encode(&enc);
                    form=form.text(opt.name.clone()+&name_suffix,line1);
                    //println!("{:?}",line1)
                }else if opt.typ == "enum" {
                    print!("\nAvailable Values are {:?} : ",opt.values.as_ref().expect("No list given"));
                    std::io::stdout().flush().unwrap();
                    let stdin = io::stdin();
                    let mut iterator = stdin.lock().lines();
                    line1 = iterator.next().unwrap().unwrap();
                    if !opt.values.as_ref().expect("Enum Values can't be NULL").contains(&line1) {
                        panic!("Input not present in available values")
                    }
                    form=form.text(opt.name.clone()+&name_suffix,line1.clone());
                    let more_args=getArgs(&(url.to_owned()+&line1+"/"));
                    if more_args.len() != 0  {
                        url_to_send.push_str(&(line1+"/"));
                        //slice.update(|x| x.extend(more_args));
                        //json_inp.to_mut().extend(more_args);
                        for j in 1..more_args.len(){
                            q.queue(more_args[j].clone());
                        }
                    }
                }
                else {
                    let stdin = io::stdin();
                    let mut iterator = stdin.lock().lines();
                    line1 = iterator.next().unwrap().unwrap();
                    form=form.text(opt.name.clone()+&name_suffix,line1);
                }
                
            }
        }else {
            let mut num_arg =1;
            if opt.typ.contains("..."){
                num_arg=options.len()-req;
            }
            for i in 0..num_arg {
                let mut name_suffix="".to_string();
                if opt.typ.contains("..."){
                    name_suffix=(i+1).to_string();
                }
                if opt.typ == "file"{
                    /*let file_fs = fs::read(options[req].clone()).unwrap_or_else(|error| {
                        panic!("File not present");
                    });
                    let part = reqwest::blocking::multipart::Part::bytes(file_fs);
                    form=form.part(opt.name.clone()+&name_suffix,part);*/
                    let part = reqwest::blocking::multipart::Part::file(options[req].clone()).unwrap_or_else(|error| {
                        panic!("File not present");
                    });;
                form=form.part(opt.name.clone()+&name_suffix,part);
                }
                else if opt.typ == "enum" {
                    if !opt.values.as_ref().expect("Enum Values can't be NULL").contains(&options[req]) {
                        panic!("Input not present in available values")
                    }
                    form=form.text(opt.name.clone()+&name_suffix,options[req].clone());
                    let more_args=getArgs(&(url.to_owned()+&options[req].clone()+"/"));
                    if more_args.len() != 0  {
                        url_to_send.push_str(&(options[req].clone()+"/"));
                        //slice.update(|x| x.extend(more_args));
                        //json_inp.to_mut().extend(more_args);
                        for j in 1..more_args.len(){
                            q.queue(more_args[j].clone());
                        }
                    }
                }else if opt.secret {
                    let t_n=chrono::offset::Local::now().timestamp_nanos().to_string();
                    let mut line1 = options[req].clone() + SALT;
                    line1  =digest(line1) +"_"+ &t_n;
                    //println!("{:?}",line1);
                    
                    let enc = pub_key.encrypt(&mut rng, Pkcs1v15Encrypt, &line1.as_bytes()).expect("failed to encrypt");
                    line1=hex::encode(&enc);
                    form=form.text(opt.name.clone()+&name_suffix,line1);
                }else{
                    form=form.text(opt.name.clone()+&name_suffix,options[req].clone());
                }
                
                req+=1;
            }
        }
        
    }    
    //println!("{:#?}", ans);
    /*let mut querystring:String =url.to_owned()+"?";
    for (k,v) in ans{
        querystring +=&(k+"="+&v+"&")
    }*/
    //println!("{:#?}",form);
  
    
    let t_n=chrono::offset::Local::now().timestamp_nanos().to_string();
    //print!("Utc value: {}\n",t_n);
    let enc = pub_key.encrypt(&mut rng, Pkcs1v15Encrypt, &t_n.as_bytes()).expect("failed to encrypt");
    form=form.text("utc",hex::encode(&enc));
    let client = reqwest::blocking::Client::new();
    let mut headers: HeaderMap = HeaderMap::new();
    //let header_string = format!("{}", jwt);
    let header_value = HeaderValue::from_str(&jwt).unwrap_or_else(|error| {
        panic!("Header error!! Pls login Again");
    });

    headers.insert(AUTHORIZATION, header_value.clone());
    //print!("authorisation: {:?}\n",header_value);
    //headers.insert(CONTENT_TYPE, HeaderValue::from_static("application/json"));

    let resp = client.post(url_to_send)
      .headers(headers.clone())
      .multipart(form).send().unwrap_or_else(|error| {
        panic!("Network not reachable");
    });
    /*if options[0] == "login" && !resp.status().is_success() {
        panic!("Login Failed!!!");
    }
    if options[0] == "logout" && !resp.status().is_success() {
        panic!("Logout Failed!!!");
    }*/
    match resp.headers().get(header::CONTENT_TYPE) {
        Some(ct) if ct == "application/pdf" =>{
            let line:String=resp.headers().get(header::CONTENT_DISPOSITION).unwrap_or(&HeaderValue::from_str("attachment; filename=output.pdf").unwrap()).to_str().unwrap().to_string();
            let name = &line[21..];
            let content = resp.bytes().unwrap_or_else(|error| {
                panic!("No data present in response!!");
            });
        
            // Save the content to a local file
            let mut file = File::create(name).unwrap_or_else(|error| {
                panic!("Output File Creation Failed!! {:?}",name);
            });
            file.write_all(&content).unwrap_or_else(|error| {
                panic!("FIle write failed!!");
            });
            
            println!("File Downloaded as {:?}",name);
            return
        }
        Some(ct) if ct == "image/jpeg"  => {
            let line:String=resp.headers().get(header::CONTENT_DISPOSITION).unwrap_or(&HeaderValue::from_str("attachment; filename=output.jpg").unwrap()).to_str().unwrap().to_string();
            let name = &line[21..];
            let content = resp.bytes().unwrap_or_else(|error| {
                panic!("No data present in response!!");
            });
        
            // Save the content to a local file
            let mut file = File::create(name).unwrap_or_else(|error| {
                panic!("Output File Creation Failed!! {:?}",name);
            });
            file.write_all(&content).unwrap_or_else(|error| {
                panic!("File write failed!!");
            });
            println!("File Downloaded as {:?}",name);
            return
        }
        Some(ct) if ct == "text/plain"  => {
            let line:String=resp.headers().get(header::CONTENT_DISPOSITION).unwrap_or(&HeaderValue::from_str("attachment; filename=output.txt").unwrap()).to_str().unwrap().to_string();
            let name = &line[21..];
            
            let content = resp.bytes().unwrap_or_else(|error| {
                panic!("No data present in response!!");
            });
        
            // Save the content to a local file
            let mut file = File::create(name).unwrap_or_else(|error| {
                panic!("Output File Creation Failed!! {:?}",name);
            });
            file.write_all(&content).unwrap_or_else(|error| {
                panic!("File write failed!!");
            });
            println!("File Downloaded as {:?}",name);
            return
        }
        None =>{}
        Some(ct)=>{
            
        }
    }
    let result = resp.text().unwrap_or_else(|error| {
     "".to_string()
    });
    let out:Output = serde_json::from_str(&result).expect("Output parse error");
    if(out.error !=""){
        panic!(" {}",out.error );
    }
    if options[0] != "login" && options[0] != "logout" {
        println!("{}",out.status);
    }else {
        fs::create_dir_all(HOME).unwrap_or_else(|error| {
            panic!("Create Home directory failed!!");
        });
        if options[0] == "login"{
            let mut file = File::create(HOME.to_owned()+"jwt").unwrap_or_else(|error| {
                panic!("Token file create error!!");
            });
            file.write_all(out.status.as_bytes()).unwrap_or_else(|error| {
                panic!("Token file save error!!");
            });;
        } else {
            let mut file = fs::remove_file(HOME.to_owned()+"jwt").unwrap_or_else(|error| {
                panic!("Token file remove error!!");
            });;
        }
        println!("{} Successful!!",options[0])
    }
    
}
fn main() {
    let args: Vec<_> = env::args().collect();
    if args.len()==1 {
        return
    }
    if args[1]=="set_base" {
        print!("Enter bse url : ");
        std::io::stdout().flush().unwrap();
        let stdin = io::stdin();
        let mut iterator = stdin.lock().lines();
        let base = iterator.next().unwrap().unwrap();
        writeBaseURL(&base);
        return
    }
    let resp = reqwest::blocking::get(getBaseURL()).unwrap_or_else(|error| {
        panic!("Network not reachable");
    }
    ).text().unwrap_or_else(|error| {
     "".to_string()
    });
    let output : Output = serde_json::from_str(&resp).expect("Network not reachable");
    if(output.error !=""){
        panic!(" {}",output.error );
    }
    let commands: Vec<Task> = serde_json::from_str(&output.status).expect("Network not reachable");
    //println!("{:#?}", resp);
    //println!("{:#?}", commands);
    if args[1]=="help"{
        
        print!("Commands available:" );
        for comm in &commands{
            print!(" {},",comm.command);
        } 
        println!("");
        return    
    }

    for comm in &commands {
        if comm.command == args[1] {
            runTask(&comm.url,args[1..].to_vec())
        }
    }
}

// todo update itself