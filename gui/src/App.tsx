// Import styles
import "./styles/app.less"

// Components
import React from 'react';
import LoginForm from './components/LoginForm'
import Header from './components/Header'
import GroupCreate from './components/GroupCreate'
import { IProps } from "@blueprintjs/core";
import GroupSearch from "./components/GroupSearch";

interface IState {
    loginFailed: boolean
    isVerified: boolean
    tokenPayload?: ITokenPayload
    token?: string
    errorMesssage?: string    
}

interface ITokenPayload {
    is_admin: boolean
    user_id: string
    exp: number
    display_name: string
}

class App extends React.Component<{}, IState>{    
        
    private groupSearchRef: GroupSearch | null

    constructor (props: IProps) {      

        super(props)
        
        this.groupSearchRef = null

        this.state = {
            loginFailed: false,
            isVerified: false            
        }

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
        })
        
    }

    private renderHeader() {

        return (            
            <Header displayName={ this.state.tokenPayload ? this.state.tokenPayload.display_name : "" } />            
        )
    }

    private onGroupCreateHandler(success: boolean, status_code: number) {

        this.onFinishedHandler(success, status_code)
        
        if(success && this.groupSearchRef) this.groupSearchRef.fetch()
        
    }

    private onFinishedHandler(success: boolean, status_code: number) {
        
        console.log("onFinishedHandler success: " + success)
        console.log("onFinishedHandler status_code: " + status_code)

        if (success) {                    

            this.setState({
                isVerified: true,
                loginFailed: false
            })            

        } else {           
            
            let isVerified = false
            let loginFailed = false

            if(status_code != 401 ) {
                isVerified = true,           
                loginFailed = true
            }

            this.setState({
                isVerified, loginFailed
            })            

        }
        
    }  

    private renderSearchGroup() {
        return (
            <GroupSearch ref={(ref) => this.groupSearchRef = ref} onGroupsFetchHandler={ (success: boolean, status_code: number) => this.onFinishedHandler(success, status_code) } />
        )
    }

    private renderCreateGroup() {
        return (
            <GroupCreate onGroupCreateHandler={(success: boolean, status_code: number) => this.onGroupCreateHandler(success, status_code) } />
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
                <LoginForm onLoginHandler={(success: boolean, status_code: number) => this.onFinishedHandler(success, status_code) } />
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
