// Import styles
import "./styles/app.less"

// Components
import React from 'react';
import LoginForm from './components/LoginForm'
import Header from './components/Header'
import GroupCreate from './components/GroupCreate'
import { IProps, Intent } from "@blueprintjs/core";
import GroupList from "./components/GroupList";
import { AppToaster } from "./services/Toaster";

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
        
    private groupListRef: GroupList | null

    constructor (props: IProps) {      

        super(props)
        
        this.groupListRef = null

        this.state = {
            loginFailed: false,
            isVerified: false,       
        }

    } 

    public componentWillMount() {
        
        this.updateAuth()
        
    }

    private updateAuth() {

        let tokenPayload = this.getTokenPayload()
        let isVerified = false        

        if (tokenPayload) {
            isVerified = true
        }

        this.setState({
            tokenPayload, isVerified
        })
    }

    private getTokenPayload(): ITokenPayload | undefined {

        let tokenPayload = undefined

        const token = localStorage.getItem('jwt.token')    

        if(token) {            
            tokenPayload = JSON.parse(atob(token.split(".")[1]))   
        }

        return tokenPayload
    }

    private renderHeader() {

        return (            
            <Header displayName={ this.state.tokenPayload ? this.state.tokenPayload.display_name : "" } />            
        )
    }

    private onGroupCreateHandler(success: boolean) {

        if (success) {
            AppToaster.show({
                intent: Intent.SUCCESS, 
                message: "Group created successfully."
            })
        } else {
            AppToaster.show({
                intent: Intent.DANGER, 
                message: "An error occured while creating group" 
            })
        }       

        this.onFinishedHandler()
        
        if(success && this.groupListRef) {
            this.groupListRef.fetch() 
        }
        
    }

    private onGroupUpdateHandler(success: boolean) {

        if (success) {
            AppToaster.show({
                intent: Intent.SUCCESS, 
                message: "Group updated successfully."
            })
        } else {
            AppToaster.show({
                intent: Intent.DANGER, 
                message: "An error occured while updating group" 
            })
        }       

        this.onFinishedHandler()
        
        if(success && this.groupListRef) {
            this.groupListRef.fetch() 
        }
        
    }

    private onFinishedHandler() {
        this.updateAuth()        
    }  

    private renderSearchGroup() {
        return (
            <GroupList
                is_admin={this.state.tokenPayload ? this.state.tokenPayload.is_admin : false } 
                ref={(ref) => this.groupListRef = ref} 
                onGroupsFetchHandler={ (success: boolean, status_code: number) => this.onFinishedHandler() }
                onGroupUpdateHandler={ (success: boolean) => this.onGroupUpdateHandler(success) }
            />
        )
    }

    private renderCreateGroup() {
        return (
            <GroupCreate onGroupCreateHandler={(success: boolean) => this.onGroupCreateHandler(success) } />
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
                <LoginForm onLoginHandler={(success: boolean, status_code: number) => this.onFinishedHandler() } />
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
