import React from "react";
import logo from "../assets/images/logo.svg";
import { WEB_NAME } from "../constants";
import { LogoutOutlined } from "@ant-design/icons";

function TopBar(props) {
  const { isLoggedIn, handleLogout } = props;
  return (
    <header className="App-header">
      <img src={logo} className="App-logo" alt="logo" />
      <span className="App-title">{WEB_NAME}</span>
      {isLoggedIn ? (
        <LogoutOutlined className="logout" onClick={handleLogout} />
      ) : null}
    </header>
  );
}

export default TopBar;