import React from 'react';
import { HTMLSelect, IHTMLSelectProps } from "@blueprintjs/core";

import APIService, { ILDAPGroup } from '../services/APIService';

interface IState {    
    ldap_groups: ILDAPGroup[],
    ldap_group_name: string
}

export default class LDAPGroupSearch extends React.Component<IHTMLSelectProps, IState> {
    

    componentWillMount() {
        this.loadData()
    }

    constructor(props: IHTMLSelectProps) {
        super(props)
        this.state = {
            ldap_groups: [],
            ldap_group_name: ""
        }
    }

    private async loadData() {

        const ldap_groups = await APIService.getAllLDAPGroups()              

        ldap_groups.unshift({ldap_group_name: "Select a LDAP group"})

        if(APIService.success) {
            this.setState({
                ldap_groups
            })
        }
        
    }

    render () {        

        return (
        
            <HTMLSelect                
                {...this.props}           
                options={this.state.ldap_groups.map(ldap_group => {
                    return { value: ldap_group.ldap_group_name, label: ldap_group.ldap_group_name }
                })}
            >   
            </HTMLSelect>           

        )        
        
    }

}