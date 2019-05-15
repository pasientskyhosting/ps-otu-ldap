import React from 'react';
import { FormGroup, Intent, Button, Elevation, Card, HTMLSelect, Callout, ButtonGroup, InputGroup, ITagProps, Tag, } from "@blueprintjs/core";

import ValidatedInputGroup from './ValidatedInputGroup';
import APIService, { IGroupCustomProps } from '../services/APIService';
import LDAPGroupSearch from './LDAPGroupSearch';

interface IProps {    
    onGroupCreateHandler: (success: boolean) => void, 
}

interface IState {    
    ldap_group_name: string
    group_name: string
    lease_time: number    
    errorMessage: string
    custPropKey: string
    custPropValue: string
    tags: IGroupCustomProps[]
}

type Nullable<T> = T | null

const lease_time_default = 720

export default class GroupCreate extends React.Component<IProps, IState> {
    
    private custPropsKeyRef: Nullable<HTMLInputElement>   

    constructor (props: IProps) {
                
        super(props)

        this.custPropsKeyRef = null

        this.state = {
            tags: [],
            ldap_group_name: "",
            group_name: "",
            lease_time: 720,
            errorMessage: "",
            custPropKey: "",
            custPropValue: ""
        }        

    }    

    private addTag(key: string, value: string) {

        let tag: IGroupCustomProps = { key, value }                
        let tags = this.state.tags.concat(tag);
        
        this.setState({
            tags: tags,
            custPropKey: "",
            custPropValue: ""            
        })

        this.custPropsKeyRef && this.custPropsKeyRef.focus()

    }

    private removeTag(key?: string) {
        
        let tags = this.state.tags

        this.state.tags.map((custprops: IGroupCustomProps, index: number) => {                        
            let search = custprops.key+custprops.value
            if(search == key) {                
                tags.splice(index, 1);
            }
        })
        
        this.setState({
            tags: tags,
            custPropKey: "",
            custPropValue: ""   
        })

    }

    private renderTags() {

        let tags = [] as JSX.Element[];        
        const onRemove = (e: MouseEvent<HTMLButtonElement>, tagProps: ITagProps) => { this.removeTag(tagProps.id || undefined) }      

        // Create all tags
        this.state.tags.map((custprops: IGroupCustomProps) => {                        
            tags.push(
                <Tag
                    className="tags" 
                    minimal={true}
                    onRemove={onRemove}           
                    key={custprops.key+custprops.value}
                    id={custprops.key+custprops.value}            
                     >
                    {custprops.key}={custprops.value}                    
                </Tag>
            )            
        })

        return tags

    }

    private async onSubmit() {        
    
       const group = await APIService.groupCreate(this.state.ldap_group_name, this.state.group_name, this.state.lease_time, this.state.tags)
    
        if(!APIService.success) {
            this.setState({
                errorMessage: APIService.error.error.messages[0].value
            })
        } else {
            this.setState({
                errorMessage: "",
                ldap_group_name: "None",
                custPropKey: "",
                custPropValue: "",
                group_name: "",
                lease_time: lease_time_default
                
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
                    
                    <ButtonGroup                    
                    >
                        <ValidatedInputGroup
                            inputRef={(input) => this.custPropsKeyRef = input}                             
                            value={this.state.custPropKey}                                               
                            onKeyDown={(e: React.KeyboardEvent<Element>) => {
                                        
                            }}
                            validate={(currentValue: string) => {
                            return true                      
                            }}
                            errorMessage={(currentValue: string) =>{
                                return "Not a valid key"
                            }}
                            placeholder="Key"
                            onChange={(e) => {     
                                this.setState({
                                    custPropKey: e.target.value,                                    
                                })               
                            }}
                        >
                        </ValidatedInputGroup>
                        &nbsp;
                        <ValidatedInputGroup                                                
                            value={this.state.custPropValue}                                               
                            onKeyDown={(e: React.KeyboardEvent<Element>) => {
                                if(e.keyCode == 13) {
                                    this.addTag(this.state.custPropKey, this.state.custPropValue)
                                }        
                            }}
                            validate={(currentValue: string) => {
                                return true
                            }}
                            errorMessage={(currentValue: string) =>{
                                return "Not a valid key"
                            }}
                            placeholder="Value"
                            onChange={(e) => {                    
                                this.setState({
                                    custPropValue: e.target.value                                    
                                })
                            }}
                        >
                        </ValidatedInputGroup>      
                    </ButtonGroup>                                        
                    </FormGroup>                      

                    {this.renderTags()}                 
                         
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