import logo from '../assets/images/logo.svg';
import TopBar from "./TopBar";
import Main from "./Main";
import React, { useState } from "react";

import { TOKEN_KEY } from "../constants";
import "../styles/App.css";

//usestate hook
//usestate is a function, which allows func component to use states
//[10, f]: f can setstate for 10. 10: state initial num.
//const[a, setA]
function App() {
  //jie gou
  const [isLoggedIn, setIsLoggedIn] = useState(
      localStorage.getItem(TOKEN_KEY) ? true : false
  );

  //define logout func
  const logout = () => {
    console.log("log out");
    localStorage.removeItem(TOKEN_KEY);
    //=> in class-base: this.setState
    setIsLoggedIn(false);
  };

  //local storage: store token vs cookie store token.
  //token is from backend, renamed to TOKEN_KEY
  const loggedIn = (token) => {
    if (token) {
      localStorage.setItem(TOKEN_KEY, token);
      setIsLoggedIn(true);
    }
  };
  return (
      <div className="App">
        <TopBar isLoggedIn={isLoggedIn}
                handleLogout={logout} />
        <Main
            isLoggedIn = {isLoggedIn}
            handleLoggedIn = {loggedIn}/>
      </div>
  );
}
//handeloggedin is callback func.

export default App;

