import React from 'react';
import { Button, Card, Elevation, FormGroup, InputGroup, Intent, Text} from "@blueprintjs/core";

interface IProps {
  onSubmit: (username: string, password: string) => Promise<void>
}

interface IState {
  username: string,
  password: string
}

export default class LoginForm extends React.Component<IProps, IState> {
  
    constructor (props) {      
      super(props)
      this.state = {
        username: "",
        password: ""
      }
    }

    private onButtonClick() {
      this.props.onSubmit(this.state.username, this.state.password)
    }

    render() {

      return (
        <div className="login-form" >         

          <Card interactive={false} elevation={Elevation.THREE}>
          
          {/* <div className="login-logo"></div> */}
          <h1>Sign In</h1>          
            <FormGroup                                                                            
            >
              <InputGroup                  
                id="text-input" 
                placeholder="Username"  
                large={true}            
                leftIcon="person"
                onChange={(e) => {
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
                type="password"                
                large={true}                
                leftIcon="lock"
                onChange={(e) => {
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
              onClick={() => this.onButtonClick()}
            >Login</Button>
          </Card>
          <Card interactive={false} elevation={Elevation.THREE} style={{ backgroundColor: "#394B59", align: "center" }} >
          <p><a href="#">Lost your password?</a></p>
          </Card>
        </div>
      )
    }

}
