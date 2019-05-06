import React from 'react';
import { FormGroup, Intent, Button, Elevation, Card, HTMLSelect, ControlGroup, Callout, Label } from "@blueprintjs/core";
import ValidatedInputGroup from './ValidatedInputGroup';
import APIService from '../services/APIService';

interface IProps {    
    onGroupCreateHandler: (success: boolean, status_code: number) => void,        
}

interface IState {    
    group_name: string
    lease_time: number
    connecting?: boolean
    errorMessage: string
}

export default class GroupCreate extends React.Component<IProps, IState> {
   
    constructor (props: IProps) {
        
        super(props)

        this.state = {
            group_name: "",
            lease_time: 720,
            errorMessage: ""
        }        

    }

    private async onSubmit() {

        this.setState({
            connecting: true
        })
    
       const group = await APIService.groupCreate(this.state.group_name, this.state.lease_time)

       console.log(group)
    
        if(!APIService.success) {
            this.setState({
                errorMessage: "Error while creating group"
            })
        }
          
        // call login handler
        this.props.onGroupCreateHandler(APIService.success, APIService.status)
  
    }

    render () {

        return (                   
            
            <div className="groups-create card">            
                <Card interactive={true} elevation={Elevation.FOUR}>                
                    <div className="groups-create-content">
                    <h2>Create Group</h2>
                    {this.state.errorMessage ? <Callout title="Error" className="group-create-error-message" intent={Intent.DANGER} >{this.state.errorMessage}</Callout> : null }
                   
                    <FormGroup                     
                     label="Group name"
                     labelFor="text-input"                     
                    >                                       
                    <ValidatedInputGroup                                                                                                          
                        placeholder="Insert group name here..."  
                        validate={(currentValue: string) => {
                            if(currentValue.length == 0 || (currentValue.length > 3 && currentValue.match(/^[_\-0-9a-z]+$/g)) ) {
                                return true
                            } else {
                                return false
                            }                       
                        }}
                        errorMessage={(currentValue: string) =>{
                            return "Length must be greater then 3, and be URL friendly"
                        }}
                        large={false}                        
                        onKeyDown={(e: React.KeyboardEvent) => {                
                            if(e.keyCode == 13) this.onSubmit()
                        }}
                        onChange={(e: React.ChangeEvent<HTMLInputElement>) => {                            
                            this.setState({
                                group_name: e.target.value
                            })
                        }}
                    />       
                    </FormGroup>                                 
                    <FormGroup                     
                     label="Lease time"
                     labelFor="lease-time"                     
                    >
                    <HTMLSelect                        
                        id="lease-time"
                        value={this.state.lease_time}                        
                        onChange={(e: React.ChangeEvent<HTMLSelectElement>) => {                            
                            this.setState({
                                lease_time: parseInt(e.target.value)
                            })
                        }}
                    >                    
                    <option value="60">1 hour</option>
                    <option value="720">12 hours</option>
                    <option value="1440">24 hours</option>
                    <option value="10080">1 week</option>
                    <option value="20160">2 weeks</option>
                    </HTMLSelect>                                        
                   
                    </FormGroup>                     
                         
                    <Button                          
                        style={{ width: "100%", marginTop: "20px" }}              
                        large={false}
                        intent={Intent.PRIMARY}                
                        onClick={() => this.onSubmit()}
                    >Create</Button>
                    </div>
                </Card>
            </div>        
        
        )
        
    }

}



// private async onGroupCreate(group_name: string, lease_time: number) {

//     try {
        
//         let response = await fetch('/v1/api/groups', {
//             method: 'post',    
//             headers: { "Authorization": `Bearer ${this.state.token}`  },                   
//             body: JSON.stringify({  
//                 group_name, lease_time
//             })
//         })

//         if (response.status !== 201) {

//             let responseParsed = await response.json()

//             switch(response.status) {
//                 case 401:                        
//                     this.setState({                                
//                         isVerified: false,
//                         loginFailed: false,                                                
//                     })
//                     break
//                 case 400:                        
//                     this.setState({                                
//                         errorMesssage: responseParsed.validation_error.group_name,                            
//                     })
//                     break
//                 case 409:                        
//                     this.setState({                                
//                         errorMesssage: "Group does already exist with this name",
//                     })
//                     break
//             }                

//         } else {                                    
            
//             this.setState({
//                 isVerified: true,
//                 errorMesssage: "",
//             })
//         }                        

//     } catch (error)
//     {
//         this.setState({                                
//             loginFailed: true,
//             errorMesssage: "Wrong username or password"
//         })
//         console.log(error)
//     }

// }