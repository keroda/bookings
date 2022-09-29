function Prompt() {
    let toast = function(c) {
    const {
     msg = "",
     icon = "success",
     pos = "top-end"
    } = c;
    
    const Toast = Swal.mixin({
     toast: true,
     title: msg,
     icon: icon,
     position: pos,//'top-end',
     showConfirmButton: false,
     timer: 3000,
     timerProgressBar: true,
     didOpen: (toast) => {
        toast.addEventListener('mouseenter', Swal.stopTimer)
         toast.addEventListener('mouseleave', Swal.resumeTimer)
     }
    })
    
    Toast.fire({})
    }
    
    let success = function(c) {
    const {
     icon = "success",
     msg = "",
     title = "",
     buttonText = "Cool!",
     footer = ""
    } = c;
    
    Swal.fire({
     title: title,//'Error!',
     html: msg,//'Do you want to continue',
     icon: icon,//'error',
     confirmButtonText: buttonText//'Cool'
    })
    }
    
    let error = function(c) {
    const {
     icon = "error",
     msg = "",
     title = "",
     buttonText = "Oh, OK!",
     footer = ""
    } = c;
    
    Swal.fire({
     title: title,
     html: msg,
     icon: icon,
     confirmButtonText: buttonText
    })
    }
    
    let custom = async function(c) {
    const {
     icon = "",
     msg = "",
     title = "",
     buttonText = "OK",
     footer = "",
     showConfirmButton = true,
    } = c;
    
    const { value: result } = await Swal.fire({
      icon: icon,
      title: title,
      html: msg,
      backdrop: false,
      focusConfirm: false,
      showCancelButton: true,
      showConfirmButton: showConfirmButton,
      
      willOpen: () => {
      if (c.willOpen !== undefined){
        c.willOpen();
      }
    },
    preConfirm: () => {
    return [
     document.getElementById('start-m').value,
     document.getElementById('end-m').value
    ]
    },
    didOpen: () => {
      if (c.didOpen !== undefined){
        c.didOpen();
      }
    }
    })
    if (result){
      if (result.dismiss !== Swal.DismissReason.cancel){
        if (result.value !== ""){
          if (c.callback !== undefined){
            c.callback(result);
          }
        } else {
          c.callback(false);
        }
      } else {
        c.callback(false);
      }
    }
    
    }
    
    return {
    toast: toast,
    success: success,
    error: error,
    custom: custom
    }
    }