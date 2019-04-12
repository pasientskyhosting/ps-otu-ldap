import React from 'react';
import { InputGroup, FormGroup, Intent, Button, Elevation, Card, HTMLSelect, ControlGroup, Callout, Tooltip, Position } from "@blueprintjs/core";
import { INPUT } from '@blueprintjs/core/lib/esm/common/classes';
import ValidatedInputGroup from './ValidatedInputGroup';
import ValidatedNumericInput from './ValidatedNumericInput';

interface IProps {    
    onSubmit: (group_name: string, lease_time: number) => Promise<void>,
    errorMessage?: string
}

interface IState {    
    group_name: string
    lease_time: number
    connecting?: boolean
}

export default class GroupCreate extends React.Component<IProps, IState> {
   
    constructor (props: IProps) {
        
        super(props)

        this.state = {
            group_name: "",
            lease_time: 720
        }

    }

    private async onSubmit() {

        this.setState({
          connecting: true
        })      
  
        await this.props.onSubmit(this.state.group_name, this.state.lease_time)
  
        this.setState({
            connecting: false
        })          
  
    }

    render () {

        return (                   
            
            <div className="groups-create card">            
                <Card interactive={true} elevation={Elevation.FOUR}>                
                    <div className="groups-create-content">
                    <h2>Create Group</h2>
                    {this.props.errorMessage ? <Callout title="Error" className="group-create-error-message" intent={Intent.DANGER} >{this.props.errorMessage}</Callout> : null }
                    <FormGroup>                    
                    <ValidatedInputGroup                                                           
                        placeholder="Group name"  
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
                    <FormGroup>
                    <ControlGroup 
                        fill={true} 
                        vertical={false}
                    >                    
                    <HTMLSelect defaultValue="720">
                    <option value="0">Custom</option>
                    <option value="60">1 hour</option>
                    <option value="720">12 hours</option>
                    <option value="1440">24 hours</option>
                    <option value="10080">1 week</option>
                    <option value="20160">2 weeks</option>
                    </HTMLSelect>
                    <ValidatedInputGroup                                                           
                        placeholder="Lease in min."  
                        validate={(currentValue: string) => {
                            if(currentValue.length == 0 || (currentValue.length > 3 && currentValue.match(/^[_\-0-9a-z]+$/g)) ) {
                                return true
                            } else {
                                return false
                            }                       
                        }}
                        errorMessage={(currentValue: string) =>{
                            return "Lease time must be greater than 60 and less than 20160"
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
                    </ControlGroup>
                   
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