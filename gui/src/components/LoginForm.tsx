import React from 'react';
import { Button, Card, Elevation, FormGroup, InputGroup, Intent } from "@blueprintjs/core";

export default class LoginForm extends React.Component {
  
    render() {        

      return (
        <div style={{ margin: "auto", width: "400px"}}>
          <Card interactive={true} elevation={Elevation.ONE}>
          <h3><a href="#">Login</a></h3>
                    <p>
                        Please provide your LDAP credentials
                    </p>
            <FormGroup                                
              labelFor="text-input"                              
            >
              <InputGroup                  
                id="text-input" 
                placeholder="Username"  
                large={true}                
              />
            </FormGroup>
            <FormGroup                                
              labelFor="text-input"                
            >
              <InputGroup                      
                placeholder="Enter your password..."            
                type="password"                
                large={true}                
              />
            </FormGroup>
            <Button            
              style={{ width: "100px"}}
              icon="lock"
              large={true}
              intent={Intent.PRIMARY}
            >Login</Button>
          </Card>
        </div>
      )
    }

}