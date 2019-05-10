import React from 'react';
import { InputGroup, HTMLTable, Icon, Card, Elevation, Intent } from "@blueprintjs/core";
import UserOptions from './UserOptions'
import GroupOptions from './GroupOptions'
import APIService, { IGroup, IUser } from '../services/APIService';
import { AppToaster } from "../services/Toaster";

interface IProps {  
  onGroupsFetchHandler: (success: boolean, status_code: number) => void,
  is_admin?: boolean
}

interface IState {    
    groups: IGroup[]
    users: IUser[],
    errorMessage?: string
}

export default class GroupList extends React.Component<IProps, IState> {
           
    constructor(props: IProps) {

        super(props)

        this.state = {
            groups: [],
            users: [],
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
    
    private renderSearch () {

        return (
            <InputGroup                      
            placeholder="Search for LDAP groups..."            
            type="search"                            
            large={false}                
            leftIcon="search"           
          />
        )

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
                        <td>User Options</td>
                    </tr>
                </thead>
                <tbody>                
                {this.state.groups.map((group: IGroup) => {

                    let username = ""
                    let password = ""

                    this.state.users.map((user: IUser) => {
                        
                        if (group.group_name == user.group_name) {
                            username = user.username
                            password = user.password
                        }
                    })

                    return (
                        <tr key={group.group_name}>
                            <td>{group.ldap_group_name}</td>
                            <td>{group.group_name}</td>
                            <td>{this.formatSeconds(group.lease_time)}</td>
                            <td>
                                {
                                <UserOptions
                                    username={username} 
                                    password={password}
                                    group_name={group.group_name}
                                />
                                }
                            </td>                            
                        </tr>
                    )
                })}
                </tbody>
            </HTMLTable>
        )
    }

    private renderAdminGroupTable () {

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
                        <td>One Time User</td>                 
                        <td>Group Options</td>
                    </tr>
                </thead>
                <tbody>                
                {this.state.groups.map((group: IGroup) => {

                    let username = ""
                    let password = ""

                    this.state.users.map((user: IUser) => {
                        
                        if (group.group_name == user.group_name) {
                            username = user.username
                            password = user.password
                        }
                    })

                    return (
                        <tr key={group.group_name}>
                            <td>{group.ldap_group_name}</td>
                            <td>{group.group_name}</td>
                            <td>{this.formatSeconds(group.lease_time)}</td>
                            <td>
                                {
                                <UserOptions
                                    username={username} 
                                    password={password}
                                    group_name={group.group_name}
                                />
                                }
                            </td>
                            <td>
                                {
                                <GroupOptions 
                                onGroupDeleteHandler={(success: boolean, status_code: number) => this.onGroupDeleteHandler(success, status_code) } 
                                group_name={group.group_name} 
                                />
                                }
                            </td>
                        </tr>
                    )
                })}
                </tbody>
            </HTMLTable>
        )
    }

    render () {

        if(this.props.is_admin) {

            return (       
                <div className="groups-search card">
                
                <Card interactive={true} elevation={Elevation.FOUR}>
                    <div className="groups-search-content">
                        <h2>LDAP Groups</h2>
                        <div className="groups-search-table">                         
                            {this.renderAdminGroupTable()}                    
                        </div>     
                    </div>            
                </Card>
                </div>
            )

        } else {

            return (       
                <div className="groups-search card">
                
                <Card interactive={true} elevation={Elevation.FOUR}>
                    <div className="groups-search-content">
                        <h2>LDAP Groups</h2>
                        <div className="groups-search-table">                         
                            {this.renderGroupTable()}                    
                        </div>     
                    </div>            
                </Card>
                </div>
            )
    
        }

    }

}   