// Import styles
import "./styles/app.less"

// Components
import React from 'react';
import LoginForm from './components/LoginForm'
import Header from './components/Header'
import GroupSearch, { IGroup, IUser } from './components/GroupSearch'
import GroupCreate from './components/GroupCreate'
import { IProps } from "@blueprintjs/core";

interface IState {
    loginFailed: boolean
    isVerified: boolean
    tokenPayload?: ITokenPayload
    token?: string
    errorMesssage?: string
    groups: IGroup[]
    users: IUser[]
}

interface ITokenPayload {
    is_admin: boolean
    user_id: string
    exp: number
    display_name: string
}

class App extends React.Component<{}, IState>{    
        
    constructor (props: IProps) {      

        super(props)
        
        this.state = {
            loginFailed: false,
            isVerified: false,
            groups: [],
            users: []
        }

    } 
    
    private async loadGroupData() {
        
        let response = await fetch('/v1/api/groups', {
            method: 'get',                       
            headers: { "Authorization": `Bearer ${this.state.token}`  }             
        })

        this.setState({
            groups: await response.json()
        })
    }

    private async loadUserData() {
        
        let response = await fetch('/v1/api/users', {
            method: 'get',                       
            headers: { "Authorization": `Bearer ${this.state.token}`  }             
        })

        this.setState({
            users: await response.json()
        })
    }

    public componentWillMount() {
        
        let isVerified = false
        let tokenPayload = undefined

        const token = localStorage.getItem('jwt.token')

        if(!token) {
            isVerified = false
        } else {            
            tokenPayload = JSON.parse(atob(token.split(".")[1]))            
            isVerified = true
        }

        this.setState({
            isVerified,
            tokenPayload: tokenPayload || undefined,
            token: token || undefined            
        }, () => {
            if (token) {
                this.loadGroupData()
                this.loadUserData()
            }
        })
        
    }

    private renderHeader() {
        return (            
            <Header displayName={ this.state.tokenPayload ? this.state.tokenPayload.display_name : "" } />            
        )
    }

    private async onLogin(username: string, password: string) {

        try {

            let response = await fetch('/v1/api/auth/authorize', {
                method: 'post',                        
                body: JSON.stringify({  
                    username, password
                })
            })

            if (response.status !== 200) {
                this.setState({                                
                    loginFailed: true,
                    errorMesssage: "Wrong username or password"
                })

            } else {

                let responseParsed = await response.json()
                localStorage.setItem('jwt.token', responseParsed.token)
                                   
                let tokenPayload = JSON.parse(atob(responseParsed.token.split(".")[1]))                            

                this.setState({
                    isVerified: true,
                    errorMesssage: "",
                    tokenPayload,
                    token: responseParsed.token                   
                }, () => {
                    this.loadGroupData()
                })
            }                        

        } catch (error)
        {
            this.setState({                                
                loginFailed: true,
                errorMesssage: "Wrong username or password"
            })
            console.log(error)
        }

    }

    private async onGroupCreate(group_name: string, lease_time: number) {

        try {
                      
            let ldap_group_name = group_name

            let response = await fetch('/v1/api/groups', {
                method: 'post',    
                headers: { "Authorization": `Bearer ${this.state.token}`  },                   
                body: JSON.stringify({  
                    ldap_group_name, group_name, lease_time
                })
            })

            if (response.status !== 201) {

                switch(response.status) {
                    case 401:
                        this.setState({                                
                            isVerified: false,
                            loginFailed: false,                                                
                        })
                    case 400:
                        this.setState({                                
                            errorMesssage: "Bad request error",                            
                        })
                    break
                }                

            } else {                                    
                
                this.setState({
                    isVerified: true,
                    errorMesssage: "",
                }, () => {
                    this.loadGroupData()
                })
            }                        

        } catch (error)
        {
            this.setState({                                
                loginFailed: true,
                errorMesssage: "Wrong username or password"
            })
            console.log(error)
        }

    }

    private renderSearchGroup() {
        return (
            <GroupSearch groups={this.state.groups} users={this.state.users} />
        )
    }

    private renderCreateGroup() {
        return (
            <GroupCreate errorMessage={this.state.errorMesssage} onSubmit={(group_name, lease_time) => this.onGroupCreate(group_name, lease_time)} />
        )
    }

    private renderContent() {
        
        if ( this.state.isVerified == true ) {

            if (this.state.tokenPayload && this.state.tokenPayload.is_admin === true ) {    
                return (
                    <>
                        {this.renderCreateGroup()}
                        {this.renderSearchGroup()}                        
                    </>
                )            
            } else {
                return (                    
                    this.renderSearchGroup()                    
                )
            }
            
        } else {
            return (                
                <LoginForm errorMessage={this.state.errorMesssage} onSubmit={(username, password) => this.onLogin(username, password)} />                
            )
        }
    }    

    render() { 
                
        return (  
            <div className="wrapper bp3-dark">      
                <div className="header">
                    {this.renderHeader()}                                          
                </div>                  
                <div className="content">
                    {this.renderContent()}                
                </div>            
            </div>            
        )
    }

}

export default App
