import React from 'react';
import { FormGroup, Intent, Button, Elevation, Card, HTMLSelect, Callout, ButtonGroup, InputGroup, } from "@blueprintjs/core";

import ValidatedInputGroup from './ValidatedInputGroup';
import APIService from '../services/APIService';
import LDAPGroupSearch from './LDAPGroupSearch';

interface IProps {    
    onGroupCreateHandler: (success: boolean) => void, 
}

interface IState {    
    ldap_group_name: string
    group_name: string
    lease_time: number    
    errorMessage: string
}

export default class GroupCreate extends React.Component<IProps, IState> {


    constructor (props: IProps) {
        
        super(props)

        this.state = {
            ldap_group_name: "",
            group_name: "",
            lease_time: 720,
            errorMessage: ""
        }        

    }    

    private async onSubmit() {        
    
       const group = await APIService.groupCreate(this.state.ldap_group_name, this.state.group_name, this.state.lease_time)
    
        if(!APIService.success) {
            this.setState({
                errorMessage: APIService.error.error.messages[0].value
            })
        } else {
            this.setState({
                errorMessage: ""
            })
        }
          
        // call login handler
        this.props.onGroupCreateHandler(APIService.success)
  
    }

    private onLDAPGroupChange(ldap_group_name: string) {
        this.setState({
            ldap_group_name: ldap_group_name
        })
    }

    render () {        

        return (                   
            
            <div className="groups-create card">            
                <Card interactive={true} elevation={Elevation.FOUR}>                
                    <div className="groups-create-content">
                    <h2>Create Group</h2>
                    {this.state.errorMessage ? <Callout title="Error" className="group-create-error-message" intent={Intent.DANGER} >{this.state.errorMessage}</Callout> : null }
                    
                                        
                    
                    <FormGroup                     
                     label="LDAP group"
                     labelFor="ldap-groups"                                          
                    >                    
                    <LDAPGroupSearch
                        id="ldap-groups"                         
                        onChange={(e: React.ChangeEvent<HTMLSelectElement>) => this.onLDAPGroupChange(e.target.value)}                        
                    />                    
                    </FormGroup>                       
                   
                    <FormGroup                     
                     label="Group name"
                     labelFor="text-input"                     
                    >                                       
                    <ValidatedInputGroup                                                                                                          
                        placeholder="Group name..."  
                        validate={(currentValue: string) => {
                            if(currentValue.length == 0 || (currentValue.length > 2 && currentValue.match(/^[_\-0-9a-z]+$/g)) ) {
                                return true
                            } else {
                                return false
                            }                       
                        }}
                        errorMessage={(currentValue: string) =>{
                            return "Length must be greater than 2, and be URL friendly"
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
                                                                       
                    <FormGroup                     
                     label="Custom properties"
                     labelFor="custom-props"                     
                    >
                    
                    <ButtonGroup>
                    <InputGroup
                    placeholder="Key"
                    ></InputGroup>&nbsp;
                    <InputGroup
                    placeholder="Value"
                    ></InputGroup>
                   
                    </ButtonGroup>
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