import React, { ButtonHTMLAttributes } from 'react';
import { HTMLTable, Card, Elevation, Intent, Button, HTMLSelect, Tag, ButtonGroup, Popover, Position, Classes, Switch, Navbar, Alignment } from "@blueprintjs/core";
import UserOptions from './UserOptions'
import APIService, { IGroup, IUser } from '../services/APIService';
import GroupOptions from './GroupOptions';
import ValidatedInputGroup from './ValidatedInputGroup';
import { AppToaster } from '../services/Toaster';

interface IProps {  
  onGroupsFetchHandler: (success: boolean, status_code: number) => void,  
  is_admin: boolean
}

interface IState {    
    groups: IGroup[]
    users: IUser[],
    errorMessage?: string
    groupNameChanging: string
    leaseTimeChanging: number
    editMode: boolean
}

export default class GroupList extends React.Component<IProps, IState> {
           
    constructor(props: IProps) {

        super(props)

        this.state = {
            groups: [],
            users: [],
            groupNameChanging: "",
            leaseTimeChanging: 0,
            editMode: false
        }
    }

    componentWillMount() {
        this.loadData()
    }

    public fetch() {
        this.loadData()
    }

    private async loadData() {

        const groups = await APIService.getAllGroups()        
        const users = await APIService.getAllActiveUsers()

        if(APIService.success) {
            this.setState({
                groups, users
            })
        }
        
        // call login handler
        this.props.onGroupsFetchHandler(APIService.success, APIService.status) 

        
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
        
        AppToaster.show(
            {
                intent: Intent.SUCCESS, 
                message: "Group deleted successfully." 
            }
        )

        this.fetch()

    }

    private renderLeaseTimeUpdateField(lease_time: number) {

        return (

            <ButtonGroup>
            
            <HTMLSelect                                                  
                value={lease_time}                        
                onChange={(e: React.ChangeEvent<HTMLSelectElement>) => {                            
                    this.setState({
                        leaseTimeChanging: parseInt(e.target.value)
                    })
                }}
            >                    
            <option value="60">1 hour</option>
            <option value="720">12 hours</option>
            <option value="1440">24 hours</option>
            <option value="10080">1 week</option>
            <option value="20160">2 weeks</option>
            </HTMLSelect>
            <Button small={true} icon="tick" text="" className="bp3-minimal" onClick={ (e: React.MouseEvent<HTMLElement, MouseEvent>) => {
                this.setState({
                    leaseTimeChanging: 0
                })
            }} />
            <Button small={true} icon="cross" text="" className="bp3-minimal" onClick={ (e: React.MouseEvent<HTMLElement, MouseEvent>) => {
                this.setState({
                    leaseTimeChanging: 0
                })
            }} />
            </ButtonGroup>

        )

    }

    private renderGroupUpdateField(group_name: string) {

        return (

            <ButtonGroup>
            <ValidatedInputGroup                
                defaultValue={group_name}
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
                        alert("Call API - Changing " + group_name + " to " + e.target.value)
                        this.setState({groupNameChanging: ""})
                    }
                    if (e.keyCode === 27) {
                        alert("Cancel - revert to original " + group_name)
                        this.setState({groupNameChanging: ""})
                    }
                }}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => {}}
            />
            <Button small={true} icon="tick" text="" className="bp3-minimal" onClick={ (e: React.MouseEvent<HTMLElement, MouseEvent>) => {
                this.setState({
                    groupNameChanging: ""
                })
            }} />
            <Button small={true} icon="cross" text="" className="bp3-minimal" onClick={ (e: React.MouseEvent<HTMLElement, MouseEvent>) => {
                this.setState({
                    groupNameChanging: ""
                })
            }} />
            </ButtonGroup>

        )

    }

    private renderStaticLeaseTimeField(lease_time: number) {

        return (    
            <>
            {this.formatSeconds(lease_time)}
            {this.state.editMode && 
                <Button 
                    className="bp3-minimal"                                         
                    icon="edit" 
                    small={false}
                    onClick={ (e: React.MouseEvent<HTMLElement, MouseEvent>) => {                                            
                        this.setState({
                            leaseTimeChanging: lease_time
                        })
                    }}
                />
            }
            </>     
        )                        

    }


    private renderStaticGroupNameField(group_name: string) {

        return (    
            <>
            {group_name}
            {this.state.editMode && 
                <Button 
                    className="bp3-minimal"                                         
                    icon="edit" 
                    small={false} 
                    onClick={ (e: React.MouseEvent<HTMLElement, MouseEvent>) => {                                            
                        this.setState({
                            groupNameChanging: group_name
                        })
                    }}
                />
            }
            </>     
        )                        

    }

    private renderGroupTable () {
        
        return (

            <HTMLTable            
            bordered={true}
            striped={true}
            interactive={true}
            >
                <thead>
                    <tr>
                        <td>LDAP group</td>
                        <td>Group name</td>
                        <td>Lease time</td>
                        <td>Custom props</td>
                        <td>User Options</td>
                        { (this.props.is_admin && this.state.editMode) &&
                            <td>Group Options</td>
                        }
                    </tr>
                </thead>
                <tbody>                
                {this.state.groups.map((group: IGroup) => {

                    let username = ""
                    let password = ""
                    var groupNameChange = (this.state.groupNameChanging == group.group_name) ? true : false;
                    var leaseTimeChange = (this.state.leaseTimeChanging == group.lease_time) ? true : false;

                    this.state.users.map((user: IUser) => {
                        
                        if (group.group_name == user.group_name) {
                            username = user.username
                            password = user.password
                        }
                    })

                    

                    return (
                        <tr key={group.group_name}>
                            <td>{group.ldap_group_name}</td>
                            <td>                                
                                {groupNameChange && this.props.is_admin ? (
                                    this.renderGroupUpdateField(group.group_name)
                                ) : (    
                                    this.renderStaticGroupNameField(group.group_name)
                                )}                                
                            </td>
                            <td>                            
                                {leaseTimeChange && this.props.is_admin ? (
                                    this.renderLeaseTimeUpdateField(group.lease_time)
                                ) : (    
                                    this.renderStaticLeaseTimeField(group.lease_time)
                                )}
                            </td>
                            <td>

                                <Tag
                                    className="tags"
                                    key="yes2s24"                                                                                                                                                
                                    minimal={true}         
                                    round={true}                                    
                                    onRemove={ (e: React.MouseEvent<HTMLButtonElement>) => {
                                        alert(e)
                                    }}                          
                                    
                                >
                                    readonly=true
                                </Tag>                                

                                <Tag
                                    className="tags"
                                    key="yess"                                                                                                                                                
                                    minimal={true}             
                                    round={true}             
                                    onRemove={ (e: React.MouseEvent<HTMLButtonElement>) => {
                                        alert(e)
                                    }}                                                
                                >
                                    key=value
                                </Tag>                                
                                <Tag
                                    className="tags"
                                    key="yess2"                                                                                                                                                
                                    minimal={true}             
                                    round={true}             
                                    onRemove={ (e: React.MouseEvent<HTMLButtonElement>) => {
                                        alert(e)
                                    }}                                                
                                >
                                    sql_granta_access_onusers=GRANT ALL
                                </Tag>

                            </td>
                            <td>
                                {
                                <UserOptions
                                    username={username} 
                                    password={password}
                                    group_name={group.group_name}
                                />
                                }
                            </td>        
                            { (this.props.is_admin && this.state.editMode) &&
                                <td>                                    
                                <GroupOptions                                    
                                    group_name={group.group_name}
                                    onGroupDeleteHandler={(success: boolean, status_code: number) => this.onGroupDeleteHandler(success, status_code) }                                   
                                />
                                </td>
                            }                    
                        </tr>
                    )
                })}
                </tbody>
            </HTMLTable>
        )
    }    
    
    render () {

        return (       
            <div className="groups-search card">
            <Card interactive={true} elevation={Elevation.FOUR}>
            <div className="edit-switch">
            <Switch 
                onChange={ (e: React.FormEvent<HTMLInputElement>) => {
                    this.setState({
                        editMode: e.target.checked
                    })
                }} 
                labelElement={<strong>Enable edit mode</strong>} />
            </div>
                <div className="groups-search-content">
                    <div><h2>Groups</h2></div>
                    <div className="groups-search-table">                         
                        {this.renderGroupTable()}                    
                    </div>     
                </div>            
            </Card>
            </div>
        )

    }

}   