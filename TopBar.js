import React from 'react';
import logo from "../assets/images/logo.svg";

import { LogoutOutlined } from "@ant-design/icons";

//function component
//use hooks under function component
function TopBar(props) {
    //props ->isLoggedIn
    const { isLoggedIn, handleLogout } = props;
    return (
        <header className="App-header">
            <img src={logo} className="App-logo" alt="logo" />
            <span className="App-title">Around Web</span>
            {
                isLoggedIn ?
                    <LogoutOutlined className='logout' onClick={handleLogout}/>
                    :
                    null
            }
        </header>
    );
}

export default TopBar;
