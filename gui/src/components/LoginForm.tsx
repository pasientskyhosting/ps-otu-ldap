import React, { ReactEventHandler } from 'react';
import { Button, Card, Elevation, FormGroup, InputGroup, Intent, Spinner, Callout } from "@blueprintjs/core";

interface IProps {
  onSubmit: (username: string, password: string) => Promise<void>,
  errorMessage?: string
}

interface IState {
  username: string,
  password: string,
  connecting: boolean  
}

type Nullable<T> = T | null

export default class LoginForm extends React.Component<IProps, IState> {
    
    private usernameInputRef: Nullable<HTMLInputElement>
    private passwordInputRef: Nullable<HTMLInputElement>

    constructor (props: IProps) {      
      
      super(props)

      this.usernameInputRef = null
      this.passwordInputRef = null

      this.state = {
        username: "",
        password: "",
        connecting: false
      }
    }

    private async onSubmit() {

      this.setState({
        connecting: true
      })      

      await this.props.onSubmit(this.state.username, this.state.password)

      setTimeout(() => {
        console.log("timeout")

        this.setState({
          connecting: false
        }, () => {
          if (this.props.errorMessage) {
            if(this.state.username) this.passwordInputRef && this.passwordInputRef.focus()
            else this.usernameInputRef && this.usernameInputRef.focus()
          } 
        })

      }, 1500);

    }

    public componentDidMount() {
      this.usernameInputRef && this.usernameInputRef.focus()      
    }

    private renderForm() {

      return (
        <div id="login-form-content">
        {this.props.errorMessage ? <Callout title="Unauthorized" className="login-error-message" intent={Intent.DANGER} >{this.props.errorMessage}</Callout> : null }     
          <FormGroup       
                                                                                 
          >
          <InputGroup                  
            id="text-input"              
            placeholder="Username"  
            value={this.state.username}
            inputRef={(input) => this.usernameInputRef = input}
            large={true}            
            leftIcon="person"              
            onKeyDown={(e) => {                
              if(e.keyCode == 13) this.onSubmit()
            }}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
            this.setState({
            username: e.target.value
            })
          }}                   
          />
          </FormGroup>
          <FormGroup
          >
          <InputGroup                      
            placeholder="Enter your password..."
            inputRef={(input) => this.passwordInputRef = input}
            type="password"                
            value={this.state.password}
            large={true}                
            leftIcon="lock"
            onKeyDown={(e: React.KeyboardEvent) => {                
              if(e.keyCode == 13) this.onSubmit()
            }}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
            this.setState({
            password: e.target.value
            })                  
          }}
          />
          </FormGroup>                           
          <Button                          
            style={{ width: "100%", marginTop: "20px" }}              
            large={true}
            intent={Intent.PRIMARY}
            onClick={() => this.onSubmit()}
          >Login</Button>                           
        </div>
        
      )   

    }

    private renderSpinner(){      
        return (          
          <div id="login-spinner-content">
          <Spinner intent={Intent.PRIMARY} size={170} />          
          </div>
        ) 
    }
    
    private renderCardContent() {      
            
      if(!this.state.connecting) {
        return this.renderForm()            
      } else {    
        return this.renderSpinner()
      }

    }

    render() {

      return (
        <div className="login-form " >          
          <Card interactive={false} elevation={Elevation.THREE}>
            <h1>Sign In</h1>   
            <div id="login-card-content">
            {this.renderCardContent()}
            </div>
          </Card>
          <Card interactive={false} elevation={Elevation.THREE} style={{ backgroundColor: "#394B59", align: "center" }} >
          <p><a href="#">Lost your password?</a></p>
          </Card>
        </div>
      )
    }

}
