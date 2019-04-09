// Import styles
import "./styles/app.less"

// Components
import React from 'react';
import LoginForm from './components/LoginForm'
import Header from './components/Header'

interface IProps {
    isAuthenticated: boolean
}

class App extends React.Component<IProps>{    

    private renderContent() {        
        if (this.props.isAuthenticated ) {
            return (
                // top div
                <div>yes</div>
            )
        } else {
            return (
                <LoginForm onSubmit={async(username, password) => {
                    console.log(username,password)                    
                }} />
            )
        }
    }

    render() { 
        return (  
            <div className="wrapper bp3-dark">
                <div className="header">
                    <Header />     
                </div>                                
                <div className="content">
                    {this.renderContent()}
                </div>
            </div>            
        )
    }

}

export default App
