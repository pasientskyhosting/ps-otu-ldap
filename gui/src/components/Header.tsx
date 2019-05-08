import React from 'react';
import { Button, Navbar, Alignment, Intent } from "@blueprintjs/core";
import { AppToaster } from '../services/Toaster';

interface IProps {
  displayName: string
}

const Header: React.FunctionComponent<IProps> = ({ displayName }) => {
    return (         
        <Navbar>           
            <Navbar.Group align={Alignment.LEFT}>                
              <div className="logo-white"></div>
              <Navbar.Heading >PatientSky One Time User</Navbar.Heading>
              <Navbar.Divider />                                            
            </Navbar.Group>
            <Navbar.Group align={Alignment.RIGHT}>                
                {displayName ? 
                <>
                <Button className="bp3-minimal" icon="person" text={displayName} />
                <Button 
                  className="bp3-minimal" 
                  icon="book" 
                  text="Docs"
                  onClick={ () => {                                        
                    window.open(
                      '/docs/',
                      '_blank'
                    );
                  }}
                 />
                <Button 
                  className="bp3-minimal" 
                  icon="log-out" 
                  text="" 
                  onClick={ () => {
                  localStorage.removeItem('jwt.token')
                  AppToaster.show(
                    {
                        intent: Intent.WARNING, 
                        message: "Goodbye old friend!",
                        icon: "hand" 
                    }
                  )

                  setTimeout(() => {
                    location.href = "/"
                  }, 2000)
                  
                }} /></> : null }
              
            </Navbar.Group>
        </Navbar>              
      )
}

export default Header