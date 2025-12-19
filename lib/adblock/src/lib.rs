use adblock::Engine;
use adblock::lists::ParseOptions;
use adblock::request::Request;
use std::ffi::CStr;
use std::os::raw::c_char;

pub struct AdblockEngine {
    engine: Engine,
}

#[no_mangle]
pub extern "C" fn adblock_engine_create(rules: *const c_char) -> *mut AdblockEngine {
    if rules.is_null() {
        return std::ptr::null_mut();
    }

    let c_str = unsafe { CStr::from_ptr(rules) };
    let rules_str = match c_str.to_str() {
        Ok(s) => s,
        Err(_) => return std::ptr::null_mut(),
    };

    let filter_lines: Vec<String> = rules_str.lines().map(|s| s.to_string()).collect();
    let filter_refs: Vec<&str> = filter_lines.iter().map(|s| s.as_str()).collect();

    let engine = Engine::from_rules(&filter_refs, ParseOptions::default());

    Box::into_raw(Box::new(AdblockEngine { engine }))
}

#[no_mangle]
pub extern "C" fn adblock_engine_check(
    engine: *mut AdblockEngine,
    url: *const c_char,
    source_url: *const c_char,
    resource_type: *const c_char,
) -> bool {
    if engine.is_null() || url.is_null() || source_url.is_null() || resource_type.is_null() {
        return false;
    }

    let engine = unsafe { &*engine };
    
    let url_str = unsafe { CStr::from_ptr(url) }.to_str().unwrap_or("");
    let source_url_str = unsafe { CStr::from_ptr(source_url) }.to_str().unwrap_or("");
    let resource_type_str = unsafe { CStr::from_ptr(resource_type) }.to_str().unwrap_or("");

    let request = match Request::new(url_str, source_url_str, resource_type_str) {
        Ok(r) => r,
        Err(_) => return false,
    };

    let blocker_result = engine.engine.check_network_request(&request);
    
    blocker_result.matched
}

#[no_mangle]
pub extern "C" fn adblock_engine_destroy(engine: *mut AdblockEngine) {
    if !engine.is_null() {
        unsafe {
            drop(Box::from_raw(engine));
        }
    }
}
