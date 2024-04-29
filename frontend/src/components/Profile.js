import React, { useState, useEffect } from "react";
import { Navigate } from 'react-router-dom';
import AuthService from "../services/auth.service";
import UserService from '../services/user.service';
import "./Profile.css"; 

const genderOptions = [
  { value: 0, label: "Not known" },
  { value: 1, label: "Male" },
  { value: 2, label: "Female" },
  { value: 9, label: "Not applicable" }
];

const Profile = () => {
  const [userData, setUserData] = useState(null);
  const currentUser = AuthService.getCurrentUser();

  useEffect(() => {
    const fetchUserData = async () => {
      try {
        const data = await UserService.getUserBoard();
        setUserData(data.data);
      } catch (error) {
        console.error('Error fetching user data:', error);
      }
    };

    fetchUserData();
  }, []);

  const handleSearch = () => {
    window.location.href = '/profile/unitinfo';
  };

  const handleUpdatePet = () => {
    window.location.href = '/profile/updatepet';
  };

  const handleRegisterPet = () => {
    window.location.href = '/profile/registerpet';
  };
  const handleSearchInterests = () => {
    window.location.href = '/profile/searchinterests';
  };
  const handleAvgPrice = () => {
    window.location.href = '/profile/avgprice';
  };
  if(!currentUser) {
    return <Navigate to="/login" replace={true} />
  }

  return (
    <div className="container">
      <header className="jumbotron">
        <h3>
          <strong>{userData?.username || 'Username'}</strong> Profile
        </h3>
      </header>
      {userData && (
        <div>
          <p><strong>Username:</strong> {userData.username}</p>
          <p><strong>First Name:</strong> {userData.first_name}</p>
          <p><strong>Last Name:</strong> {userData.last_name}</p>
          <p><strong>Date of Birth:</strong> {new Date(userData.dob).toLocaleDateString('en-US', {
                          year: 'numeric',
                          month: '2-digit',
                          day: '2-digit'
                        })}</p>
          <p><strong>Gender:</strong> {genderOptions.find(option => option.value === userData.gender)?.label}</p>
          <p><strong>Email:</strong> {userData.email}</p>
          <p><strong>Phone:</strong> {userData.phone}</p>
        </div>
      )}
      <button className="btn btn-primary custom-btn" onClick={handleSearch}>Search Certain Apartment Units</button><br /><br />
      <button className="btn btn-primary custom-btn" onClick={handleRegisterPet}>Register Pet</button><br /><br />
      <button className="btn btn-primary custom-btn" onClick={handleUpdatePet}>Update Pet</button><br /><br />
      <button className="btn btn-primary custom-btn" onClick={handleSearchInterests}>Search Interest</button><br /><br />
      <button className="btn btn-primary custom-btn" onClick={handleAvgPrice}>Average Rent with Zipcode</button><br /><br />
    </div>
  );

};

export default Profile;
