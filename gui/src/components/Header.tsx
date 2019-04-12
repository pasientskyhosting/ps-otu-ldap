import React from 'react';
import { Button, Navbar, Alignment } from "@blueprintjs/core";

interface IProps {
  displayName: string
}

const Header: React.FunctionComponent<IProps> = ({ displayName }) => {
    return (         
        <Navbar>           
            <Navbar.Group align={Alignment.LEFT}>                
              <div className="logo-white"></div>
              <Navbar.Heading >Patientsky One Time User</Navbar.Heading>
              <Navbar.Divider />                                            
            </Navbar.Group>
            <Navbar.Group align={Alignment.RIGHT}>                
                {displayName ? 
                <>
                <Button className="bp3-minimal" icon="person" text={displayName} />
                <Button className="bp3-minimal" icon="cog" text="" />                
                <Button className="bp3-minimal" icon="log-out" text="" onClick={ () => {
                  localStorage.removeItem('jwt.token')
                  location.href = "/"
                }} /></> : null }
              
            </Navbar.Group>
        </Navbar>              
      )
}

export default Header