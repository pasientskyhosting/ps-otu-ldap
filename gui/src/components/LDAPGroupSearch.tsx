import React from 'react';
import { HTMLSelect } from "@blueprintjs/core";

import APIService, { ILDAPGroup } from '../services/APIService';

interface IProps {          
}

interface IState {    
    ldap_groups: ILDAPGroup[],
    ldap_group_name: string
}

export default class LDAPGroupSearch extends React.Component<IProps, IState> {
    

    componentWillMount() {
        this.loadData()
    }

    constructor(props: IProps) {
        super(props)
        this.state = {
            ldap_groups: [{ldap_group_name: "None"}],
            ldap_group_name: ""
        }
    }

    private async loadData() {

        const ldap_groups = await APIService.getAllLDAPGroups()              

        if(APIService.success) {
            this.setState({
                ldap_groups
            })
        }
        
        // call login handler
        //this.props.onLDAPGroupsFetchHandler(APIService.success) 

        
    }

    render () {        

        return (
        
            <HTMLSelect                      
                id="ldap-groups"   
                {...this.props}           
            >                    
                {this.state.ldap_groups.map(ldap_group => (
                    <option key={ldap_group.ldap_group_name} value={ldap_group.ldap_group_name}>{ldap_group.ldap_group_name}</option>
                ))}
            </HTMLSelect>           

        )        
        
    }

}