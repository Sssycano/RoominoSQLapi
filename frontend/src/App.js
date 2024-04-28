import React, { useState, useEffect } from "react";
import { Routes, Route, Link } from "react-router-dom";
import "bootstrap/dist/css/bootstrap.min.css";
import "./App.css";

import AuthService from "./services/auth.service";
import UnitInfo from "./components/UnitInfo"; 
import UpdatePet from "./components/UpdatePet"; 
import RegisterPet from "./components/RegisterPet";
import Login from "./components/Login";
import Interests from "./components/Interests";
import Register from "./components/Register";
import Profile from "./components/Profile";
import SearchInterests from "./components/SearchInterests";
import InterestsUserProfile from "./components/InterestUserProfile";
import AvgPrice from "./components/Avgprice";
import EventBus from "./common/EventBus";

const App = () => {
  const [currentUser, setCurrentUser] = useState(undefined);

  useEffect(() => {
    const user = AuthService.getCurrentUser();

    if (user) {
      setCurrentUser(user);
    }

    EventBus.on("logout", () => {
      logOut();
    });

    return () => {
      EventBus.remove("logout");
    };
  }, []);

  const logOut = () => {
    AuthService.logout();
    setCurrentUser(undefined);
  };

  return (
    <div>
      <nav className="navbar navbar-expand navbar-dark bg-dark">
        <Link to={"/"} className="navbar-brand">
          Roomino
        </Link>
        <div className="navbar-nav mr-auto">
          <li className="nav-item">
            <Link to={"/"} className="nav-link">
              Home
            </Link>
          </li>
        </div>

        {currentUser ? (
          <div className="navbar-nav ml-auto">
            <li className="nav-item">
              <Link to={"/profile"} className="nav-link">
                {currentUser.username}
              </Link>
            </li>
            <li className="nav-item">
              <a href="/login" className="nav-link" onClick={logOut}>
                LogOut
              </a>
            </li>
          </div>
        ) : (
          <div className="navbar-nav ml-auto">
            <li className="nav-item">
              <Link to={"/login"} className="nav-link">
                Login
              </Link>
            </li>

            <li className="nav-item">
              <Link to={"/register"} className="nav-link">
                Sign Up
              </Link>
            </li>
          </div>
        )}
      </nav>

      <div className="container mt-3">
        <Routes>
          <Route path="/" element={<Login />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/profile" element={<Profile />} />
          <Route path="/profile/unitinfo" element={<UnitInfo />} /> 
          <Route path="/profile/updatepet" element={<UpdatePet />} /> 
          <Route path="/profile/registerpet" element={<RegisterPet />} /> 
          <Route path="/profile/interests" element={<Interests />} /> 
          <Route path="/profile/searchinterests" element={<SearchInterests />} /> 
          <Route path="/profile/interestprofile" element={<InterestsUserProfile />} /> 
          <Route path="/profile/avgprice" element={<AvgPrice />} /> 
        </Routes>
      </div>
    </div>
  );
};

export default App;
