use std::ffi::{CStr, CString};
use std::os::raw::{c_char, c_int, c_uchar};
use std::path::Path;
use vtracer::{convert_image_to_svg, Config, ColorMode, Hierarchical};
use visioncortex::{ColorImage, PathSimplifyMode};

#[repr(C)]
#[derive(Debug, Clone, Copy)]
pub struct VtracerConfig {
    pub color_mode: c_uchar,
    pub hierarchical: c_uchar,
    pub filter_speckle: usize,
    pub color_precision: c_int,
    pub layer_difference: c_int,
    pub mode: c_uchar,
    pub corner_threshold: c_int,
    pub length_threshold: f64,
    pub max_iterations: usize,
    pub splice_threshold: c_int,
    pub path_precision: u32,
}

impl From<VtracerConfig> for Config {
    fn from(c_config: VtracerConfig) -> Self {
        Config {
            color_mode: match c_config.color_mode {
                0 => ColorMode::Color,
                _ => ColorMode::Binary,
            },
            hierarchical: match c_config.hierarchical {
                0 => Hierarchical::Stacked,
                _ => Hierarchical::Cutout,
            },
            filter_speckle: c_config.filter_speckle,
            color_precision: c_config.color_precision,
            layer_difference: c_config.layer_difference,
            mode: match c_config.mode {
                0 => PathSimplifyMode::None,
                1 => PathSimplifyMode::Polygon,
                _ => PathSimplifyMode::Spline,
            },
            corner_threshold: c_config.corner_threshold,
            length_threshold: c_config.length_threshold,
            max_iterations: c_config.max_iterations,
            splice_threshold: c_config.splice_threshold,
            path_precision: Some(c_config.path_precision),
        }
    }
}

#[no_mangle]
pub extern "C" fn vtracer_default_config() -> VtracerConfig {
    let default_config = Config::default();
    VtracerConfig {
        color_mode: match default_config.color_mode {
            ColorMode::Color => 0,
            ColorMode::Binary => 1,
        },
        hierarchical: match default_config.hierarchical {
            Hierarchical::Stacked => 0,
            Hierarchical::Cutout => 1,
        },
        filter_speckle: default_config.filter_speckle,
        color_precision: default_config.color_precision,
        layer_difference: default_config.layer_difference,
        mode: match default_config.mode {
            PathSimplifyMode::None => 0,
            PathSimplifyMode::Polygon => 1,
            PathSimplifyMode::Spline => 2,
        },
        corner_threshold: default_config.corner_threshold,
        length_threshold: default_config.length_threshold,
        max_iterations: default_config.max_iterations,
        splice_threshold: default_config.splice_threshold,
        path_precision: default_config.path_precision.unwrap_or(2),
    }
}

#[no_mangle]
pub extern "C" fn vtracer_convert_file(
    input_path: *const c_char,
    output_path: *const c_char,
    config: *const VtracerConfig,
) -> c_int {
    if input_path.is_null() || output_path.is_null() || config.is_null() {
        return -1; // Invalid parameters
    }

    unsafe {
        let input_str = match CStr::from_ptr(input_path).to_str() {
            Ok(s) => s,
            Err(_) => return -2, // Invalid input path
        };

        let output_str = match CStr::from_ptr(output_path).to_str() {
            Ok(s) => s,
            Err(_) => return -3, // Invalid output path
        };

        let rust_config = Config::from(*config);

        match convert_image_to_svg(Path::new(input_str), Path::new(output_str), rust_config) {
            Ok(_) => 0,
            Err(_) => -4, // Conversion failed
        }
    }
}

#[no_mangle]
pub extern "C" fn vtracer_convert_bytes(
    image_data: *const c_uchar,
    image_len: usize,
    width: usize,
    height: usize,
    config: *const VtracerConfig,
    output: *mut *mut c_char,
) -> c_int {
    if image_data.is_null() || config.is_null() || output.is_null() {
        return -1; // Invalid parameters
    }

    unsafe {
        // Check data length
        if image_len != width * height * 4 {
            return -2; // Data length error
        }

        let bytes = std::slice::from_raw_parts(image_data, image_len);
        // ColorImage expects RGBA byte array, so use raw data directly
        let pixels = bytes.to_vec();

        let color_image = ColorImage { pixels, width, height };
        let rust_config = Config::from(*config);

        match vtracer::convert(color_image, rust_config) {
            Ok(svg) => {
                let svg_string = svg.to_string();
                match CString::new(svg_string) {
                    Ok(c_string) => {
                        *output = c_string.into_raw();
                        0
                    }
                    Err(_) => -3, // String conversion failed
                }
            }
            Err(_) => -4, // Image conversion failed
        }
    }
}

#[no_mangle]
pub extern "C" fn vtracer_free_string(s: *mut c_char) {
    if !s.is_null() {
        unsafe {
            let _ = CString::from_raw(s);
        }
    }
}
