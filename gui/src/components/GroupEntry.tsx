import React from 'react';
import { HTMLSelect, Tag, ButtonGroup, FormGroup, InputGroup, ITagProps, Button } from "@blueprintjs/core";
import UserOptions from './UserOptions'
import APIService, { IGroup, IUser, IGroupCustomProps } from '../services/APIService';
import GroupOptions from './GroupOptions';
import ValidatedInputGroup from './ValidatedInputGroup';

interface IProps {    
  onGroupUpdateHandler: (success: boolean, status_code: number) => void,  
  onGroupDeleteHandler: (success: boolean, status_code: number) => void,    
  editMode: boolean,
  group: IGroup
  user?: IUser | undefined  
}

interface IState {    
    group: IGroup    
    mangledGroup: IGroup
    user?: IUser | undefined
    errorMessage?: string
    custPropKey: string
    custPropValue: string    
    disableUpdate: boolean
}

type Nullable<T> = T | null

export default class GroupEntry extends React.Component<IProps, IState> {
    
    private custPropsKeyRef: Nullable<HTMLInputElement>   

    constructor(props: IProps) {

        super(props)    

        this.custPropsKeyRef = null

        this.state = {              
            group: this.props.group,
            mangledGroup: Object.assign({}, this.props.group),
            user: this.props.user || undefined,
            custPropKey: "",
            custPropValue: "",
            disableUpdate: true
        }

    }

    private formatSeconds(sec: number) {

        var seconds = sec

        var days = Math.floor(seconds / (60*24))        
        seconds  -= days*60*24

        var hrs   = Math.floor(seconds / 60)
        seconds  -= hrs*60        
        
        if (days === 0) {
            return hrs+" hour(s)"
        } else {
            return days+" day(s)"
        }

    }

    private onGroupDeleteHandler(success: boolean, status_code: number) {        
        this.props.onGroupDeleteHandler(success, status_code)
    }

    private onGroupUpdateHandler(success: boolean, status_code: number) {        
        this.props.onGroupUpdateHandler(success, status_code)
    }

    private renderEditLeaseTime() {

        return (
            <HTMLSelect                                                  
                value={this.state.mangledGroup.lease_time}  
                onChange={(e: React.ChangeEvent<HTMLSelectElement>) => {                    
                    
                    let mangled = this.state.mangledGroup
                    mangled.lease_time = parseInt(e.target.value)

                    this.setState({
                        mangledGroup: mangled,
                        disableUpdate: false               
                    })                    

                }}                
            >                    
            <option value="60">{this.formatSeconds(60)}</option>
            <option value="720">{this.formatSeconds(720)}</option>
            <option value="1440">{this.formatSeconds(1440)}</option>
            <option value="10080">{this.formatSeconds(10080)}</option>
            <option value="20160">{this.formatSeconds(20160)}</option>
            </HTMLSelect>
        )

    }   

    private renderEditGroupName() {

        return (            
            <ValidatedInputGroup       
                key={this.state.group.group_name}                 
                value={this.state.mangledGroup.group_name}                 
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
                onKeyEnter={() => {                
                    this.onGroupUpdateAPIHandler()
                }}  
                onChange={(e) => {                    
                    
                    let mangled = this.state.mangledGroup
                    mangled.group_name = e.target.value

                    this.setState({
                        mangledGroup: mangled,
                        disableUpdate: false 
                    })

                    console.log(this.state)

                }}
                
            />  
        )
    }    

    private async onGroupUpdateAPIHandler() {

        let response = await APIService.groupUpdate(this.state.group.group_name, this.state.mangledGroup)
        
        // Update group name
        if(APIService.success) {
            
            if(response) {                
                this.setState({
                    group: response,
                    disableUpdate: true                    
                })
            }            
        }        

        this.props.onGroupUpdateHandler(APIService.success, APIService.status)

    }

    private renderTags() {

        let tags = [] as JSX.Element[];        
        const onRemove = (this.props.editMode) ? (e: React.MouseEvent<HTMLButtonElement>, tagProps: ITagProps) => { this.removeTag(tagProps.id || undefined) } : undefined        

        // Create all tags
        this.state.mangledGroup.custom_properties.map((custprops: IGroupCustomProps) => {                        
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

    private addTag(key: string, value: string) {

        let tag: IGroupCustomProps
        tag = { key, value }
        
        let mangledGroup = Object.assign({},this.state.mangledGroup)
        mangledGroup.custom_properties = mangledGroup.custom_properties.concat(tag);
        
        this.setState({
            mangledGroup: mangledGroup,
            custPropKey: "",
            custPropValue: "",
            disableUpdate: false   
        })

        this.custPropsKeyRef && this.custPropsKeyRef.focus()

    }

    private removeTag(key?: string) {
        
        let mangledGroup = Object.assign({},this.state.mangledGroup)                

        this.state.mangledGroup.custom_properties.map((custprops: IGroupCustomProps, index: number) => {                        
            let search = custprops.key+custprops.value
            if(search == key) {                
                mangledGroup.custom_properties.splice(index, 1);       
            }
        })
        
        this.setState({
            mangledGroup: mangledGroup,
            custPropKey: "",
            custPropValue: "",
            disableUpdate: false     
        })

    }

    private renderAddTags() {
        return (
            <FormGroup>            
                <ButtonGroup                    
                >
                    <ValidatedInputGroup 
                        inputRef={(input) => this.custPropsKeyRef = input}
                        value={this.state.custPropKey}                        
                        validate={(currentValue: string) => {
                            return (currentValue == "" ) ? false : true
                        }}
                        errorMessage={(currentValue: string) =>{
                            return "Not a valid key"
                        }}
                        placeholder="Key"
                        onChange={(e) => {     
                            this.setState({
                                custPropKey: e.target.value,
                                disableUpdate: false
                            })               
                        }}
                    >
                    </ValidatedInputGroup>
                    &nbsp;
                    <ValidatedInputGroup                                                
                        value={this.state.custPropValue}                                                             
                        onKeyEnter={() => {                
                            this.addTag(this.state.custPropKey, this.state.custPropValue)
                        }}
                        validate={(currentValue: string) => {
                            return (currentValue == "" ) ? false : true
                        }}
                        errorMessage={(currentValue: string) =>{
                            return "Not a valid value"
                        }}
                        placeholder="Value"
                        onChange={(e) => {                    
                            this.setState({
                                custPropValue: e.target.value,
                                disableUpdate: false
                            })
                        }}
                    >
                    </ValidatedInputGroup>   
                    <Button
                        icon="plus"                                                
                        text=""                        
                        onClick={(e: React.MouseEvent<HTMLElement, MouseEvent> ) => {
                            if(this.validateTagKey(this.state.custPropKey) && this.validateTagValue(this.state.custPropValue)) {
                                this.addTag(this.state.custPropKey, this.state.custPropValue)
                            }                                    
                        }}
                    ></Button>             
                </ButtonGroup>
            </FormGroup> 
        )
    }

    private validateTagKey(str: string): boolean {
        return (str == "" || !(str.match(/^[_\-0-9a-z]+$/g))) ? false : true 
    }

    private validateTagValue(str: string): boolean {
        return (str == "" || !(str.match(/^[_\-0-9a-z]+$/g))) ? false : true
    }

    private renderGroupRow() {       
                    
        let user = this.state.user ? this.state.user : undefined        

        return (
            <tr key={this.state.group.group_name}>
                <td>{this.state.group.ldap_group_name}</td>
                <td>                                
                    {this.props.editMode ? (
                        this.renderEditGroupName()
                    ) : (                            
                        this.state.group.group_name
                    )}                                
                </td>
                <td>                            
                    {this.props.editMode ? (
                        this.renderEditLeaseTime()
                    ) : (    
                        this.formatSeconds(this.state.group.lease_time)
                    )}
                </td>
                <td>
                    { this.props.editMode && this.renderAddTags() }                    
                    {this.renderTags()}                    
                </td>
                <td>
                    
                    {user? (
                        <UserOptions
                        username={user.username} 
                        password={user.password} 
                        group_name={this.state.group.group_name}
                        />
                    ) : (
                        <UserOptions
                        username=""
                        password=""
                        group_name={this.state.group.group_name}
                        />
                    )}
                </td>        
                { this.props.editMode &&
                    <td>                                    
                    <GroupOptions                                    
                        group={this.state.group}
                        onGroupDeleteHandler={(success: boolean, status_code: number) => this.onGroupDeleteHandler(success, status_code) }                          
                        onGroupUpdateHandler={() => this.onGroupUpdateAPIHandler() }                   
                        disableUpdate={this.state.disableUpdate}   
                    />
                    </td>
                }                    
            </tr>
        )
    }    
    
    render () {
        
        return (       
            this.renderGroupRow()
        )
    }

}   