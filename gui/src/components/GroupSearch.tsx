import React from 'react';
import { InputGroup, HTMLTable, Icon, Card, Elevation } from "@blueprintjs/core";
import UserOptions from './UserOptions'
import GroupOptions from './GroupOptions'
import APIService, { IGroup, IUser } from '../services/APIService';

interface IProps {  
  onGroupsFetchHandler: (success: boolean, status_code: number) => void,  
}

interface IState {    
    groups: IGroup[]
    users: IUser[],
    connecting: boolean,
    errorMessage?: string
}

export default class GroupSearch extends React.Component<IProps, IState> {
           
    constructor(props: IProps) {

        super(props)

        this.state = {
            groups: [],
            users: [],
            connecting: false
        }
    }

    componentWillMount() {
        this.loadGroupData()
    }

    public fetch() {
        this.loadGroupData()
    }

    private async loadGroupData() {

        this.setState({
            connecting: true
        })

        const groups = await APIService.getAllGroups()
        
        this.setState({
            connecting: false
        })          

        if(!APIService.success) {
            this.setState({
                errorMessage: "Error while fetching groups"
            })
        } else {    
            this.setState({
                groups
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
                    <tr><td>Group name</td><td>Lease time</td><td>User Options</td><td>Group Options</td></tr>
                </thead>
                <tbody>                
                {this.state.groups.map((group: IGroup) => {
                    return <tr key={group.group_name}><td>{group.group_name}</td><td>{this.formatSeconds(group.lease_time)}</td><td>{<UserOptions />}</td><td>{<GroupOptions onGroupDeleteHandler={(success: boolean, status_code: number) => this.onGroupDeleteHandler(success, status_code) } groupName={group.group_name} />}</td></tr>
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