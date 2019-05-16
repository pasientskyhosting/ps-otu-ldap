import React from 'react';
import { Button, Navbar, Alignment, Intent, Popover, Menu, MenuItem, MenuDivider, Position } from "@blueprintjs/core";
import { AppToaster } from '../services/Toaster';

interface IProps {
  displayName: string
}

const subDocMenu = (
  <Menu>
    <MenuItem 
        icon="book" 
        text="README" 
        onClick={ () => {                                        
          window.open(
            'https://github.com/pasientskyhosting/ps-otu-ldap',
            '_blank'
          );
        }}
      />    
      <MenuItem 
        icon="graph" 
        text="API" 
        onClick={ () => {                                        
          window.open(            
            '/docs/',
            '_blank'
          );
        }}
      />      
  </Menu>
);

const Header: React.FunctionComponent<IProps> = ({ displayName }) => {
    return (         
        <Navbar>           
            <Navbar.Group align={Alignment.LEFT}>                
              <div 
                className="logo-white pointer"
                onClick={ () => {                                        
                  location.href = "/"
                }}
              ></div>
              <Navbar.Heading 
                className="pointer"
                onClick={ () => {                                        
                  location.href = "/"
                }}
                >PatientSky One Time User
              </Navbar.Heading>
              <Navbar.Divider />                                            
            </Navbar.Group>            
            <Navbar.Group align={Alignment.RIGHT}>                
                {displayName ? 
                <>
                <Button className="bp3-minimal" icon="person" text={displayName} />
                <Popover content={subDocMenu} position={Position.RIGHT_BOTTOM}>
                    <Button icon="book" text="Docs" className="bp3-minimal"  />
                </Popover>         
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