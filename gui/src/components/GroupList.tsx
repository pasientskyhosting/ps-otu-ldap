import React, { ButtonHTMLAttributes } from 'react';
import { HTMLTable, Card, Elevation, Intent, Button, HTMLSelect, Tag, ButtonGroup, Popover, Position, Classes, Switch, Navbar, Alignment, TagInput } from "@blueprintjs/core";
import APIService, { IGroup, IUser, IGroupCustomProps } from '../services/APIService';
import GroupEntry from './GroupEntry';
import { UNDERLINE } from '@blueprintjs/icons/lib/esm/generated/iconContents';
import { AppToaster } from '../services/Toaster';

interface IProps {  
  onGroupsFetchHandler: (success: boolean, status_code: number) => void,  
  onGroupUpdateHandler: (success: boolean) => void,  
  is_admin: boolean
}

interface IState {    
    groups: IGroup[]
    users: IUser[],
    errorMessage?: string    
    editMode: boolean
}

export default class GroupList extends React.Component<IProps, IState> {
           
    constructor(props: IProps) {

        super(props)

        this.state = {
            groups: [],
            users: [],            
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
                            <td>&nbsp;</td>
                        }
                    </tr>
                </thead>
                <tbody>                
                {this.state.groups.map((group: IGroup) => {

                    let user: IUser | undefined

                    // find matching user
                    this.state.users.map((userMap: IUser) => {
                        if (group.group_name == userMap.group_name) {                                                        
                            user = userMap                            
                        }
                    })       
                    
                    return (
                       <GroupEntry
                        key={group.group_name}                
                        group={group}       
                        user={user}
                        onGroupDeleteHandler={(success: boolean, status_code: number) => {
                            this.onGroupDeleteHandler(success, status_code)
                        }}
                        onGroupUpdateHandler={(success: boolean, status_code: number) => {
                            this.onGroupUpdateHandler(success, status_code)
                        }}                 
                        editMode={this.state.editMode}
                       >
                       </GroupEntry>
                    )
                })}
                </tbody>
            </HTMLTable>
        )
    }    

    private onGroupUpdateHandler(success: boolean, status_code: number) {
        
        AppToaster.show(
            {
                intent: Intent.SUCCESS, 
                message: "Group updated successfully." 
            }
        )

        this.fetch()

        this.setState({
            editMode: false
        })

    }

    private onGroupDeleteHandler(success: boolean, status_code: number) {
        
        AppToaster.show(
            {
                intent: Intent.SUCCESS, 
                message: "Group deleted successfully." 
            }
        )

        this.loadData()

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
                checked={this.state.editMode}
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