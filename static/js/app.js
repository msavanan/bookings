function Prompt(){
    let toast = function(c){
                            const{msg = "", icon="success", position = "top-end"} = c
                              const Toast = Swal.mixin({
                                toast: true,
                                position: position,
                                icon: icon,
                                title: msg,
                                showConfirmButton: false,
                                timer: 3000,
                                timerProgressBar: true,
                                didOpen: (toast) => {
                                  toast.addEventListener('mouseenter', Swal.stopTimer)
                                  toast.addEventListener('mouseleave', Swal.resumeTimer)
                              }
                             });
  
    Toast.fire({});
   }
  
    let success = function(c){
                              const{ title="", text="", footer=""} = c
                              Swal.fire({
                                          icon: "success",
                                          title: title,
                                          text: text,
                                          footer: footer
                                        });
                              }
      let error = function(c){
                              const{ title="", text="", footer=""} = c
                              Swal.fire({
                                          icon: "error",
                                          title: title,
                                          text: text,
                                          footer: footer
                                        });
                             }
      let custom = async function(c){
        const{
        icon= "",showConfirmButton=true,  msg="", title=""} = c;
        
        const { value: result} = await Swal.fire({
                                                        title: title,
 
                                                        icon: icon,
                                                        html:msg,
                                                        backdrop: false,
                                                        showCancelButton: true,
                                                        showConfirmButton:showConfirmButton,
                                                        focusConfirm: false,
                                                        willOpen: () =>{
                                                         if(c.willOpen !== undefined){
                                                         c.willOpen()
                                                         }
                                                        },
                                                        didOpen: () =>{
                                                          if (c.didOpen !== undefined){
                                                          c.didOpen()
                                                          }
                                                        },
                                                        // preConfirm: () => {
                                                        //   return [
                                                        //     document.getElementById('start').value,
                                                        //     document.getElementById('end').value
                                                        //   ]
                                                        // }
                                                      })
  
       if (result){
         if (result.dismiss !== Swal.DismissReason.cancel ){
           if (result.value !== ""){
             if (c.callback !== undefined){
               c.callback(result)
             }
           }else{
             c.callback(false)
           }
       }else{
         c.callback(false)
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

function checkAvailability(id){
      //notify("Helllo", "success")
  // notifyModal("This is a title", "<em>This is HTML</em>", "error", "yes")
 // attention.toast({msg:"Helllo this is a message", icon:"error"})
 //attention.success({msg:"OOPs....",title: "This is title"})
 //attention.error({msg:"OOPs....",title: "This is title"})

 let html =`
 <form id="check-availability-form" action="" method="POST" novalidate class="needs-validations">
    <div class="row">
      <div class="col">
        <div class="row" id="reservation-dates-modal">
          <div class="col">
            <input disabled required class="form-control" type="text" name="start" id="start" placeholder="Arrival">
          </div>
          <div class="col">
            <input disabled required class="form-control" type="text" name="end" id="end" placeholder="Departure">
          </div>
        </div>
      </div>
    </div>
 `

 
        attention.custom({
          title: "Choose your Dates", 
          msg: html,
          willOpen: function(){
           const elem = document.getElementById("reservation-dates-modal")
           const rp = new DateRangePicker(elem,{
              format:"yyyy-mm-dd",
              showOnFocus:true,
              minDate: new Date(),
            })
          },
          didOpen: function(){
          document.getElementById("start").removeAttribute("disabled")
          document.getElementById("end").removeAttribute("disabled")
          },
          callback: function(result){
          let form = document.getElementById("check-availability-form")
     let formData = new FormData(form)
     formData.append("csrf_token", "{{.CSRFToken}}");
     formData.append("room_id", id);

            fetch("/search-availability-json", {
              method: "post",
              body: formData
            })
            .then(response => response.json())
            .then(data =>{
            if (data.ok){
            console.log("room is available")
            attention.custom({
                icon: "success",
                showConfirmButton: false,
                msg:'<p>Room is available</p>' 
                + '<p><a href="/book-room?id='
                + data.room_id
                + '&s='
                + data.start_date
                + '&e='
                + data.end_date
                +' "class="btn btn-primary">'
                + 'Book Now!</a></p>'
                
                 ,

            })

            }else{
            console.log("room is not available")
            attention.error({
            msg: "No room Available",
            })
            }
            //console.log(data.ok)
            //console.log(data.message)
            })
          }
        });
}