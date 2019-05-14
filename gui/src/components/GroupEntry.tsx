import React from 'react';
import { Intent, Button, HTMLSelect, Tag, ButtonGroup } from "@blueprintjs/core";
import UserOptions from './UserOptions'
import APIService, { IGroup, IUser, IGroupCustomProps } from '../services/APIService';
import GroupOptions from './GroupOptions';
import ValidatedInputGroup from './ValidatedInputGroup';
import { AppToaster } from '../services/Toaster';
import { UNDERLINE } from '@blueprintjs/icons/lib/esm/generated/iconContents';

interface IProps {    
  onGroupUpdateHandler: (success: boolean, status_code: number) => void,  
  onGroupDeleteHandler: (success: boolean, status_code: number) => void,    
  onReturnEditModeHandler: () => boolean,      
  group: IGroup
  user?: IUser | undefined  
}

interface IState {
    tags: JSX.Element[]; 
    group: IGroup
    mangledGroup: IGroup
    user?: IUser | undefined
    errorMessage?: string        
    editMode: boolean,
}

export default class GroupEntry extends React.Component<IProps, IState> {
           
    constructor(props: IProps) {

        super(props)

        var tags = [] as JSX.Element[];          
        const onRemove = (e: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {            
            alert("Yesman")
        }
        

        this.state = {            
            tags: tags,            
            group: this.props.group,
            mangledGroup: {
                ldap_group_name: this.props.group.create_by,
                group_name: this.props.group.group_name,
                lease_time: this.props.group.lease_time,
                custom_properties: this.props.group.custom_properties,
                create_time: this.props.group.create_time,
                create_by: this.props.group.create_by
            },
            editMode: this.props.onReturnEditModeHandler(),
            user: this.props.user || undefined  
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
                        mangledGroup: mangled
                    })

                    console.log(this.state)

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
                onKeyDown={(e: React.KeyboardEvent) => {                
                    
                    if(e.keyCode == 13) {                                                
                        this.onGroupUpdateAPIHandler()
                    } else if (e.keyCode === 27) {                        
                        // TODO exit edit mode
                    }
                    
                }}  
                onChange={(e) => {                    
                    
                    let mangled = this.state.mangledGroup
                    mangled.group_name = e.target.value

                    this.setState({
                        mangledGroup: mangled
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
                console.log("Setting state")
                this.setState({
                    group: response                    
                })
            }            
        }        

        this.props.onGroupUpdateHandler(APIService.success, APIService.status)

    }

    private renderGroupRow () {       
                    
        let user = this.state.user ? this.state.user : undefined

        // Create all tags
        this.props.group.custom_properties.map((props: IGroupCustomProps) => {                        
            this.tags.push(
                <Tag
                    className="tags" 
                    minimal={true}
                    onRemove={ this.props.onReturnEditModeHandler() && onRemove}           
                    key={this.props.group.group_name+"-"+props.key+"+"+props.value} >
                    {props.key}={props.value}
                </Tag>
            )            
        })

        return (
            <tr key={this.state.group.group_name}>
                <td>{this.state.group.ldap_group_name}</td>
                <td>                                
                    {this.props.onReturnEditModeHandler() ? (
                        this.renderEditGroupName()
                    ) : (                            
                        this.state.group.group_name
                    )}                                
                </td>
                <td>                            
                    {this.props.onReturnEditModeHandler() ? (
                        this.renderEditLeaseTime()
                    ) : (    
                        this.formatSeconds(this.state.group.lease_time)
                    )}
                </td>
                <td>
                    {this.state.tags}
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
                { this.props.onReturnEditModeHandler() &&
                    <td>                                    
                    <GroupOptions                                    
                        group={this.state.group}
                        onGroupDeleteHandler={(success: boolean, status_code: number) => this.onGroupDeleteHandler(success, status_code) }                          
                        onGroupUpdateHandler={() => this.onGroupUpdateAPIHandler() }   
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