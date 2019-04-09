import React from 'react';import { Button, Navbar, Alignment } from "@blueprintjs/core";

export default class LoginForm extends React.Component {
  
    render() {        

      return (         
        <Navbar>
            <Navbar.Group align={Alignment.LEFT}>                
                <Navbar.Heading >Patientsky One Time User</Navbar.Heading>
                <Navbar.Divider />                
                <Button className="bp3-minimal" icon="home" text="Home" />
            </Navbar.Group>
            <Navbar.Group align={Alignment.RIGHT}>
                <Button className="bp3-minimal" icon="cog" text="" />
                <Button className="bp3-minimal" icon="person" text="Jeppe Larsen" />
            </Navbar.Group>
        </Navbar>              
      )
    }

}