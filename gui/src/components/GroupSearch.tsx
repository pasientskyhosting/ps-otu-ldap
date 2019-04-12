import React from 'react';
import { InputGroup, HTMLTable, Icon, Card, Elevation } from "@blueprintjs/core";

export interface IGroup {
    group_name: string
    ldap_group_name: string
    lease_time: number
    create_time: number
    create_by: string
}

export interface IUser {
    username: string
    password: string
    group_name: string
    expire_time: number
    create_time: number
    create_by: string
}

interface IProps {
  groups: IGroup[]
  users: IUser[]
}

interface IState {
    filteredGroups: IGroup[]
}

export default class GroupSearch extends React.Component<IProps, IState> {
           

    constructor (props: IProps) {      
      
      super(props)

      this.state = {
          filteredGroups: props.groups          
      }
    
    }

    private getFilteredList() {
        
        return this.props.groups.filter((group) => {
            return group.group_name.includes("")
        })
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

        var days = Math.floor(seconds / (3600*24))        
        seconds  -= days*3600*24

        var hrs   = Math.floor(seconds / 3600)
        seconds  -= hrs*3600        
        
        if (days === 0) {
            return hrs+" hour(s)"
        } else {
            return days+" day(s)"
        }

    }


    private renderGroupTable () {

        return (                
            <HTMLTable            
            bordered={true}
            striped={true}
            interactive={true}
            >
                <thead>
                    <tr><td>Group name</td><td>LDAP group</td><td>Lease time</td><td>User</td></tr>
                </thead>
                <tbody>                
                {this.props.groups.map((group: IGroup) => {
                    return <tr key={group.group_name}><td>{group.group_name}</td><td>{group.ldap_group_name}</td><td>{this.formatSeconds(group.lease_time)}</td><td>test</td></tr>
                })}
                </tbody>
            </HTMLTable>
        )
    }

    render () {

        return (       
            <div className="groups-search card">
            <Card interactive={true} elevation={Elevation.FOUR}>
                <div className="groups-search-content">
                    <div className="groups-search-bar search">
                        {this.renderSearch()}
                    </div>     
                    <div className="groups-search-table">
                        {this.renderGroupTable()}                    
                    </div>     
                </div>            
            </Card>
            </div>
        )
    }

}